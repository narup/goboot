package pweb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"runtime/debug"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/phil-inc/plib/core/util"
)

type body struct {
	Key string
}
type params struct {
	Key string
}

// Body key for request body
var Body = body{Key: "Body"}

//Params key for params
var Params = params{Key: "Params"}

// Errors represents json errors
type Errors struct {
	Errors []*Error `json:"errors"`
}

// Response response interface
type Response interface {
	Write(w http.ResponseWriter, r *http.Request)
}

// Error represents the API level error for the client apps
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

//ErrorHandler handler for all the middelware errors
type ErrorHandler interface {
	HandleError(r *http.Request, err error)
}

var (
	// UnAuthorized resource not found error
	UnAuthorized = &Error{"un_authorized", 401, "Request UnAuthorized", "Request must be authorized"}
	// Forbidden resource not found error
	Forbidden = &Error{"forbidden", 403, "Request Forbidden", "Request Forbidden"}
	// ErrNotFound resource not found error
	ErrNotFound = &Error{"not_found", 404, "Not found", "Data not found"}
	// ErrBadRequest bad request error
	ErrBadRequest = &Error{"bad_request", 400, "Bad request", "Request body is not well-formed. It must be JSON."}
	// ErrUnsupportedMediaType error
	ErrUnsupportedMediaType = &Error{"not_supported", 405, "Not supported", "Unsupported media type"}
	// ErrInternalServer error to represent server errors
	ErrInternalServer = &Error{"internal_server_error", 500, "Internal Server Error", "Something went wrong."}
)

// APIKeyAuth basically checks the authorization header for API Key
func APIKeyAuth(ctx context.Context, e ErrorHandler) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			apiKey, err := extractAPIKeyFromAuthHeader(r)
			if err != nil {
				log.Printf("[ERROR] Invalid authorization header %s\n", err)
				e.HandleError(r, err)
				WriteError(w, UnAuthorized)
			} else {
				if apiKey != util.Config("app.apiKey") {
					e.HandleError(r, errors.New("Invalid API Key. Unauthorized access"))
					WriteError(w, UnAuthorized)
				} else {
					// Delegate request to the given handle
					next.ServeHTTP(w, r)
				}
			}
		}
		return http.HandlerFunc(fn)
	}

	return m
}

// JWTAuthHandler checks and validate JWT token
func JWTAuthHandler(ctx context.Context, e ErrorHandler) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// check JSON web token data
			claims, err := checkJWT(w, r)
			// If there was an error, do not continue.
			if err != nil && err.Error() != "Token is expired" {
				log.Printf("[ERROR] Invalid authentication token: %v", err)
				e.HandleError(r, err)
				WriteError(w, UnAuthorized)
				return
			}
			if next != nil {
				if claims != nil {
					b := context.WithValue(r.Context(), util.SessionUserKey, claims)
					r = r.WithContext(b)
				}
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}

	return m
}

// JSONBodyHandler is a middleware to decode the JSON body, then set the body
// into the context.
func JSONBodyHandler(ctx context.Context, v interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(v)
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				WriteError(w, ErrBadRequest)
				return
			}

			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)
			if err != nil {
				log.Printf("[ERROR] Error decoding JSON data: %s", err)
				WriteError(w, ErrBadRequest)
				return
			}

			if next != nil {
				b := context.WithValue(r.Context(), Body, val)
				r = r.WithContext(b)
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// ResponseHandler handles the response from services
func ResponseHandler(f func(http.ResponseWriter, *http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := f(w, r)
		response.Write(w, r)
	}
}

// RecoverHandler is a deferred function that will recover from the panic,
// respond with a HTTP 500 error and log the panic. When our code panics in production
// (make sure it should not but we can forget things sometimes) our application
// will shutdown. We must catch panics, log them and keep the application running.
// It's pretty easy with Go and our middleware system.
func RecoverHandler(ctx context.Context, e ErrorHandler) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rr := recover(); rr != nil {
					log.Printf("PANIC: %s", debug.Stack())

					var err error
					switch x := rr.(type) {
					case string:
						err = errors.New(x)
					case error:
						err = x
					default:
						err = errors.New("Unknown panic")
					}
					if err != nil {
						e.HandleError(r, err)
					}
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// LoggingHandler middleware to log request/response
func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request:[%s] %q\n", r.Method, r.URL.String())
		//start time
		t1 := time.Now()
		//invoke next handler on a chain
		next.ServeHTTP(w, r)
		//end time
		t2 := time.Now()
		//log it!
		log.Printf("Reponse Time:[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

//ContentTypeHandler make sure content type is appplication/json for PUT/POST data
func ContentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			WriteError(w, ErrUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// WriteJSON writes resource to the output stream as JSON data.
func WriteJSON(w http.ResponseWriter, resource interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

// WriteError writes error response
func WriteError(w http.ResponseWriter, err *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(Errors{[]*Error{err}})
}

func checkJWT(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, error) {
	if r.Method == "OPTIONS" {
		return nil, nil
	}

	// Use the specified token extractor to extract a token from the request
	token, err := extractTokenFromAuthHeader(r)
	// If debugging is turned on, log the outcome
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, errors.New("Invalid auth token")
	}

	// Now parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[ERROR] Invalid signing method: %s", token.Signature)
		}
		return []byte(util.Config("auth.hmacToken")), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		return nil, errors.New("Invalid auth token")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid auth token")
}

// extractAPIKeyFromAuthHeader extract Phil API Key from the header
func extractAPIKeyFromAuthHeader(r *http.Request) (string, error) {
	authHeaderParts, err := getAuthHeaderParts(r)
	if err != nil {
		return "", err
	}
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "philkey" {
		return "", errors.New("Incorrect authorization header format. Invalid API Key")
	}
	return authHeaderParts[1], nil
}

// extractTokenFromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func extractTokenFromAuthHeader(r *http.Request) (string, error) {
	authHeaderParts, err := getAuthHeaderParts(r)
	if err != nil {
		return "", err
	}
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Incorrect authorization header format. Invalid access token")
	}

	return authHeaderParts[1], nil
}

func getAuthHeaderParts(r *http.Request) ([]string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return []string{""}, nil // No error, just no token
	}
	return strings.Split(authHeader, " "), nil
}

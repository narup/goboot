// Package aliz containing the core middleware functions for JWT and API key based authentication
// JSON body parser, request logger, recover handler and more.
// Copyright (c) Puran  2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package aliz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	gorilla "github.com/gorilla/context"
)

// APIKeyAuth basically checks the authorization header for API Key
func APIKeyAuth(ctx context.Context, apiKey string, e ErrorHandler) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			key, err := extractAPIKeyFromAuthHeader(r)
			if err != nil {
				log.Printf("[ERROR] Invalid authorization header %s\n", err)
				e.HandleError(r, err)
				WriteError(w, UnAuthorized)
			} else {
				if key != apiKey {
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
func JWTAuthHandler(ctx context.Context, secretAuthToken string, e ErrorHandler) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// check JSON web token data
			claims, err := checkJWT(w, r, secretAuthToken)
			// If there was an error, do not continue.
			if err != nil && err.Error() != "Token is expired" {
				log.Printf("[ERROR] Invalid authentication token: %v", err)
				e.HandleError(r, err)
				WriteError(w, UnAuthorized)
				return
			}
			if next != nil {
				if claims != nil {
					b := context.WithValue(r.Context(), SessionUserKey, claims)
					r = r.WithContext(b)
				}
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}

	return m
}

// ClearHandler wraps an http.Handler and clears request values at the end
// of a request lifetime.
func ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gorilla.Clear(r)
		h.ServeHTTP(w, r)
	})
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

// ResponseHandler handles the response from services and write it to the network output
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

func checkJWT(w http.ResponseWriter, r *http.Request, secretAuthToken string) (jwt.MapClaims, error) {
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
		return []byte(secretAuthToken), nil
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

// extractAPIKeyFromAuthHeader extract API Key from the header
func extractAPIKeyFromAuthHeader(r *http.Request) (string, error) {
	authHeaderParts, err := getAuthHeaderParts(r)
	if err != nil {
		return "", err
	}
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "apikey" {
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

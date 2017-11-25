// Package aliz containing the custom router which is a wrapper around httprouter (github.com/julienschmidt/httprouter)
// to make it compatible with http.Handler function. It uses wrapHandler function to wraps the middleware functions
// implementing http.Handler into a httprouter.Handler function.
//
// Example to setup routes using aliz.Router using alice(https://github.com/justinas/alice)
// middleware chaining:
//
//  package main
//
//  import (
//      "github.com/narup/aliz"
//      "net/http"
//      "fmt"
//  )
//
//  type AppPanicHandler struct {
//  }
//
//  func (eh AppPanicHandler) HandleError(r *http.Request, err error) {
//     //Handle panic error best suited for your application
//  }
//
//  func Index(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprint(w, "Welcome!\n")
//  }
//
//  var panicHandler = &AppPanicHandler{}
//  ctx := context.Background()
//  chain := alice.New(aliz.ClearHandler, aliz.LoggingHandler)
// 	chain.Append(aliz.RecoverHandler(ctx, panicHandler))
//
//  r := aliz.DefaultRouter(ctx)
//  r.Get("/", chain.ThenFunc(Index))
//
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package aliz

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

type sessionUser struct {
	Key string
}

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

// SessionUserKey key for context
var SessionUserKey = sessionUser{Key: "SessionUser"}

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

// Router wraps httprouter.Router, which is non-compatible with http.Handler to make it
// compatible by implementing http.Handler into a httprouter.Handler function.
type Router struct {
	r              *httprouter.Router
	Ctx            context.Context
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// DefaultRouter returns new aliz.Router with default settings
func DefaultRouter(ctx context.Context) *Router {
	ar := new(Router)
	ar.Ctx = ctx
	ar.r = httprouter.New()
	ar.AllowedOrigins = "*"
	ar.AllowedMethods = "POST, GET, OPTIONS, PUT, DELETE"
	ar.AllowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	return ar
}

//ServeHTTP handler function that takes care of headers
func (ar *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	if origin == "" || origin == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		if strings.Contains(ar.AllowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			WriteError(w, Forbidden)
			return
		}
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", ar.AllowedMethods)
	w.Header().Set("Access-Control-Allow-Headers", ar.AllowedHeaders)
	if req.Method == "OPTIONS" {
		w.(http.Flusher).Flush()
	}
	ar.r.ServeHTTP(w, req)
}

// Get wraps httprouter's GET function
func (ar *Router) Get(path string, handler http.Handler) {
	ar.r.GET(path, wrapHandler(ar.Ctx, handler))
}

// Post wraps httprouter's POST function
func (ar *Router) Post(path string, handler http.Handler) {
	ar.r.POST(path, wrapHandler(ar.Ctx, handler))
}

// Put wraps httprouter's PUT function
func (ar *Router) Put(path string, handler http.Handler) {
	ar.r.PUT(path, wrapHandler(ar.Ctx, handler))
}

// Delete wraps httprouter's DELETE function
func (ar *Router) Delete(path string, handler http.Handler) {
	ar.r.DELETE(path, wrapHandler(ar.Ctx, handler))
}

// wrapHandler wraps http.Handler middleware function inside httprouter.Handle
func wrapHandler(ctx context.Context, h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//instead of passing extra params to handler function use context
		if ps != nil {
			ctxParams := context.WithValue(r.Context(), Params, ps)
			r = r.WithContext(ctxParams)
		}
		h.ServeHTTP(w, r)
	}
}

// ErrMissingRequiredData error to represent missing data error
var ErrMissingRequiredData = errors.New("missing required data")

//ErrNotRecognized error for any unrecognized client
var ErrNotRecognized = errors.New("not recognized")

// APIResponse response data representation for API
type APIResponse struct {
	Error  string      `json:"error,omitempty"`
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// Write - Reponse interface implementation
func (res APIResponse) Write(w http.ResponseWriter, r *http.Request) {
	if res.Status == "ERROR" {
		log.Printf("[ERROR][API][PATH: %s]:: Error handling request. ERROR: %s. User agent: %s", r.RequestURI, res.Error, r.Header.Get("User-Agent"))
	}
	WriteJSON(w, res)
}

// DataResponse creates new API data response using the resource
func DataResponse(data interface{}) APIResponse {
	return APIResponse{Error: "", Status: "OK", Data: data}
}

// StringErrorResponse constructs error response based on input
func StringErrorResponse(err string) APIResponse {
	return APIResponse{Error: err, Status: "ERROR", Data: nil}
}

//ErrorResponse constructs error response from the API
func ErrorResponse(err error) APIResponse {
	return APIResponse{Error: err.Error(), Status: "ERROR", Data: nil}
}

// RequestBody returns the request body
func RequestBody(r *http.Request) interface{} {
	return r.Context().Value(Body)
}

// SessionUserID returns user id of the current session
func SessionUserID(r *http.Request) string {
	if jwtClaims, ok := r.Context().Value(SessionUserKey).(jwt.MapClaims); ok {
		return jwtClaims["uid"].(string)
	}
	return ""
}

// UserRoles current user roles
func UserRoles(r *http.Request) []string {
	if jwtClaims, ok := r.Context().Value(SessionUserKey).(jwt.MapClaims); ok {
		return jwtClaims["uid"].([]string)
	}
	return make([]string, 0)
}

// QueryParamByName returns the request param by name
func QueryParamByName(name string, r *http.Request) string {
	return r.URL.Query().Get(name)
}

// QueryParamsByName returns the request param by name
func QueryParamsByName(name string, r *http.Request) []string {
	values := r.URL.Query()
	return values[name]
}

// ParamByName returns the request param by name
func ParamByName(name string, r *http.Request) string {
	params := r.Context().Value(Params).(httprouter.Params)
	return params.ByName(name)
}

//Authorize checks if given request is authorized
func Authorize(w http.ResponseWriter, r *http.Request) {
	sid := SessionUserID(r)
	uid := ParamByName("uid", r)

	if sid != uid {
		WriteError(w, Forbidden)
	}
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

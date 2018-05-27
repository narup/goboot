# goboot 
## ![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103) ![Open Source Love](https://badges.frapsoft.com/os/v2/open-source.svg?v=103)
Ultralight Go web framework for bootstrapping microservices. It's built on top of popular httprouter [httprouter](github.com/julienschmidt/httprouter) with all the necessary middlewares hooked to get you started easily. 

## Features 
* Easy routes setup with [alice middleware chaining](https://github.com/justinas/alice)

    ```go
    ctx := context.Background()
    panicHandler := &AppPanicHandler{}

    chain := alice.New(go.ClearHandler, go.LoggingHandler)
    chain.Append(go.RecoverHandler(ctx, panicHandler))

    r.Get("/", chain.ThenFunc(Index))
    ```
* Simple and consistent controller spec. Just return go.Response type from the controller
	```go
	func SignUp(w http.ResponseWriter, r *http.Request) go.Response {
		us := go.RequestBody(r).(*User)
		
		savedUser, err := service.SaveUser(us)
		if err != nil {
			return go.ErrorResponse(err)
		}
	
		return go.DataResponse(savedUser)
	}
	```
	Each controller can either return an error or success API response. JSON output that is written to the client is:
    
	* ErrorResponse
	```JSON
	{
	    "status" : "ERROR",
	    "error"  : "<api error message>",
	    "data"   : "error data, if name"
	}
	```
	* DataResponse
	```JSON
	{
	    "status" : "OK",
	    "data"   : {"name":"Puran S", "email":"puran@myemail.com"}
	}
	```
* Simple JSON POST body handler marshal's your JSON request to go struct 
     - Define your request data struct
     ```go
     type newsArticle struct {
		ID            bson.ObjectId `json:"id"`
		Title         string        `json:"title"`
		PublishedDate string        `json:"date_published"`
     }
     ```
    
    - Define route
    ```go
    jsonHandler := go.JSONBodyHandler(ctx, newsArticle{})
    r.Post("/api/v1/articles", chain.Append(jsonHandler).ThenFunc(go.ResponseHandler(SaveNewsArticle)))
    ```
	
    - POST request handler
    ```go
    func SaveNewsArticle(w http.ResponseWriter, r *http.Request) go.Response {
	  newArtcile := go.RequestBody(r).(*newsArticle)
          ........
    }
    ```
    
* Supports API Key or JWT Auth Token for security
    * For API Key based security append the middleware ```APIKeyAuth``` chain
    ```go
	apiKeyAuthC := chain.Append(go.APIKeyAuth(ctx, apiErrHandler))
    ```
    * For auth token based security use ```JWTAuthHandler```
    ```go
    jwtChain := chain.Append(go.JWTAuthHandler(ctx, apiErrHandler))
    ```


#### Sample API service web.go putting it all together:

 ```go
 package main

 import (
     "github.com/narup/go"
     "net/http"
     "fmt"
 )

 type AppPanicHandler struct {
 }
 
 func (eh AppPanicHandler) HandleError(r *http.Request, err error) {
    //Handle panic error best suited for your application
 }
 
 func Index(w http.ResponseWriter, r *http.Request) {
     fmt.Fprint(w, "Welcome!\n")
 }

 func SignUp(w http.ResponseWriter, r *http.Request) go.Response {
	us := go.RequestBody(r).(*User)
	
	savedUser, err := service.SaveUser(us)
	if err != nil {
	    return go.ErrorResponse(err)
	}
	
	return go.DataResponse(savedUser)
 }
 
 func handlers() *go.Router {
     ctx := context.Background()
     var panicHandler = &AppPanicHandler{}
   
     chain := alice.New(go.ClearHandler, go.LoggingHandler)
     chain.Append(go.RecoverHandler(ctx, panicHandler))

     r := go.DefaultRouter(ctx)
     
     //setup routes
     r.Get("/", chain.ThenFunc(Index))
     
     userJSONHandler := go.JSONBodyHandler(ctx, User{})
     r.Post("/api/v1/users", chain.Append(userJSONHandler).ThenFunc(go.ResponseHandler(SignUp)))
 }
 
 func main() {
 	port := "8080"
	fmt.Printf("Starting server on port: %s.... %s \n", port)
	log.Println("Press ctrl+E to stop the server.")
	srv := &http.Server{
		Handler:      handlers(),
		Addr:         port,
		ReadTimeout:  4 * time.Minute,
		WriteTimeout: 8 * time.Minute,
	}
	if serr := srv.ListenAndServe(); serr != nil {
		log.Fatalf("Error starting server: %s\n", serr)
	}
} 
```

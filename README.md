# Aliz
Ultralight Go web framework for writing microservices

#### Sample API service web.go

```go
 package main

 import (
     "github.com/narup/aliz"
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

 func SignUp(w http.ResponseWriter, r *http.Request) aliz.Response {
	us := pweb.RequestBody(r).(*User)
	
	savedUser, err := service.SaveUser(us)
	if err != nil {
	    return aliz.ErrorResponse(err)
	}
	
	return aliz.DataResponse(savedUser)
 }
 
 func handlers() *aliz.Router {
     ctx := context.Background()
     var panicHandler = &AppPanicHandler{}
   
     chain := alice.New(aliz.ClearHandler, aliz.LoggingHandler)
     chain.Append(aliz.RecoverHandler(ctx, panicHandler))

     r := aliz.DefaultRouter(ctx)
     
     //setup routes
     r.Get("/", chain.ThenFunc(Index))
     
     userJSONHandler := aliz.JSONBodyHandler(ctx, User{})
     r.Post("/api/v1/users", chain.Append(userJSONHandler).ThenFunc(aliz.ResponseHandler(SignUp)))
 }
 
 func main() {
 	port := "8080"
	fmt.Printf("Starting phil API server on port: %s.... %s \n", port)
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

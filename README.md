# Aliz
Ultralight Go web framework for writing microservices

```
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
 
 var panicHandler = &AppPanicHandler{}
 ctx := context.Background()
 chain := alice.New(aliz.ClearHandler, aliz.LoggingHandler)
 chain.Append(aliz.RecoverHandler(ctx, panicHandler))

 r := aliz.DefaultRouter(ctx)
 
 //setup routes
 r.Get("/", chain.ThenFunc(Index))
 
```

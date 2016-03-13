# FrameGoWork
Golang micro web framework.


## Use
- Import ```"router"``` library.
- Write a function with params: ```http.ResponseWriter, *http.Request & map[string]string```:
``` 
  func echo(w http.ResponseWriter, r *http.Request, params map[string]string) {
    fmt.Fprintf(w, "Hello, "+params["msg"])
  } 
```
- Write, if you want, a middleware function with params: ```http.ResponseWriter, *http.Request & router.NextHandler```
```
  func middleware(w http.ResponseWriter, r *http.Request,  next router.NextHandler) {
    fmt.Fprintf(w, "Hey!\n")
    next.Exec(w, r)
  }
```
- Instance ```router``` struct:
```
  router := router.New()
```
- Set path, method and callback:
```
  router.GET("/echo/:msg", echo)
```
or
```
  router.GET("/echo/:msg", echo, middleware)
```
- Set port number and start server:
```
  router.Run("9999")
```
- Open ```localhost:9999/echo/world```on your browser.

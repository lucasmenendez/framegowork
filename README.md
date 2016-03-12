# FrameGoWork
Golang micro web framework.


## Use
- Import ```"router"``` library.
- Write a function with params: ```http.ResponseWriter, *http.Request & map[string]string```:
``` 
  func helloWorld(w http.ResponseWriter, r *http.Request, params map[string]string) {
    fmt.Fprintf(w, "Hello, " + params["msg"])
  } 
```
- Instance ```router``` struct:
```
  router := router.New()
```
- Set path, method and callback:
```
  router.GET("/echo/:msg", helloWorld)
```
- Set port number and start server:
```
  router.Run("9999")
```
- Open ```localhost:9999/echo/world```on your browser.

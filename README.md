# FrameGoWork
Golang micro web framework.


## Use
- Import ```"router"``` library.
- Write a function with params: ```http.ResponseWriter``` & ```*http.Request```:
``` 
  func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")
  } 
```
- Instance ```router``` struct:
```
  router := router.New()
```
- Set path, method and callback:
```
  router.GET("/hellow-world", helloWorld)
```
- Set port number and start server:
```
  router.RunServer("9999")
```
- Open ```localhost:9999/hello-world```on your browser.

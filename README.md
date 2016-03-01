# FrameGoWork
Golang micro web framework.


## Use
0. Import ```"router"``` library.
1. Write a function with params: ```http.ResponseWriter``` & ```*http.Request```:
``` 
  func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")
  } 
```
2. Instance ```router``` struct:
```
  router := router.New()
```
3. Set path, method and callback:
```
  router.GET("/hellow-world", helloWorld)
```
4. Set port number and start server:
```
  router.RunServer("9999")
```
5. Open ```localhost:9999/hello-world```on your browser.

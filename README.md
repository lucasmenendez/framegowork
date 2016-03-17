# FrameGoWork
Golang micro web framework.


## Demo
- Import all repository
- Exec ```go run main.go```
- Open ```localhost:9999/echo/world```


## Tools
- ```username, password, ok := tools.BasicAuth(authHeader string) (string, string, bool)
- ```json, err := tools.ToJSON(data interface{}) ([]byte, error)


## Use

#### Import
Import ```"server"``` library.

#### Custom route type
Write a function with params: ```http.ResponseWriter, *http.Request & map[string]string```:
``` 
  func echo(w http.ResponseWriter, r *http.Request, params map[string]string) {
    fmt.Fprintf(w, "Hello, "+params["msg"])
  } 
```

#### Middleware
Write, if you want, a middleware function with params: ```http.ResponseWriter, *http.Request & server.NextHandler```
```
  func middleware(w http.ResponseWriter, r *http.Request,  next server.NextHandler) {
    fmt.Fprintf(w, "Hey!\n")
    next.Exec(w, r)
  }
```

#### Create router
Instance ```server``` struct with config params:
```
  server := server.New()
```

#### Config server
##### Headers
Add universal headers to server. Can be overwrited with common method.
```
  server.SetHeader("Content-Type", "text/plain")
  server.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
  server.SetHeader("Access-Control-Allow-Origin", "*")
  server.SetHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
```
##### Port
Set port server. Default ```9999```
```
  server.SetPort("9999")
```

#### Instance routes
Set path, method and callback:
```
  server.GET("/echo/:msg", echo)
```
or
```
  server.POST("/echo/:msg", echo, middleware)
```

#### Run server
Set port number and start server:
```
  server.Run()
```

Open ```localhost:9999/echo/world```on your browser.

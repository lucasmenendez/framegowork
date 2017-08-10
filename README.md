# framework.go
Golang micro web framework.

## Install

```
go get github.com/lucasmenendez/framework.go	
```

## Demo

```go
package main

import (
	"fmt"
	f "github.com/lucasmenendez/framework.go"
)

func echo(w f.Response, r f.Request, params f.Params) {
	fmt.Fprintf(w, "Hello, "+params["msg"])
}

func middleware(w f.Response, r f.Request, next f.NextHandler) {
	fmt.Fprintf(w, "Hey!\n")
	next.Exec(w, r)
}

func main() {
	server := f.New()
	server.SetPort(9999)
	server.DebugMode(true)

	server.GET("/echo/:msg", echo, middleware)
	server.Run()
}
```

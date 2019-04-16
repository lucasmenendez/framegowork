[![GoDoc](https://godoc.org/github.com/lucasmenendez/shgf?status.svg)](https://godoc.org/github.com/lucasmenendez/shgf)
[![Build Status](https://travis-ci.org/lucasmenendez/shgf.svg?branch=master)](https://travis-ci.org/lucasmenendez/shgf)
[![Report](https://goreportcard.com/badge/github.com/lucasmenendez/shgf)](https://goreportcard.com/report/github.com/lucasmenendez/shgf)
[![codebeat badge](https://codebeat.co/badges/2cf14aa6-9240-447e-ac7e-0c620fc6cf99)](https://codebeat.co/projects/github-com-lucasmenendez-shgf-master)


# SHGF: Simple HTTP golang framework
**S**imple **H**TTP **G**olang **F**ramework. Provides simple API to create an HTTP server and routes with dynamic paths, registered by HTTP method.

## Install

```sh
go get github.com/lucasmenendez/shgf	
```

## Documentation
Read all the reference documents into [GoDoc](https://godoc.org/github.com/lucasmenendez/shgf) article.

### Import the package

You can import it like this:

```go
import "github.com/lucasmenendez/shgf"
```

### Main Handler

Main handler function represents the main entry point to the route assigned to it.
It must be a `shgf.Handler`. You can checkout the documentation about the `shgf.Handler` in its [GoDoc page](https://godoc.org/github.com/lucasmenendez/shgf#Handler).

```go
func req(ctx *shgf.Context) (res *shgf.Response) {
	var err error
	var msg = fmt.Sprintf("Hello %s!", ctx.Params["msg"])
	if res, err = shgf.NewResponse(200, msg); err != nil {
		res, _ = shgf.NewResponse(500, err)
	}

	return res
}
```

### Middleware Handler

Middleware function is executed before main handler function, but provides entire access to request data and functions. Commonly used for validate request data before process the main function, for example in a process of authentication. In that case, the middleware is used for parsing the route params calling `ctx.ParseParams()` function. Like the *Main Handler*, it must be a `shgf.Handler`. 

```go
func mid(ctx *shgf.Context) (res *shgf.Response) {
	if err := ctx.ParseParams(); err != nil {
		res, _ := shgf.NewResponse(500, err)
		return res
	}

	return ctx.Next()
}
```

### Route

Registering new route, the server is prepared to listen to requests on the route's path and handling it calling the route's handler. The route is registered calling to the different functions named like the HTTP verbs.

#### Route params

Route registration admits dynamic paths to define typed params inside. The params must have the following format: `<type:alias>`.

```go
s.GET("/hello/<string:msg>", req, mid)
```

The supported types are:

|  **Type**  |   **Example**  |
|:----------:|:--------------:|
|   `bool`   |  `<bool:foo>`  |
|    `int`   |   `<int:foo>`  |
|   `float`  |  `<float:foo>` |
|  `string`  | `<string:foo>` |


### Server

Server instance allows to developer to configure the behavior of the server and register new routes to handle it. For more information about `shgf.Server` checkout the documentation [here](https://godoc.org/github.com/lucasmenendez/shgf#Server).

```go
s, err := shgf.New(shgf.LocalConf())
if err != nil {
	fmt.Println(err)
}

// New routes go here

s.Listen()

```

#### Server Configuration

Configuration instance allows to set some parameters to define server behavior and properties. For an **advance configuration**, read full docs about `shgf.Config` into the GoDoc [article](https://godoc.org/github.com/lucasmenendez/shgf#Config).

##### Local Host configuration

LocalConf method returns the default configuration to deploy localhost server. It is configured with debug enabled, and serve into `127.0.0.1:8080`.

```go
/*
var conf = &shgf.Config{
	Hostname: "127.0.0.1",
	Port: 8080,
	PortTLS: 8081,
	Debug: true,
}
*/

var conf = shgf.LocalConf()
```

##### Basic HTTP configuration

BasicConf method receives server hostname and port number, and returns the configuration to deploy a simple HTTP server with that parameters.

```go
/*
var conf = &shgf.Config{
	Hostname: "0.0.0.0",
	Port: 8080,
}
*/

var conf = shgf.BasicConf("0.0.0.0", 8080)
```

##### Enable TLS

To enable TLS, append a TLS port number and a `cert` and `key` files path to the any `shgf.Config`:

```go
var conf = shgf.BasicConf("0.0.0.0", 8080)
conf.TLSPort = 8081
conf.TLSCert = "/path/to/cert.pem"
conf.TLSKey = "/path/to/key.pem"
```

##### Enable HTTP2

To enable HTTP2, TLS must be configured.

```go
var conf = &shgf.Config{
	Hostname: "127.0.0.1",
	Port:     8080,
	PortTLS:  8081,
	TLSCert:  "/path/to/cert.pem",
	TLSKey:   "/path/to/key.pem",
	HTTP2:    true,
}
```
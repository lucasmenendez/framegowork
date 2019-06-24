[![GoDoc](https://godoc.org/github.com/lucasmenendez/shgf?status.svg)](https://godoc.org/github.com/lucasmenendez/shgf)
[![Build Status](https://travis-ci.org/lucasmenendez/shgf.svg?branch=master)](https://travis-ci.org/lucasmenendez/shgf)
[![Report](https://goreportcard.com/badge/github.com/lucasmenendez/shgf)](https://goreportcard.com/report/github.com/lucasmenendez/shgf)
[![codebeat badge](https://codebeat.co/badges/2cf14aa6-9240-447e-ac7e-0c620fc6cf99)](https://codebeat.co/projects/github-com-lucasmenendez-shgf-master)


# SHGF: Simple HTTP golang framework
**S**imple **H**TTP **G**olang **F**ramework. Provides simple API to create an HTTP server and routes with dynamic paths, registered by HTTP method.


## Main features

* Handle URL by path and method
* Register dynamic paths with typed params
* Parse forms easely
* TLS & HTTP/2 

## Reference
Read all the reference documents into [GoDoc](https://godoc.org/github.com/lucasmenendez/shgf) article.

## Installation

```sh
go get github.com/lucasmenendez/shgf	
```

---

## Documentation

1. [Including `shgf` on your project](#including-shgf-on-your-project)
2. [Initializing `shgf.Server`](#initializing-shgf.Server)
   - Server configuration
      - Local Host configuration
	  - Basic HTTP configuration
	  - Enable TLS
	  - Enable HTTP/2
3. [Routes](#routes)
   - Registring new routes
   - Routes params
   - Parsing route params
4. [Handlers](#handlers)
   - Main Handler
   - Middleware Handler
5. [Forms](#forms)
   - Parsing forms

### Including `shgf` on your project

You can import it like this:

```go
import "github.com/lucasmenendez/shgf"
```

### Initializing `shgf.Server`

Server instance allows to developer to configure the behavior of the server and register new routes to handle it. For more information about `shgf.Server` checkout the documentation [here](https://godoc.org/github.com/lucasmenendez/shgf#Server).

```go
server, err := shgf.New(shgf.LocalConf())
if err != nil {
	fmt.Println(err)
}

// new routes go here

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

##### Enable HTTP/2

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

### Routes

Registering new route, the server is prepared to listen to requests on the route's path and handling it calling the route's handler. The route is registered calling to the different functions named like the HTTP verbs.

#### Registring new routes 

New routes must be assigned to an existing server. Register new routes calling desired server method according to the chosen HTTP method, passing the desired path to listen and one [handler](#handlers) at least. Check out more information about this [here](https://godoc.org/github.com/lucasmenendez/shgf#Server);

```go
server.GET("/hello/<string:msg>", req, mid)
server.POST("/hello/<string:msg>", req, mid)
```

#### Routes params

Route registration admits dynamic paths to define typed params inside. The params must have the following format: `<type:alias>`. After parsing params, its possible to access them using `context.Params` (check how to parse params in the [following section](#parsing-route-params)). `context.Params` its a `map[string]interface{}` with functions like `Exists()` or `Get()` that provides secure API. 


The supported types are:

|  **Type**  |   **Example**  |
|:----------:|:--------------:|
|   `bool`   |  `<bool:foo>`  |
|    `int`   |   `<int:foo>`  |
|   `float`  |  `<float:foo>` |
|  `string`  | `<string:foo>` |


#### Parsing route params

To use the route params into a handler function, it must first be parsed calling `context.ParseParams()`. Then, all the params will be accessible from `context.Params` (a map of `string` and `interface{}`):

```go

// handling the route '/hello/<string:foo>'
if err := ctx.ParseParams(); err != nil {
	// catch error
}

foo, _ := ctx.Params.Get("foo") // get foo value safely
fmt.Println(foo) // prints path foo item value
fmt.Println(ctx.Params["foo"]) // or get value like a map
```

### Handlers

Handlers are functions that are executed when some client makes a request to the server. Handlers generally respond under conditions such as URL or method  matching those determined by the developer.

#### Main Handler

Main handler function represents the main entry point to the route assigned to it.
It must be a `shgf.Handler`. You can checkout the documentation about the `shgf.Handler` in its [GoDoc page](https://godoc.org/github.com/lucasmenendez/shgf#Handler).

```go
func req(ctx *shgf.Context) (res *shgf.Response) {
	var err error
	if res, err = shgf.NewResponse(200, "Hello world!"); err != nil {
		res, _ = shgf.NewResponse(500, err)
	}

	return res
}
```

#### Middleware Handler

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

### Forms

One of the ways to send user data between client and server is via Forms. Forms is a protocol convention that includes a specific formats based on HTTP Content-Type header previusly defined. There are two options:

| **Content-Type** | **Raw Body** |
|:---------------------:|:---------------:|
| `multipart/form-data` | `Content-Disposition: form-data; name="foo" bar` |
| `application/x-www-form-urlencoded` | `foo=bar` |

After parsing form, its possible to access them using `context.Form` (check how to parse params in the [following section](#parsing-forms)). `context.Form` its a `map[string]interface{}` with functions like `Exists()` or `Get()` that provides secure API. 

#### Parsing forms

To use data from Forms, it must first be parsed calling `context.ParseForm()`. Then, all the params will be accessible from `context.Form` (a map of `string` and `interface{}`):

```go
/** 
	Handling:
	- GET 'application/x-www-form-urlencoded' request to '/hello' with body 'foo=bar'
	- POST 'multipart/form-data' request to '/hello' with body 'Content-Disposition: form-data; name="foo" bar'
*/ 
if err := ctx.ParseForm(); err != nil {
	// catch error
}

foo, _ := ctx.Form.Get("foo") // get foo value safely
fmt.Println(foo) // prints form foo field value
fmt.Println(ctx.Form["foo"]) // or get value like a map
```
package shgf

type route struct {
	path        string
	middlewares handler
	handlers    handlers
	params      params
}

type routes []route

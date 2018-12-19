package shgf

type Handler func(Context)

//Exec next function when route has a middleware
func (c Context) Continue() {
	var next Handler = c.Handler
	c.Handler = nil
	next(c)
}

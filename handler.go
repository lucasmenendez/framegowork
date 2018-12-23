package shgf

type handler func(req request, ctx context) response
type handlers []handler

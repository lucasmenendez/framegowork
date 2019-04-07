package shgf

// routes type wraps a slices of route
type routes []*route

// exists function returns if provided route exists into the current routes
// list. A existing route has the same method and path that one route of
// current list at least.
func (rl routes) exists(n route) bool {
	for _, r := range rl {
		if r.method == n.method && r.path == n.path {
			return true
		}
	}

	return false
}

// findByPath function returns a sublist of current list filtered by path
// provided.
func (rl routes) findByPath(path string) routes {
	var res = routes{}
	for _, r := range rl {
		if r.path == path {
			res = append(res, r)
		}
	}

	return res
}

// findByMatchingPath function returns a sublist of current list routes that
// match with path provided.
func (rl routes) findByMatchingPath(path string) routes {
	var res = routes{}
	for _, r := range rl {
		if r.matcher.MatchString(path) {
			res = append(res, r)
		}
	}

	return res
}

// findByPath function returns a sublist of current list filtered by method
// provided.
func (rl routes) findByMethod(method string) routes {
	var res = routes{}
	for _, r := range rl {
		if r.method == method {
			res = append(res, r)
		}
	}

	return res
}

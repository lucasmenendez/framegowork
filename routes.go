package shgf

type routes []*route

func (rl routes) exists(n route) bool {
	for _, r := range rl {
		if r.method == n.method && r.path == n.path {
			return true
		}
	}

	return false
}

func (rl routes) findByPath(path string) routes {
	var res = routes{}
	for _, r := range rl {
		if r.path == path {
			res = append(res, r)
		}
	}

	return res
}

func (rl routes) findByMatchingPath(path string) routes {
	var res = routes{}
	for _, r := range rl {
		if r.matcher.MatchString(path) {
			res = append(res, r)
		}
	}

	return res
}

func (rl routes) findByMethod(method string) routes {
	var res = routes{}
	for _, r := range rl {
		if r.method == method {
			res = append(res, r)
		}
	}

	return res
}

package shgf

type shgfErr struct {
	err []interface{}
	dev []string
}

func (e shgfErr) Error() string { return "" }

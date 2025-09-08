package context

type Context struct {
	Debug bool
}

var (
	context *Context
)

func GetContext() *Context {
	if context != nil {
		return context
	}
	context = &Context{Debug: false}
	return context
}

func (c *Context) SetContext(debug bool) {
	c.Debug = debug
}

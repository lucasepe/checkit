package render

import "io"

type Renderer interface {
	Render(io.Reader) error
}

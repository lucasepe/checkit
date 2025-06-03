package render

import (
	"github.com/lucasepe/checkit/internal/parser"
)

type Renderer interface {
	Render(*parser.CheckList) error
}

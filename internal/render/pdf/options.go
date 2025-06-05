package pdf

import (
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

type RenderOption func(g *pdfRenderImpl)

func Square(on bool) RenderOption {
	return func(g *pdfRenderImpl) {
		if on {
			g.opts = squareOptions()
		} else {
			g.opts = defaultOptions()
		}
	}
}

func OutputDir(dir string) RenderOption {
	return func(g *pdfRenderImpl) {
		g.outputDir = strings.TrimSpace(dir)
		if g.outputDir == "" {
			g.outputDir, _ = os.Getwd()
		}
	}
}

type pdfRenderOptions struct {
	pageWidth  float64
	pageHeight float64

	marginTop    float64
	marginBottom float64
	marginLeft   float64

	groupTitleFontSize float64
	groupTitleMargin   float64

	itemFontSize float64
	itemMargin   float64

	itemNoteFontSize float64
	itemNoteMargin   float64

	documentTitleFontSize float64

	lineSpacing float64
}

func defaultOptions() pdfRenderOptions {
	return pdfRenderOptions{
		pageWidth:  float64(gopdf.PageSizeA4.W),
		pageHeight: float64(gopdf.PageSizeA4.H),

		marginTop:    40.0,
		marginBottom: 40.0,
		marginLeft:   40.0,
	}
}

func squareOptions() pdfRenderOptions {
	opts := pdfRenderOptions{
		pageWidth:  1080,
		pageHeight: 1080,

		marginTop:    8,
		marginBottom: 8,
		marginLeft:   8,
	}

	return opts
}

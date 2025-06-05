package pdf

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasepe/checkit/internal/parser"
	"github.com/lucasepe/checkit/internal/render"
	"github.com/lucasepe/x/text/slugify"
	"github.com/signintech/gopdf"
	"golang.org/x/image/font/gofont/gomono"
)

func New(opts ...RenderOption) (render.Renderer, error) {
	g := &pdfRenderImpl{
		doc:  &gopdf.GoPdf{},
		opts: defaultOptions(),
	}

	for _, opt := range opts {
		opt(g)
	}

	if err := g.setup(); err != nil {
		return g, err
	}

	g.doc.Start(gopdf.Config{
		Unit: gopdf.UnitPT,
		PageSize: gopdf.Rect{
			W: g.opts.pageWidth,
			H: g.opts.pageHeight,
		},
	})
	err := g.doc.AddTTFFontData(fontName, gomono.TTF)
	if err != nil {
		return g, err
	}

	g.doc.AddPage()
	g.pageCount = 1

	return g, nil
}

const (
	fontName       = "GoMono"
	symbol         = "\u25CB" // ○
	itemIndent     = "   "    // Indent for wrapped lines
	itemNoteIndent = "     "  // Indent for wrapped lines

	defaultCreator = "Check IT (github.com/lucasepe/checkit)"
)

var _ render.Renderer = (*pdfRenderImpl)(nil)

type pdfRenderImpl struct {
	doc       *gopdf.GoPdf
	opts      pdfRenderOptions
	pageCount int
	lineCount int
	filename  string
	outputDir string
}

func (g *pdfRenderImpl) Render(lst *parser.CheckList) (err error) {
	y := g.opts.marginTop

	if lst.Title != "" {
		g.setMeta(lst.Title)
		y, err = g.handleDocumentTitle(y, lst.Title)
		if err != nil {
			return err
		}
	}

	for _, grp := range lst.Groups {
		y, err = g.handleGroup(y, grp.Title)
		if err != nil {
			return err
		}

		for _, it := range grp.Items {
			y, err = g.handleItem(y, it.Title)
			if err != nil {
				return err
			}

			for _, nt := range it.Notes {
				y, err = g.handleItemNote(y, nt)
				if err != nil {
					return err
				}
			}
		}
	}

	return g.savePDF()
}

func (g *pdfRenderImpl) setup() error {
	fontSize := 0.02 * math.Min(g.opts.pageWidth, g.opts.pageHeight)

	g.opts.itemFontSize = fontSize
	g.opts.lineSpacing = 0.5 * fontSize
	g.opts.itemMargin = 0.8 * fontSize
	g.opts.itemNoteFontSize = 0.75 * fontSize
	g.opts.itemNoteMargin = 1.15 * (0.75 * fontSize)
	g.opts.groupTitleFontSize = fontSize + 4
	g.opts.groupTitleMargin = 1.1 * (fontSize + 4)
	g.opts.documentTitleFontSize = (fontSize + 4) + 4

	return os.MkdirAll(g.outputDir, 0755)
}

func (g *pdfRenderImpl) setMeta(title string) {
	author, ok := os.LookupEnv("USER")
	if !ok {
		author, ok = os.LookupEnv("USERNAME")
	}
	if !ok {
		author = defaultCreator
	}

	g.doc.SetInfo(gopdf.PdfInfo{
		Title:        title,
		Author:       author,
		Subject:      title,
		Creator:      defaultCreator,
		Producer:     defaultCreator,
		CreationDate: time.Now(),
	})
}

func (g *pdfRenderImpl) handleGroup(y float64, title string) (float64, error) {
	if y+g.opts.groupTitleFontSize+g.opts.lineSpacing > g.opts.pageHeight-g.opts.marginBottom {
		g.doc.AddPage()
		g.pageCount++

		y = g.opts.marginTop
	} else {
		y += g.opts.groupTitleMargin
	}

	g.doc.SetFont(fontName, "", g.opts.groupTitleFontSize)
	g.doc.SetX(g.opts.marginLeft)
	g.doc.SetY(y)
	g.doc.Text(title)

	y += g.opts.groupTitleFontSize //+ 0.5*g.opts.groupTitleMargin

	return y, nil
}

func (g *pdfRenderImpl) handleDocumentTitle(y float64, title string) (float64, error) {
	err := g.doc.SetFont(fontName, "", g.opts.documentTitleFontSize)
	if err != nil {
		return y, err
	}

	g.filename = fmt.Sprintf("%s.pdf", slugify.Sprint(title))

	// Calcola larghezza testo
	tw, err := g.doc.MeasureTextWidth(title)
	if err != nil {
		return y, err
	}

	th, err := g.doc.MeasureCellHeightByText(title)
	if err != nil {
		return y, err
	}

	y += th

	// Centra il testo orizzontalmente
	titleX := (g.opts.pageWidth - tw) / 2

	// Posiziona in alto
	g.doc.SetX(titleX)
	g.doc.SetY(y)
	g.doc.Text(title)

	// Sposta `y` sotto il titolo
	y += g.opts.documentTitleFontSize + g.opts.groupTitleMargin

	return y, nil
}

func (g *pdfRenderImpl) handleItem(y float64, line string) (float64, error) {
	err := g.doc.SetFont(fontName, "", g.opts.itemFontSize) // reset to item font
	if err != nil {
		return y, err
	}

	prefix := fmt.Sprintf(" %s ", symbol)
	maxTextWidth := g.opts.pageWidth - 2*g.opts.marginLeft
	wrappedLines := wrapTextWithPrefix(g.doc, line, prefix, itemIndent, maxTextWidth)

	y += g.opts.itemMargin

	for _, l := range wrappedLines {
		if y+g.opts.itemFontSize+g.opts.lineSpacing > g.opts.pageHeight-g.opts.marginBottom {
			g.doc.AddPage()
			g.pageCount++

			th, err := g.doc.MeasureCellHeightByText(l)
			if err != nil {
				return y, err
			}

			y = g.opts.marginTop + th
		}

		g.doc.SetX(g.opts.marginLeft)
		g.doc.SetY(y)
		g.doc.Text(l)
		y += g.opts.itemFontSize + g.opts.itemMargin
	}

	return y, nil
}

func (g *pdfRenderImpl) handleItemNote(y float64, line string) (float64, error) {
	err := g.doc.SetFont(fontName, "", g.opts.itemNoteFontSize) // reset to item font
	if err != nil {
		return y, err
	}

	maxTextWidth := g.opts.pageWidth - 2*g.opts.marginLeft
	wrappedLines := wrapTextWithPrefix(g.doc, line, itemNoteIndent, itemNoteIndent, maxTextWidth)

	for _, l := range wrappedLines {
		if y+g.opts.itemNoteFontSize+g.opts.lineSpacing > g.opts.pageHeight-g.opts.marginBottom {
			g.doc.AddPage()
			g.pageCount++
			y = g.opts.marginTop
		}

		g.doc.SetX(g.opts.marginLeft)
		g.doc.SetY(y)
		g.doc.Text(l)
		y += g.opts.itemNoteFontSize + g.opts.itemNoteMargin
	}

	return y, nil
}

func (g *pdfRenderImpl) savePDF() error {
	err := g.doc.SetFont(fontName, "", 8.0)
	if err != nil {
		return err
	}

	for i := range g.pageCount {
		g.doc.SetPage(i + 1)
		pageText := fmt.Sprintf("%d / %d", i+1, g.pageCount)
		pageNumWidth, err := g.doc.MeasureTextWidth(pageText)
		if err != nil {
			return err
		}

		g.doc.SetX(g.opts.pageWidth - g.opts.marginLeft - pageNumWidth)
		g.doc.SetY(g.opts.pageHeight - g.opts.marginBottom + 8) // un po' sopra il bordo inferiore
		g.doc.Text(pageText)
	}

	return g.doc.WritePdf(filepath.Join(g.outputDir, g.filename))
}

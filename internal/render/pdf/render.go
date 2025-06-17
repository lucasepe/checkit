package pdf

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lucasepe/checkit/internal/parser"
	"github.com/lucasepe/checkit/internal/render"
	"github.com/lucasepe/x/text/slugify"
	"github.com/mattn/go-runewidth"
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
	err := g.doc.AddTTFFontData(defaultFontName, gomono.TTF)
	if err != nil {
		return g, err
	}

	g.doc.AddPage()
	g.pageCount = 1

	return g, nil
}

const (
	defaultFontName         = "GoMono"
	defaultSymbol           = " \u25CB " // ○
	defaultLineSpacingRatio = 0.5
	defaultCreator          = "Check IT (github.com/lucasepe/checkit)"
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

		for _, el := range grp.Notes {
			y, err = g.handleGroupNote(y, el)
			if err != nil {
				return err
			}
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
	g.opts.itemNoteFontSize = 0.75 * fontSize
	g.opts.groupTitleFontSize = fontSize + 4
	g.opts.groupNoteFontSize = fontSize + 2
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

func (g *pdfRenderImpl) handleDocumentTitle(y float64, title string) (float64, error) {
	g.filename = fmt.Sprintf("%s.pdf", slugify.Sprint(title))

	prefix := ""
	y += defaultLineSpacingRatio * g.opts.documentTitleFontSize
	return g.renderText(y, title, prefix, g.opts.documentTitleFontSize, true, 53, 57, 53)
}

func (g *pdfRenderImpl) handleGroup(y float64, title string) (float64, error) {
	prefix := ""
	y += defaultLineSpacingRatio * g.opts.groupTitleFontSize
	return g.renderText(y, title, prefix, g.opts.groupTitleFontSize, false, 53, 57, 53)
}

func (g *pdfRenderImpl) handleGroupNote(y float64, line string) (float64, error) {
	prefix := " "
	return g.renderText(y, line, prefix, g.opts.groupNoteFontSize, false, 178, 190, 181)
}

func (g *pdfRenderImpl) handleItem(y float64, line string) (float64, error) {
	prefix := defaultSymbol

	return g.renderText(y, line, prefix, g.opts.itemFontSize, false, 54, 69, 79)
}

func (g *pdfRenderImpl) handleItemNote(y float64, line string) (float64, error) {
	prefix := strings.Repeat(" ", 5)
	return g.renderText(y, line, prefix, g.opts.itemNoteFontSize, false, 132, 136, 139)
}

func (rdr *pdfRenderImpl) renderText(y float64, line string, firstPrefix string, fontSize float64, center bool, r, g, b uint8) (float64, error) {
	err := rdr.doc.SetFont(defaultFontName, "", fontSize)
	if err != nil {
		return y, err
	}

	rdr.doc.SetTextColor(r, g, b)

	indent := strings.Repeat(" ", runewidth.StringWidth(firstPrefix))
	maxTextWidth := rdr.opts.pageWidth - 2*rdr.opts.marginLeft
	wrappedLines := wrapTextWithPrefix(rdr.doc, line, firstPrefix, indent, maxTextWidth)

	deltaY := defaultLineSpacingRatio * fontSize

	for _, l := range wrappedLines {

		if y+fontSize+2*deltaY > rdr.opts.pageHeight-rdr.opts.marginBottom {
			rdr.doc.AddPage()
			rdr.pageCount++
			y = rdr.opts.marginTop
		}

		x := rdr.opts.marginLeft
		if center {
			tw, err := rdr.doc.MeasureTextWidth(l)
			if err != nil {
				return y, err
			}

			x = ((rdr.opts.pageWidth - rdr.opts.marginLeft) - tw) / 2
		}
		rdr.doc.SetX(x)
		rdr.doc.SetY(y + deltaY)
		rdr.doc.Text(l)
		y += (fontSize + 2*deltaY)
	}

	return y, nil
}

func (g *pdfRenderImpl) savePDF() error {
	err := g.doc.SetFont(defaultFontName, "", 8.0)
	if err != nil {
		return err
	}

	g.doc.SetTextColor(54, 69, 79)

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

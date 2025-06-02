package render

import (
	"errors"
	"fmt"
	"os"

	pdfrender "github.com/lucasepe/checkit/internal/render/pdf"
	getoptutil "github.com/lucasepe/checkit/internal/util/getopt"
	ioutil "github.com/lucasepe/checkit/internal/util/io"
	"github.com/lucasepe/x/getopt"
	"github.com/lucasepe/x/text/conv"
)

func Do(args []string) error {
	extras, opts, err := getopt.GetOpt(args,
		"s:o:",
		[]string{"font-size", "output"},
	)
	if err != nil {
		return err
	}

	fontSize := conv.Float64(getoptutil.FindOptVal(opts, []string{"-s", "--font-size"}), 12.0)
	output := getoptutil.FindOptVal(opts, []string{"-o", "--output-dir"})

	var filename string
	if len(extras) > 0 {
		filename = extras[0]
	}

	src, cleanup, err := ioutil.FileOrStdin(filename)
	if err != nil {
		if errors.Is(err, ioutil.ErrNoInputDetected) {
			fmt.Fprintf(os.Stderr, "warning: %s.\n", err.Error())
		}
		return err
	}
	defer cleanup()

	render, err := pdfrender.New(
		pdfrender.FontSize(fontSize),
		pdfrender.OutputDir(output),
	)
	if err != nil {
		return err
	}

	return render.Render(src)
}

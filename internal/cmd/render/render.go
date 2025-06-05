package render

import (
	"errors"
	"fmt"
	"os"

	"github.com/lucasepe/checkit/internal/parser"
	pdfrender "github.com/lucasepe/checkit/internal/render/pdf"
	getoptutil "github.com/lucasepe/checkit/internal/util/getopt"
	ioutil "github.com/lucasepe/checkit/internal/util/io"
	"github.com/lucasepe/x/getopt"
)

func Do(args []string) error {
	extras, opts, err := getopt.GetOpt(args,
		"o:s",
		[]string{"output", "square"},
	)
	if err != nil {
		return err
	}

	square := getoptutil.HasOpt(opts, []string{"-s", "--square"})
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

	lst, err := parser.Parse(src)
	if err != nil {
		return err
	}

	render, err := pdfrender.New(
		pdfrender.Square(square),
		pdfrender.OutputDir(output),
	)
	if err != nil {
		return err
	}

	return render.Render(&lst)
}

package cmd

import (
	"fmt"
	"io"
	"strings"

	xtext "github.com/lucasepe/x/text"
)

const (
	appName = "checkit"
)

func usage(wri io.Writer) {
	var (
		desc = []string{
			"Render your checklists into clean printable PDF documents.\n",
			"Ideal for quick task lists, audit checklists, packing lists, etc.",
		}

		donateInfo = []string{
			"If you find this tool helpful consider supporting with a donation.",
			"Every bit helps cover development time and fuels future improvements.\n",
			"Your support truly makes a difference — thank you!\n",
			"  * https://www.paypal.com/donate/?hosted_button_id=FV575PVWGXZBY\n",
		}
	)

	fmt.Fprintln(wri)
	fmt.Fprint(wri, "┌─┐┬ ┬┌─┐┌─┐┬┌─\n")
	fmt.Fprint(wri, "│  ├─┤├┤ │  ├┴┐\n")
	fmt.Fprint(wri, "└─┘┴ ┴└─┘└─┘┴ ┴ IT\n")

	fmt.Fprintln(wri)
	for _, el := range desc {
		if el[0] == 194 {
			fmt.Fprintf(wri, "%s\n\n", xtext.Indent(xtext.Wrap(el, 60), "  "))
			continue
		}
		fmt.Fprintf(wri, "%s\n\n", xtext.Wrap(el, 76))
	}
	fmt.Fprintln(wri)

	fmt.Fprint(wri, "USAGE:\n\n")
	fmt.Fprintf(wri, "  %s [FLAGS] [INPUT_FILE]\n\n", appName)

	fmt.Fprint(wri, "FLAGS:\n\n")
	fmt.Fprint(wri, "  -s, --font-size     Base font size in points (default: 12).\n")
	fmt.Fprint(wri, "  -o, --output-dir    Output directory for generated PDF.\n")
	fmt.Fprint(wri, "  -h, --help          Show help and exit.\n")
	fmt.Fprint(wri, "  -v, --version       Show version and exit.\n")
	fmt.Fprintln(wri)

	fmt.Fprint(wri, "EXAMPLES:\n\n")
	fmt.Fprint(wri, " » Generate a PDF checklist from a file:\n\n")
	fmt.Fprintf(wri, "     %s /path/to/my-checklist.md\n\n", appName)
	fmt.Fprint(wri, " » Pipe input from another command:\n\n")
	fmt.Fprintf(wri, "     cat /path/to/my-checklist.md | %s\n\n", appName)

	fmt.Fprint(wri, "SUPPORT:\n\n")
	fmt.Fprint(wri, xtext.Indent(strings.Join(donateInfo, "\n"), "  "))
	fmt.Fprint(wri, "\n\n")

	fmt.Fprintln(wri, "Copyright (c) 2025 Luca Sepe")
}

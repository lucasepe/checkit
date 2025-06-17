package pdf

import (
	"strings"

	"github.com/signintech/gopdf"
)

// wrapTextWithPrefix spezza il testo tenendo conto di un prefisso sulla prima riga e di indentazioni successive
func wrapTextWithPrefix(doc *gopdf.GoPdf, text, firstPrefix, indent string, maxWidth float64) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string
	prefix := firstPrefix

	for _, word := range words {
		testLine := strings.TrimSpace(currentLine + " " + word)
		width, _ := doc.MeasureTextWidth(prefix + testLine)

		if width > maxWidth && currentLine != "" {
			lines = append(lines, prefix+currentLine)
			currentLine = word
			prefix = indent
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, prefix+currentLine)
	}
	return lines
}

func FitText(text string, maxWidth float64, measurer TextMeasurer) (int, bool, error) {
	runes := []rune(text)
	low, high := 0, len(runes)
	var fitCount int

	for low <= high {
		mid := (low + high) / 2
		substr := string(runes[:mid])
		width, _, err := measurer.Measure(substr)
		if err != nil {
			return 0, false, err
		}

		if width <= maxWidth {
			fitCount = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	//fitText := string(runes[:fitCount])
	fitsAll := fitCount == len(runes)
	return fitCount, fitsAll, nil
}

type TextMeasurer interface {
	Measure(text string) (width float64, height float64, err error)
}

type Margins struct {
	Left, Right, Top, Bottom float64
}

func WrapText(text, firstPrefix, indent string, maxWidth float64, measurer TextMeasurer) ([]string, error) {
	words := strings.Fields(text)
	var lines []string
	var currentLine string
	prefix := firstPrefix

	for _, word := range words {
		testLine := strings.TrimSpace(currentLine + " " + word)
		width, _, err := measurer.Measure(prefix + testLine)
		if err != nil {
			return nil, err
		}

		if width > maxWidth && currentLine != "" {
			lines = append(lines, prefix+currentLine)
			currentLine = word
			prefix = indent
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, prefix+currentLine)
	}
	return lines, nil
}

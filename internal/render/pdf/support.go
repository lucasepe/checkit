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

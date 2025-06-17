package pdf

import (
	"testing"
)

// mockMeasurer è un Measurer fittizio che assegna una larghezza fissa a ogni carattere (runa)
type mockMeasurer struct {
	perCharWidth float64
}

func (m mockMeasurer) Measure(text string) (float64, float64, error) {
	runes := []rune(text)
	return float64(len(runes)) * m.perCharWidth, 10.0, nil // altezza fissa 10.0
}

func TestFitText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		maxWidth     float64
		perCharWidth float64
		wantCount    int
		wantFitsAll  bool
	}{
		{
			name:         "short text fits all",
			text:         "Hello",
			maxWidth:     100,
			perCharWidth: 10,
			wantCount:    5,
			wantFitsAll:  true,
		},
		{
			name:         "text trimmed at 3 characters",
			text:         "Hello",
			maxWidth:     30,
			perCharWidth: 10,
			wantCount:    3,
			wantFitsAll:  false,
		},
		{
			name:         "empty text",
			text:         "",
			maxWidth:     50,
			perCharWidth: 10,
			wantCount:    0,
			wantFitsAll:  true,
		},
		{
			name:         "zero width box",
			text:         "Test",
			maxWidth:     0,
			perCharWidth: 10,
			wantCount:    0,
			wantFitsAll:  false,
		},
		{
			name:         "exact fit",
			text:         "Golang",
			maxWidth:     60,
			perCharWidth: 10,
			wantCount:    6,
			wantFitsAll:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			measurer := mockMeasurer{perCharWidth: tt.perCharWidth}
			count, fitText, err := FitText(tt.text, tt.maxWidth, measurer)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if count != tt.wantCount {
				t.Errorf("fitCount = %d, want %d", count, tt.wantCount)
			}
			if fitText != tt.wantFitsAll {
				t.Errorf("fitText = %t, want %t", fitText, tt.wantFitsAll)
			}
		})
	}
}

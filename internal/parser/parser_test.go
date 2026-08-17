package parser

import (
	charreader "calculator/internal/char-reader"
	"testing"
)

func TestNumbers(t *testing.T) {
	tests := []struct {
		name string
		text string
		want float64
	}{
		{
			name: "positive integer",
			text: "1337",
			want: 1337,
		},
		{
			name: "negative integer",
			text: "-42",
			want: -42,
		},
		{
			name: "float",
			text: "3.1415926",
			want: 3.1415926,
		},
		{
			name: "float with E",
			text: "3.2e26",
			want: 3.2e26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(charreader.New())
			got, err := p.Parse(tt.text)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Parse got %v, want %v", got, tt.want)
			}
		})
	}
}

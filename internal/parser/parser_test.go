package parser

import (
	charreader "calculator/internal/char-reader"
	"testing"

	"github.com/stretchr/testify/require"
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
			name: "unary operator",
			text: "-----10",
			want: -10,
		},
		{
			name: "float",
			text: "3.1415926",
			want: 3.1415926,
		},
		{
			name: "strange zero",
			text: "-0.00",
			want: 0,
		},
		//{
		//	name: "float with E",
		//	text: "3.2e26",
		//	want: 3.2e26,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(charreader.New())
			got, err := p.Parse(tt.text)
			require.NoError(t, err, "Unexpected error: %v", err)
			require.Equal(t, got, tt.want, "Parse got %v, want %v", got, tt.want)
		})
	}
}

func TestIncorrectNumbers(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr string
	}{
		{
			name:    "wrong base",
			text:    "0x13",
			wantErr: "Expected: +-*/^, got x",
		},
		{
			name:    "two dots",
			text:    "-42.3.2",
			wantErr: "invalid syntax",
		},
		{
			name:    "letter",
			text:    "A",
			wantErr: "Unknown char: A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(charreader.New())
			_, err := p.Parse(tt.text)
			require.NotNil(t, err, "expected error")
			require.ErrorContains(t, err, tt.wantErr, "Parse error got %v, want %v", err.Error(), tt.wantErr)
		})
	}
}

func TestBrackets(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    float64
		wantErr string
	}{
		{
			name: "single brackets",
			text: "(1337)",
			want: 1337,
		},
		{
			name: "multiply brackets",
			text: "(((42)))",
			want: 42,
		},
		{
			name: "brackets with negative",
			text: "(-42)",
			want: -42,
		},
		{
			name: "brackets with unary minus",
			text: "(-((-(-10))))",
			want: -10,
		},
		{
			name:    "no closed brackets",
			text:    "(1337",
			wantErr: "Expected: ), got EOF",
		},
		{
			name:    "no open brackets",
			text:    "1337)",
			wantErr: "Expected EOF, got )",
		},
		{
			name:    "missplaced bracket",
			text:    "133)7",
			wantErr: "Expected EOF, got )",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(charreader.New())
			got, err := p.Parse(tt.text)
			if err != nil {
				if tt.wantErr == "" {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				require.ErrorContains(t, err, tt.wantErr, "Parse error got %v, want %v", err.Error(), tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got, "Parse got %v, want %v", got, tt.want)
		})
	}
}

func TestBinaryOperators(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    float64
		wantErr string
	}{
		{
			name: "add",
			text: "1+2",
			want: 3,
		},
		{
			name: "sub",
			text: "1-2",
			want: -1,
		},
		{
			name: "mul",
			text: "1*2",
			want: 2,
		},
		{
			name: "div",
			text: "10/2",
			want: 5,
		},
		{
			name: "addsub",
			text: "10-2+3+6-6",
			want: 11,
		},
		{
			name: "muldiv",
			text: "10/2*3",
			want: 15,
		},
		{
			name: "big expression",
			text: "45-56*9*(-3)+638/(-76)*13*(625+9-(8267+654)-786.262*67)",
			want: 6654933.301,
		},
		{
			name: "unary operator add",
			text: "20+-1+30---10",
			want: 39,
		},
		{
			name: "unary operator with div",
			text: "21/-3",
			want: -7,
		},
		{
			name:    "division by zero",
			text:    "37*9+(10/0)-13*2",
			wantErr: "division by zero",
		},
		{
			name:    "two operators",
			text:    "10+*2",
			wantErr: "Unknown char: *",
		},
		{
			name:    "strange operator",
			text:    "10@2",
			wantErr: "Expected: +-*/^, got @",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(charreader.New())
			got, err := p.Parse(tt.text)
			if err != nil {
				if tt.wantErr == "" {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				require.ErrorContains(t, err, tt.wantErr, "Parse error got %v, want %v", err.Error(), tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got, "Parse got %v, want %v", got, tt.want)
		})
	}
}

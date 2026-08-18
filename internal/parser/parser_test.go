package parser

import (
	charreader "calculator/internal/char-reader"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name    string
	text    string
	want    float64
	wantErr string
}

func RunTests(t *testing.T, testCases []testCase) {
	t.Helper()
	for _, tt := range testCases {
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
			require.LessOrEqual(t, math.Abs(tt.want-got), 1e-9, "Parse got %v, want %v", got, tt.want)
		})
	}
}

func TestNumbers(t *testing.T) {
	tests := []testCase{
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
		//{
		//	name: "float with E",
		//	text: "3.2e26",
		//	want: 3.2e26,
		//},
	}
	RunTests(t, tests)
}

func TestBrackets(t *testing.T) {
	tests := []testCase{
		{
			name: "single brackets",
			text: "(1337)",
			want: 1337,
		},
		{
			name: "multiple brackets",
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
	RunTests(t, tests)
}

func TestBinaryOperators(t *testing.T) {
	tests := []testCase{
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
	RunTests(t, tests)
}

func TestPower(t *testing.T) {
	tests := []testCase{
		{
			name: "simple",
			text: "3^2",
			want: 9,
		},
		{
			name: "negative base",
			text: "(-2)^5",
			want: -32,
		},
		{
			name: "negative power",
			text: "2^-2",
			want: 0.25,
		},
		{
			name: "zero",
			text: "0^1000",
			want: 0,
		},
		{
			name: "one",
			text: "1^1000000000",
			want: 1,
		},
		{
			name: "float power",
			text: "9^(1/2)",
			want: 3,
		},
		{
			name: "multiple power",
			text: "4^2^3",
			want: 65536,
		},
		{
			name: "big expression",
			text: "5^(-15)^3*18*(-5+37.89)/2.737^(-9)+26^0-39",
			want: -38,
		},
		{
			name:    "negative root",
			text:    "-2^(1/2)",
			wantErr: "negative",
		},
		{
			name:    "zero with negative power",
			text:    "(1+2-3)^(-4)",
			wantErr: "negative power",
		},
	}
	RunTests(t, tests)
}

func TestConstants(t *testing.T) {
	tests := []testCase{
		{
			name: "pi",
			text: "pi",
			want: math.Pi,
		},
		{
			name: "e",
			text: "e",
			want: math.E,
		},
		{
			name: "phi",
			text: "phi",
			want: math.Phi,
		},
		{
			name: "big expression",
			text: "2*pi-11*e+(1/2*phi-pi^2)",
			want: -32.67850221258432,
		},
		{
			name:    "strange constant",
			text:    "phga",
			wantErr: "unknown letters",
		},
	}
	RunTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []testCase{
		{
			name: "log2",
			text: "log2(8)",
			want: 3,
		},
		{
			name: "log10",
			text: "log10(200)",
			want: 2.301029995663981,
		},
		{
			name: "ln",
			text: "ln(e^10)",
			want: 10,
		},
		{
			name: "exp",
			text: "exp(2)",
			want: 7.38905609893065,
		},
		{
			name: "sqrt",
			text: "sqrt(34^4)",
			want: 1156,
		},
		{
			name: "cbrt",
			text: "cbrt(27^2)",
			want: 9,
		},
		{
			name: "abs",
			text: "abs(-12)",
			want: 12,
		},
		{
			name: "big expression",
			text: "ln(exp(1)*log2(log10(100))+sqrt(2*pi)-cbrt(sqrt(64)))",
			want: 1.1709050748483818,
		},
		{
			name:    "strange function",
			text:    "sas(27^2)",
			wantErr: "unknown",
		},
	}
	RunTests(t, tests)
}

func TestTrigonometry(t *testing.T) {
	tests := []testCase{
		{
			name: "sin",
			text: "sin(-pi)",
			want: 0,
		},
		{
			name: "cos",
			text: "cos(3*pi/4)",
			want: -math.Sqrt(2) / 2,
		},
		{
			name: "tan",
			text: "tan(pi)",
			want: 0,
		},
		{
			name: "asin",
			text: "asin(sin(0.67+2*pi))",
			want: 0.67,
		},
		{
			name: "acos",
			text: "acos(0)",
			want: math.Pi / 2,
		},
		{
			name: "atan",
			text: "atan(3)",
			want: 1.24904577239,
		},
		{
			name: "big expression",
			text: "sin(2*pi)+(cos(4)^2+sin(4)^2)*tan(acos(sqrt(2)/2))",
			want: 1,
		},
	}
	RunTests(t, tests)
}

func TestFactorial(t *testing.T) {
	tests := []testCase{
		{
			name: "factorial",
			text: "6!",
			want: 720,
		},
		{
			name: "big expression",
			text: "-3!^2*1!",
			want: 36,
		},
		{
			name:    "negative factorial",
			text:    "(-2)!",
			wantErr: "negative",
		},
		{
			name:    "float factorial",
			text:    "(1/2)!",
			wantErr: "float",
		},
	}
	RunTests(t, tests)
}

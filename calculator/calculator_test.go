package calculator_test

import (
	"fmt"
	"github.com/denismitr/gds/calculator"
	"testing"
)

func TestCalculate(t *testing.T) {
	tt := []struct{
		input string
		output string
	}{
		{input: "2 + 2", output: "4"},
		{input: "2 + 2 - 1", output: "3"},
		{input: "2 + 2 - 1", output: "3"},
		{input: "9 + ( ( 24 / 6 ) - 2 )", output: "11"},
		{input: "9 + 24 / ( 7 - 3 )", output: "15"},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%s -> %s", tc.input, tc.output), func(t *testing.T) {
			output, err := calculator.Calculate(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if tc.output != output {
				t.Errorf("expected %s, got %s", tc.output, output)
			}
		})
	}
}

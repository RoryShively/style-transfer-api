// +build unit

package main

import "testing"


// file testing = https://www.youtube.com/watch?v=S1O0XI0scOM 

type CalcTest struct {
	input int
	output int
}

func TestParameterizedExample(t *testing.T) {
	var tests = []CalcTest{
		{2, 4},
		{3, 9},
		{4, 16},
		{-5, 25},
	}

	for _, test := range tests {
		if result := ContrivedCalculator(test.input); result != test.output {
			t.Errorf("Test Failed: %v input, %v output, result: %v",
				test.input,
				test.output,
				result,
			)
		}
	}
}


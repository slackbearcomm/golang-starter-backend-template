package services

import "testing"

type Result struct {
	x        int64
	y        int64
	expected int64
}

var addResults = []Result{
	{1, 1, 2},
	{2, 2, 4},
	{3, 4, 7},
}

func TestAddition(t *testing.T) {
	for _, test := range addResults {
		result := Addition(test.x, test.y)
		if result != test.expected {
			t.Fatal("Addition: output is not expected result")
		}
	}
}

var subtractResults = []Result{
	{1, 1, 0},
	{2, 2, 0},
	{3, 4, -1},
	{7, 4, 3},
}

func TestSubtraction(t *testing.T) {
	for _, test := range subtractResults {
		result := Subtraction(test.x, test.y)
		if result != test.expected {
			t.Fatal("Subtraction: output is not expected result")
		}
	}
}

var multiplyResults = []Result{
	{1, 1, 1},
	{2, 2, 4},
	{3, 4, 12},
	{7, 4, 28},
}

func TestMultiplication(t *testing.T) {
	for _, test := range multiplyResults {
		result := Multiplication(test.x, test.y)
		if result != test.expected {
			t.Fatal("Multiplication: output is not expected result")
		}
	}
}

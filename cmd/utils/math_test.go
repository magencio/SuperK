package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMath_Min(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{"Left", 3, 6, 3},
		{"Right", 6, 3, 3},
		{"Equal", 6, 6, 6},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			result := Min(test.x, test.y)

			// Assert
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMath_Max(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{"Right", 3, 6, 6},
		{"Left", 6, 3, 6},
		{"Equal", 6, 6, 6},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			result := Max(test.x, test.y)

			// Assert
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMath_Mod(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{"x bigger than y and result 0", 10, 5, 0},
		{"x bigger than y and result non-0", 3, 2, 1},
		{"x smaller than y and result non-0", 3, 6, 3},
		{"x smaller than y and result 0", 0, 6, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			result := Mod(test.x, test.y)

			// Assert
			assert.Equal(t, test.expected, result)
		})
	}
}

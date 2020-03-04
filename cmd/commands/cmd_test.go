package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_New(t *testing.T) {
	//Arrange
	expectedArgs := []string{"kubectl", "-n", "kubeflow", "get", "run"}

	// Act
	command := NewCmd("kubectl", "-n", "kubeflow", "get", "run")

	// Assert
	assert.Equal(t, "/usr/local/bin/kubectl", command.Path)
	assert.EqualValues(t, expectedArgs, command.Args)
}

func TestCommand_ToString(t *testing.T) {
	//Arrange
	command := NewCmd("kubectl", "-n", "kubeflow", "get", "run")
	expected := "kubectl -n kubeflow get run"

	// Act
	result := command.ToString()

	// Assert
	assert.Equal(t, expected, result)
}

func TestCommand_Run_NoCacheFirst(t *testing.T) {
	//Arrange
	command := NewCmd("printf", "%s", "This is a test")
	expectedOutput := "This is a test"
	result1 := command.Run(false)
	time1 := *result1.RunTime
	assert.Equal(t, expectedOutput, *result1.Output)

	// Act
	result2 := command.Run(false)

	// Assert
	assert.Equal(t, expectedOutput, *result2.Output)
	assert.True(t, result2.RunTime.After(time1))
}

func TestCommand_Run_CacheFirst(t *testing.T) {
	//Arrange
	command := NewCmd("printf", "%s", "This is a test")
	expectedOutput := "This is a test"
	result1 := command.Run(true)
	time1 := *result1.RunTime
	assert.Equal(t, expectedOutput, *result1.Output)

	// Act
	result2 := command.Run(true)

	// Assert
	assert.Equal(t, expectedOutput, *result2.Output)
	assert.Equal(t, time1, *result2.RunTime)
}

func TestCommand_RunInvalid_NoCacheFirst(t *testing.T) {
	//Arrange
	command := NewCmd("printf", "%d", "This is a test")
	expectedOutput := "printf: 'This is a test': expected a numeric value\n0"
	result1 := command.Run(false)
	time1 := *result1.RunTime
	assert.Equal(t, expectedOutput, *result1.Output)

	// Act
	result2 := command.Run(false)

	// Assert
	assert.Equal(t, expectedOutput, *result2.Output)
	assert.True(t, result2.RunTime.After(time1))
}

func TestCommand_RunInvalid_CacheFirst(t *testing.T) {
	//Arrange
	command := NewCmd("printf", "%d", "This is a test")
	expectedOutput := "printf: 'This is a test': expected a numeric value\n0"
	result1 := command.Run(true)
	time1 := *result1.RunTime
	assert.Equal(t, expectedOutput, *result1.Output)

	// Act
	result2 := command.Run(true)

	// Assert
	assert.Equal(t, expectedOutput, *result2.Output)
	assert.Equal(t, time1, *result2.RunTime)
}

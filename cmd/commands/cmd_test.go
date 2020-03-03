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

func TestCommand_Run(t *testing.T) {
	//Arrange
	command := NewCmd("printf", "%s", "This is a test")
	expected := "This is a test"

	// Act
	result := command.Run()

	// Assert
	assert.Equal(t, expected, *result)
}

func TestCommand_RunInvalid(t *testing.T) {
	//Arrange
	command := NewCmd("printf", "%d", "This is a test")
	expected := "printf: 'This is a test': expected a numeric value\n0"

	// Act
	result := command.Run()

	// Assert
	assert.Equal(t, expected, *result)
}

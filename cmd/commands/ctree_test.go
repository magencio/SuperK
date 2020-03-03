package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCTree_NewNil(t *testing.T) {
	// Arrange

	// Act
	tree, err := NewCTree(nil)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, tree.Part, "kubectl")
	assert.Len(t, tree.Children, 0)
	assert.Nil(t, tree.Parent)
}

func TestCTree_NewEmpty(t *testing.T) {
	// Arrange

	// Act
	tree, err := NewCTree([]string{})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, tree.Part, "kubectl")
	assert.Len(t, tree.Children, 0)
	assert.Nil(t, tree.Parent)
}

func TestCTree_New(t *testing.T) {
	// Arrange
	commands := []string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	}

	// Act
	tree, err := NewCTree(commands)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, tree.Part, "kubectl")
	assert.Len(t, tree.Children, 2)
	assert.Equal(t, tree.Children[0].Part, "-n kubeflow")
	assert.Len(t, tree.Children[0].Children, 1)
	assert.Equal(t, tree.Children[0].Children[0].Part, "get")
	assert.Len(t, tree.Children[0].Children[0].Children, 2)
	assert.Equal(t, tree.Children[0].Children[0].Children[0].Part, "pod")
	assert.Equal(t, tree.Children[0].Children[0].Children[1].Part, "cronjob")
	assert.Equal(t, tree.Children[1].Part, "-n pipelines")
	assert.Len(t, tree.Children[1].Children, 1)
	assert.Equal(t, tree.Children[1].Children[0].Part, "get")
	assert.Len(t, tree.Children[1].Children[0].Children, 1)
	assert.Equal(t, tree.Children[1].Children[0].Children[0].Part, "pod")
}

func TestCTree_NewInvalid(t *testing.T) {
	// Arrange
	commands := []string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
		"az login",
	}

	// Act
	tree, err := NewCTree(commands)

	// Assert
	assert.Nil(t, tree)
	assert.NotNil(t, err)
}

func TestCTree_MergeCommand(t *testing.T) {
	// Arrange
	tree, err := NewCTree(nil)
	assert.Nil(t, err)

	// Act
	err = tree.MergeCommand("kubectl -n kubeflow get pod")
	assert.Nil(t, err)

	err = tree.MergeCommand("kubectl -n kubeflow get cronjob")
	assert.Nil(t, err)

	err = tree.MergeCommand("kubectl -n pipelines get pod")
	assert.Nil(t, err)

	// Assert
	assert.Equal(t, tree.Part, "kubectl")
	assert.Len(t, tree.Children, 2)
	assert.Equal(t, tree.Children[0].Part, "-n kubeflow")
	assert.Len(t, tree.Children[0].Children, 1)
	assert.Equal(t, tree.Children[0].Children[0].Part, "get")
	assert.Len(t, tree.Children[0].Children[0].Children, 2)
	assert.Equal(t, tree.Children[0].Children[0].Children[0].Part, "pod")
	assert.Equal(t, tree.Children[0].Children[0].Children[1].Part, "cronjob")
	assert.Equal(t, tree.Children[1].Part, "-n pipelines")
	assert.Len(t, tree.Children[1].Children, 1)
	assert.Equal(t, tree.Children[1].Children[0].Part, "get")
	assert.Len(t, tree.Children[1].Children[0].Children, 1)
	assert.Equal(t, tree.Children[1].Children[0].Children[0].Part, "pod")
}

func TestCTree_MergeCommandUnique(t *testing.T) {
	// Arrange
	tree, err := NewCTree(nil)
	assert.Nil(t, err)

	// Act
	err = tree.MergeCommand("kubectl -n kubeflow get pod")
	assert.Nil(t, err)

	err = tree.MergeCommand("kubectl -n kubeflow get cronjob")
	assert.Nil(t, err)

	err = tree.MergeCommand("kubectl -n pipelines get pod")
	assert.Nil(t, err)

	err = tree.MergeCommand("kubectl -n kubeflow get cronjob")
	assert.Nil(t, err)

	// Assert
	assert.Equal(t, tree.Part, "kubectl")
	assert.Equal(t, tree.Children[0].Part, "-n kubeflow")
	assert.Equal(t, tree.Children[0].Children[0].Part, "get")
	assert.Len(t, tree.Children[0].Children[0].Children, 2)
	assert.Equal(t, tree.Children[0].Children[0].Children[0].Part, "pod")
	assert.Equal(t, tree.Children[0].Children[0].Children[1].Part, "cronjob")
	assert.Equal(t, tree.Children[1].Part, "-n pipelines")
	assert.Equal(t, tree.Children[1].Children[0].Part, "get")
	assert.Equal(t, tree.Children[1].Children[0].Children[0].Part, "pod")
}

func TestCTree_MergeCommandInvalid(t *testing.T) {
	// Arrange
	tree, err := NewCTree(nil)
	assert.Nil(t, err)

	// Act
	err = tree.MergeCommand("az login")

	// Assert
	assert.NotNil(t, err)
}

func TestCTree_GetCommand(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	position := 5
	expected := NewCmd("kubectl", "-n", "kubeflow", "get", "cronjob")

	// Act
	result := tree.GetCmd(position)

	// Assert
	assert.Equal(t, expected, result)
}

func TestCTree_GetCommandInvalid(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	position := 10

	// Act
	result := tree.GetCmd(position)

	// Assert
	assert.Nil(t, result)
}

func TestCTree_GetPosition(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Position 2", "kubectl -n kubeflow", 2},
		{"Position 5", "kubectl -n kubeflow get cronjob", 5},
		{"Position 9", "kubectl --help", 9},
	}

	commands := []string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
		"kubectl --help",
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree, err := NewCTree(commands)
			assert.Nil(t, err)

			// Act
			result := tree.GetPosition(test.input)

			// Assert
			assert.Equal(t, test.expected, *result)
		})
	}
}

func TestCTree_GetPositionInvalid(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	// Act
	result := tree.GetPosition("kubectl -n kubeflow get service")

	// Assert
	assert.Nil(t, result)
}

func TestCTree_GetNextParts(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	expected := []string{"-n kubeflow", "-n pipelines"}

	// Act
	result := tree.GetNextParts("kubectl -n")

	// Assert
	assert.Len(t, result, 2)
	assert.EqualValues(t, expected, result)
}

func TestCTree_GetNextPartsNoneExist(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	// Act
	result := tree.GetNextParts("kubectl -n kubeflow d")

	// Assert
	assert.Nil(t, result)
}

func TestCTree_RemoveCommand(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	position := 3

	// Act
	err = tree.RemoveCommand(position)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, tree.Children, 2)
	assert.Equal(t, "-n kubeflow", tree.Children[0].Part)
	assert.Len(t, tree.Children[0].Children, 0)
	assert.Equal(t, "-n pipelines", tree.Children[1].Part)
	assert.Len(t, tree.Children[1].Children, 1)
	assert.Equal(t, "get", tree.Children[1].Children[0].Part)
	assert.Len(t, tree.Children[1].Children[0].Children, 1)
	assert.Equal(t, "pod", tree.Children[1].Children[0].Children[0].Part)
}

func TestCTree_RemoveCommandInvalid(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	position := 10

	// Act
	err = tree.RemoveCommand(position)

	// Assert
	assert.NotNil(t, err)
}

func TestCTree_ToStrings(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	expected := []string{
		"kubectl",
		"  -n kubeflow",
		"    get",
		"      pod",
		"      cronjob",
		"  -n pipelines",
		"    get",
		"      pod",
	}

	// Act
	result := tree.ToStrings(2)

	// Assert
	assert.Len(t, result, len(expected))
	assert.EqualValues(t, expected, result)
}

func TestCTree_Serialize(t *testing.T) {
	// Arrange
	tree, err := NewCTree([]string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	})
	assert.Nil(t, err)

	expected := []string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	}

	// Act
	result := tree.Serialize()

	// Assert
	assert.Len(t, result, len(expected))
	assert.EqualValues(t, expected, result)
}

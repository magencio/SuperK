package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackup_New(t *testing.T) {
	// Arrange
	path, err := getTmpPath("superk_test_")
	assert.Nil(t, err)
	fileName := filepath.Base(path)

	// Act
	backup := NewBackup(fileName)

	// Assert
	assert.Equal(t, path, backup.TempFile)
}

func TestBackup_SetCommands(t *testing.T) {
	// Arrange
	path, err := getTmpPath("superk_test_")
	assert.Nil(t, err)
	fileName := filepath.Base(path)
	backup := NewBackup(fileName)

	expected := []string{
		"kubectl -n kubeflow get pod",
		"kubectl -n kubeflow get cronjob",
		"kubectl -n pipelines get pod",
	}

	tree, err := NewCTree(expected)
	assert.Nil(t, err)

	// Act
	err = backup.SetCommands(tree)

	// Assert
	assert.Nil(t, err)
	result, err := backup.Commands()
	assert.Nil(t, err)
	assert.EqualValues(t, expected, result.Serialize())

	// Cleanup
	err = backup.Delete()
	assert.Nil(t, err)
}

func getTmpPath(prefix string) (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}
	path := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	return path, nil
}

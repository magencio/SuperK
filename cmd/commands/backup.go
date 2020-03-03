package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Backup structure represents a temp file to backup commands
type Backup struct {
	TempFile string
}

// NewBackup creates a new backup structure
func NewBackup(tempFileName string) *Backup {
	tempFile := path.Join(os.TempDir(), tempFileName)
	return &Backup{TempFile: tempFile}
}

// SetCommands updates the backup file with the provided commands
func (backup *Backup) SetCommands(commands *CTree) error {
	commandsToBackup := commands.Serialize()
	return backup.create(commandsToBackup)
}

func (backup *Backup) create(commands []string) error {
	file, err := os.Create(backup.TempFile)
	if err != nil {
		return err
	}
	for _, command := range commands {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", command)); err != nil {
			return err
		}
	}
	return nil
}

// Commands returns the commands from the backup file
func (backup *Backup) Commands() (*CTree, error) {
	backupCommands, err := backup.get()
	if err != nil {
		backupCommands = nil
	}
	return NewCTree(backupCommands)
}

func (backup *Backup) get() ([]string, error) {
	bytes, err := ioutil.ReadFile(backup.TempFile)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	commands := strings.FieldsFunc(content, func(ch rune) bool {
		return ch == '\n'
	})
	return commands, nil
}

// Delete deletes the backup file
func (backup *Backup) Delete() error {
	return os.Remove(backup.TempFile)
}

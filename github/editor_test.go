package github

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

func TestEditor_openAndEdit_fileExist(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "editor-test")
	ioutil.WriteFile(tempFile.Name(), []byte("hello"), 0644)
	editor := Editor{
		Program: "memory",
		File:    tempFile.Name(),
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile.Name(), file)

			return nil
		},
	}

	content, err := editor.openAndEdit()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello", string(content))

	// file is removed after edit
	_, err = os.Stat(tempFile.Name())
	assert.T(t, os.IsNotExist(err))
}

func TestEditor_openAndEdit_fileNotExist(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "editor-test")
	tempFile := filepath.Join(tempDir, "PULLREQ")
	editor := Editor{
		Program: "memory",
		File:    tempFile,
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile, file)

			return ioutil.WriteFile(file, []byte("hello"), 0644)
		},
	}

	content, err := editor.openAndEdit()
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello", string(content))

	// file is removed after edit
	_, err = os.Stat(tempFile)
	assert.T(t, os.IsNotExist(err))
}

func TestEditor_EditTitleAndBody(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "editor-test")
	tempFile := filepath.Join(tempDir, "PULLREQ")
	editor := Editor{
		Program: "memory",
		File:    tempFile,
		openEditor: func(program string, file string) error {
			assert.Equal(t, "memory", program)
			assert.Equal(t, tempFile, file)

			message := `A title
A title continues

A body
A body continues
# comment
`
			return ioutil.WriteFile(file, []byte(message), 0644)
		},
	}

	title, body, err := editor.EditTitleAndBody()
	assert.Equal(t, nil, err)
	assert.Equal(t, "A title A title continues", title)
	assert.Equal(t, "A body\nA body continues", body)

	// file is removed after edit
	_, err = os.Stat(tempFile)
	assert.T(t, os.IsNotExist(err))
}

func TestReadTitleAndBody(t *testing.T) {
	message := `A title
A title continues

A body
A body continues
# comment
`
	r := strings.NewReader(message)
	reader := bufio.NewReader(r)
	title, body, err := readTitleAndBody(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "A title A title continues", title)
	assert.Equal(t, "A body\nA body continues", body)
}

func TestGetMessageFile(t *testing.T) {
	gitPullReqMsgFile, _ := getMessageFile("PULLREQ")
	assert.T(t, strings.Contains(gitPullReqMsgFile, "PULLREQ_EDITMSG"))
}

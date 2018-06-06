package operations

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

const authvar = "GITHUB_TOKEN"

// Verbose prints some extra info
var Verbose = false

var ignorePaths = map[string]bool{
	".git":            true,
	"CONTRIBUTING.md": true,
	"LICENSE":         true,
}

// Log prints the message to stderr
func Log(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// LogError prints the error and a message to stderr
func LogError(err error, format string, args ...interface{}) {
	a := append(args, err.Error())
	fmt.Fprintf(os.Stderr, format+": %s\n", a...)
}

// Debug prints a message if the Verbose flag is set
func Debug(format string, args ...interface{}) {
	if Verbose {
		Log(format, args...)
	}
}

// copyFiles copies all the files/directories into
// the destination folder.
func copyFiles(src, dst string) error {
	Debug("Copying %s into %s", src, dst)

	if err := makeDir(dst); err != nil {
		LogError(err, "Error creating destination folder")
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		LogError(err, "Error reading files of %s", src)
		return err
	}

	for _, entry := range files {
		name := entry.Name()
		if ignorePaths[name] {
			continue
		}

		if err := copy(path.Join(src, name), path.Join(dst, name)); err != nil {
			LogError(err, "Error copying %s", name)
			return err
		}
	}

	return nil
}

// copy copies the file/directory to a different directory
// Golang, why don't you have this in the standard library?
func copy(src, dst string) error {
	Debug("Copying %s into %s", src, dst)
	cmd := exec.Command("cp", "-r", src, dst)
	if Verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	return cmd.Run()
}

// makeDir creates the full path for the integration
func makeDir(dir string) error {
	Debug("Creating %s", dir)
	return exec.Command("mkdir", "-p", dir).Run()
}

func executeTemplate(tmpl *template.Template, data interface{}) string {
	buffer := new(bytes.Buffer)

	if err := tmpl.Execute(buffer, data); err != nil {
		panic(err)
	}

	return buffer.String()
}

// GetAuthToken returns the authentication token if found, if not, it exist the
// app.
func GetAuthToken() string {
	token := os.Getenv(authvar)
	if token == "" {
		Log("Please export $%s with a personal token and try again. Exiting", authvar)
		os.Exit(1)
	}
	return token
}

func fileExists(file string) (bool, error) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		LogError(err, "Error reading file")
		return false, err

	}

	return true, nil
}

package Utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/go-repoflow/pkg/git"
)

// ListOfRequiredBinaries
const (
	// Formatting      string = "gofumpt"
	Formatting      string = "gofmt"
	Imports         string = "goimports"
	lint            string = "golangci-lint"
	Gosecurity      string = "govulncheck"
	GoSecurity      string = "gosec"
	Pythonformatter string = "ruff"
	PythonTyping    string = "mypy"
	PythonTests     string = "pytest"
	Vulnerabilities string = "osv-scanner"
	Secrets         string = "gitleaks"
	Hooks           string = "pre-commit"
	GoCoverage      string = "go test"
	PyCoverage      string = "pytest-cov"
)

var ListOfRequiredBinaries = []string{Formatting, Imports, lint, Gosecurity, GoSecurity, Pythonformatter, PythonTyping, PythonTests, Vulnerabilities, Secrets, Hooks, GoCoverage, PyCoverage}

type GoProjectBinaries struct {
	Project  bool
	Binaries []string
}

type PythonProjectBinaries struct {
	Project  bool
	Binaries []string
}

func findBinaries(binary string, paths []string) bool {
	for _, path := range paths {

		if _, err := os.Stat(path); err != nil {
			continue
		}

		err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if d.Name() == binary {
				return fmt.Errorf("Found")
			}
			return nil
		})
		if err != nil {
			// fmt.Println("Found binary: ", binary)
			return true
		}
	}
	return false
}

func runCommand(command string, args []string) (string, error) {

	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	commandRunError := cmd.Run()
	if commandRunError != nil {
		fmt.Println("Error running ", command, "command", commandRunError.Error())
	}

	errBytes, err := io.ReadAll(&stderr)
	if err != nil {
		fmt.Println("Could not read error")
	}

	if len(errBytes) != 0 {
		return "", errors.New(string(errBytes))
	}

	return out.String(), nil
}

func toolAvalibilty(Binaries []string) {
	path := os.Getenv("PATH")
	if path == "" {
		fmt.Print(("No Path env found"))
	}

	paths := strings.Split(path, ":")

	// if len(paths) < len(ListOfRequiredBinaries) {
	for _, binary := range Binaries {
		found := findBinaries(binary, paths)

		if found {
			aphrodite.PrintInfo("Found binary for " + binary + "\n")
			continue
		} else {
			aphrodite.PrintError("No binary for " + binary + "\n")
			continue
		}
	}
}

func Doctor() {

	// }

	GoVersion, err := runCommand("go", []string{"version"})
	if err != nil {
		fmt.Println("Error getting go version: ", err)
	}

	PythonVersion, err := runCommand("python3", []string{"--version"})
	if err != nil {
		fmt.Println("Error getting go version: ", err)
	}

	GitVersion, err := runCommand("git", []string{"--version"})
	if err != nil {
		fmt.Println("Error getting go version: ", err)
	}

	fmt.Println("Go version: ", GoVersion)
	fmt.Println("Python version: ", PythonVersion)
	fmt.Println("Git version: ", GitVersion)

	toolAvalibilty(ListOfRequiredBinaries)

	if found := git.FindGitFolder(); found == false {
		fmt.Println("Could not find a git folder...")
	} else {
		fmt.Println("Currently in a git directory")
	}

}

package Utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

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

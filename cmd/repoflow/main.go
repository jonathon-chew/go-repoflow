package main

import (
	"fmt"
	"os"

	"github.com/jonathon-chew/go-repoflow/internal/cli"
	"github.com/jonathon-chew/go-repoflow/pkg/git"
)

func MAIN() int {

	// Check if there are arguments have been input - if so run through the cli module
	if len(os.Args[1:]) >= 1 {
		ErrProcessingcli := cli.CLI(os.Args[1:])
		if ErrProcessingcli != nil {

			// Print that there was an issue and the command passed in
			fmt.Printf("Error parsing the command line argument, %v\n", ErrProcessingcli)

			// Return with a bad status code to allow this to be checked in other programmes whether it was succesfully even understood!
			os.Exit(1)
		} else {
			// If there is no error exit the main function - this stops the deafult behaviour from writing to the file
			return 0
		}
	}

	// modify file is true if we want to modify the file, false if we just want to check for new todos
	return git.ProcessTodosInRepo(true)
}

func main() {
	os.Exit(MAIN())
}

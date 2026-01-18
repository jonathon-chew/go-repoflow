package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	utils "github.com/jonathon-chew/go-repoflow/internal/Utils"
	cmd "github.com/jonathon-chew/go-repoflow/internal/cli"
	"github.com/jonathon-chew/go-repoflow/pkg/git"
)

func main() {

	// Check if there are arguments have been input - if so run through the cmd module
	if len(os.Args[1:]) >= 1 {
		ErrProcessingCmd := cmd.CLI(os.Args[1:])
		if ErrProcessingCmd != nil {

			// Print that there was an issue and the command passed in
			fmt.Printf("Error parsing the command line argument, %v\n", ErrProcessingCmd)

			// Return with a bad status code to allow this to be checked in other programmes whether it was succesfully even understood!
			os.Exit(1)
		} else {
			// If there is no error exit the main function - this stops the deafult behaviour from writing to the file
			return
		}
	}
	// CHECK to see if their is a git folder
	// Initilaise the known files to ignore!
	unwantedFiles := []string{".localized", ".DS_Store", ".gitignore"}
	unwantedExtentions := []string{".app", ".exe", ".elf", ".md"}
	fileList := utils.FindFilesInCurrentDirectory()

	if !git.FindGitFolder() {
		os.Exit(1)
	}

	// Check there is an origin, and exit if not!
	_, remoteOriginErr := git.GetRemoteOrigin()
	if remoteOriginErr != nil {
		fmt.Printf("[ERROR]: %s\n", remoteOriginErr)
		os.Exit(1)
	}

	// Get a list of all current issues
	listOfGithubIssues, githubErr := git.ListGithubIssues(false)
	if githubErr != nil {
		if errors.Is(githubErr, fmt.Errorf("there were no github issues")) {
			fmt.Printf("[ERROR]: There was an error getting issues: %v\n", githubErr)
			return
		}
	}

	// Get the number of existing issues
	CurrentNumberOfIssues := len(listOfGithubIssues)

	var foundNewTODO bool = false
	for _, fileName := range fileList {

		// Keep going straight away if it's a directory
		if fileName.IsDir() || !strings.Contains(fileName.Name(), ".") {
			continue
		}

		// Get the lines of the file
		var fileLine []string

		// Set the file name
		var filePath = fileName.Name()

		// Make sure it's not one of the known unwanted files to edit
		if slices.Contains(unwantedFiles, filePath) {
			continue
		}

		// Set up variables to be used to check through eveyting that's already in place
		var unwantedExtention bool = false
		var updatedFile bool = false

		// ignore binary files!
		for _, extension := range unwantedExtentions {
			if strings.Contains(filePath, extension) {
				unwantedExtention = true
				break
			}
		}

		if unwantedExtention {
			continue
		}

		// Look for to dos in the file
		file, err := os.Open(filePath)
		if err != nil {
			return
		}

		var lineNumber int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()

			// This is adding a number to the start of the todo as a way to keep track and act as a guard against duplicating issues!
			if strings.Contains(line, "TODO: ") && !strings.Contains(line, ") TODO") {

				// Find the with the TODO in it
				var replaceString string = fmt.Sprintf("(#%d) TODO", CurrentNumberOfIssues+1)

				// Replace the issue with the replace string which now has a number in it
				line = strings.Replace(line, "TODO", replaceString, 1)

				// Print this to the screen
				fmt.Printf("I would like to make a github issue for: %s\nThe title is %s\nThe body is: %s on line %d\n", strings.TrimSpace(line), strings.TrimSpace(line), fileName.Name(), lineNumber)

				// Incriment the number of current issues - for the next time this needs to be used
				CurrentNumberOfIssues += 1

				// Check whether the issue already exists...
				git.MakeGithubIssue(line, fmt.Sprintf("This is from file %s on line %d\n", fileName.Name(), lineNumber))

				// Conditional if something has been updated, some actions needs to happen outside of the loop
				updatedFile, foundNewTODO = true, true

			} else if strings.Contains(line, "TODO: ") && strings.Contains(line, ") TODO") {
				// This finds OLD TODOs

				// (#22) TODO: If github issue not in the list of old todos close issue
				/* _, removeError := RemoveLineDueToGithubIssue(line, listOfGithubIssues)
				if removeError == nil {
					// (#21) TODO: If todo in the list of old todos and no longer open on github, remove line
					line = ""
				} */

				// issue here being TOO powerful, when run on itself it deletes the if statements! Check for number?
			}

			// Regardless of whether a line has changed or not, add it into the list of lines to write back in
			fileLine = append(fileLine, line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}

		// Write modified content back to the file
		if updatedFile {

			// Write the result of the parsing of the file to the file again!
			err = os.WriteFile(filePath, []byte(strings.Join(fileLine, "\n")), 0644)
			if err != nil {
				fmt.Println("Error writing file:", err)
				return
			}
		}
	}

	if !foundNewTODO {
		fmt.Println("No new todo found in any file in this directory")
	}
}

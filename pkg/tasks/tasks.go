package tasks

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/go-repoflow/pkg/git/git_utils"
)

type TaskState struct {
	Open   bool
	Closed bool
}

type Task struct {
	Number  string
	Details string
	Added   bool
	State   TaskState
}

const folder string = ".tasks"

/*
Wrapper around a function to check if file exists, don't care what the reason is,
just need a boolean whether it exists or not
*/
func fileExists() bool {
	FileInfo, err := os.Lstat(folder)
	if err != nil {
		return false
	}

	if !FileInfo.IsDir() {
		return false
	}

	return true
}

/*
Find whether the folder .tasks is if it is at all
Changes directory into the folder that contains the .tasks
*/
func findTasks() error {

	cwd := os.Getenv("PWD")
	if cwd == "" {
		return fmt.Errorf("Could not get PWD")
	}

	folders := strings.Split(cwd, string(os.PathSeparator))

	for i := len(folders); i > len(folders); i-- {
		os.Chdir(strings.Join(folders[:i], string(os.PathSeparator)))
		if fileExists() {
			return nil
		}
	}

	return fmt.Errorf("Could not find the tasks folder")
}

func findTODOs() {
	fileList := git_utils.FindAllFilesInCurrentDirectoryAndSubdirectories()
	unwantedFiles := []string{".localized", ".DS_Store", ".gitignore", "response.txt"}
	unwantedExtentions := []string{".app", ".exe", ".elf", ".md"}
	var modifyFile bool
	var foundNewTODO bool
	// Track all TODO lines found in files (for checking if issues should be closed)
	todosFoundInFiles := make(map[string]bool)

	path := os.Getenv("PWD")
	if path != "" {
		return
	}

	err := findTasks()
	if err != nil {
		return
	}

	NumberOfIssues, err := os.ReadDir(".")
	if err != nil {
		return
	}

	err = os.Chdir(path)
	if err != nil {
		return
	}

	var CurrentNumberOfIssues int = len(NumberOfIssues)

	for _, fileName := range fileList {
		// Keep going straight away if it's a directory
		if fileName.IsDir || !strings.Contains(fileName.Name, ".") {
			continue
		}

		// Get the lines of the file
		var fileLine []string

		// Set the file name
		filePath := fileName.FullPath

		// Make sure it's not one of the known unwanted files to edit
		if slices.Contains(unwantedFiles, fileName.Name) {
			continue
		}

		// Set up variables to be used to check through everything that's already in place
		var unwantedExtention bool
		var updatedFile bool

		// ignore binary files
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
			log.Println(aphrodite.ReturnError("Error opening file: " + err.Error()))
			continue
		}
		defer file.Close()

		var lineNumber int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()

			// Track this TODO line for later comparison with GitHub issues
			trimmedLine := strings.TrimSpace(line)
			todosFoundInFiles[trimmedLine] = true

			// This is adding a number to the start of the todo as a way to keep track and act as a guard against duplicating issues
			if strings.Contains(line, "TODO: ") && !strings.Contains(line, ") TODO") {

				// Find the with the TODO in it
				replaceString := fmt.Sprintf("(#%d) TODO", CurrentNumberOfIssues+1)

				// Replace the issue with the replace string which now has a number in it
				line = strings.Replace(line, "TODO", replaceString, 1)

				// Increment the number of current issues - for the next time this needs to be used
				CurrentNumberOfIssues++

				if modifyFile {
					// Print this to the screen
					fmt.Printf("I would like to make a github issue for: %s\nThe title is %s\nThe body is: %s on line %d\n", strings.TrimSpace(line), strings.TrimSpace(line), fileName.Name, lineNumber)

					// Check whether the issue already exists...
					// MakeGithubIssue(line, fmt.Sprintf("This is from file %s on line %d\n", fileName.Name, lineNumber), []string{"Bug"})

					// Conditional if something has been updated, some actions needs to happen outside of the loop
					updatedFile, foundNewTODO = true, true
				} else {
					fmt.Printf("The title of the todo is %s\nThe location is: %s on line %d\n\n", strings.TrimSpace(line), fileName.Name, lineNumber)
				}

			} else if strings.Contains(line, "TODO: ") && strings.Contains(line, ") TODO") {
				// Track this TODO line for later comparison with GitHub issues
				trimmedLine := strings.TrimSpace(line)
				todosFoundInFiles[trimmedLine] = true

				if modifyFile {
					// This finds OLD TODOs that have already had a GitHub issue created for them.
					// Try to close the corresponding GitHub issue and, if successful, remove the line.
					removed, removeError := RemoveLineDueToGithubIssue(line, listOfGithubIssues)
					if removeError != nil {
						log.Println(aphrodite.ReturnError("Error closing GitHub issue: " + removeError.Error()))
					}
					if removed && removeError == nil {
						// The issue was successfully closed; remove the TODO from the code.
						line = ""
						updatedFile = true
						// Remove from tracking since we're deleting it
						delete(todosFoundInFiles, trimmedLine)
					}

					if removed {
						fmt.Printf("I would like to remove the line for: %s\nThe title is %s\nThe body is: %s on line %d\n", strings.TrimSpace(line), strings.TrimSpace(line), fileName.Name, lineNumber)
					}
				} else {
					fmt.Printf("The title of the todo I am already tracking is in the file %s on line %d\n\n", fileName.Name, lineNumber)
				}

			}

			if modifyFile {
				// Regardless of whether a line has changed or not, add it into the list of lines to write back in
				fileLine = append(fileLine, line)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}

		// Write modified content back to the file
		if updatedFile && modifyFile {

			// Write the result of the parsing of the file to the file again
			err = os.WriteFile(filePath, []byte(strings.Join(fileLine, "\n")), 0644)
			if err != nil {
				fmt.Println("Error writing file:", err)
				return
			}
		}
	}

	// After scanning all files, check for GitHub issues that no longer have corresponding TODOs
	if modifyFile {
		if CurrentNumberOfIssues > 0 {
			for _, issue := range NumberOfIssues {

				issue.Info()
				// Only check open issues
				/* if issue.State == "open" {
					// log.Print("THIS ISSUE IS OPEN: " + issue.State + " " + issue.Title)
					issueTitle := strings.TrimSpace(issue.Title)

					_, seenInFile := todosFoundInFiles[issue.Title]
					if !seenInFile {
						// This issue doesn't have a corresponding TODO in the codebase anymore
						fmt.Printf("Closing GitHub issue #%d '%s' - TODO line no longer exists in codebase\n", issue.Number, issueTitle)

						if err := CloseGithubIssue(&GithubIssueUpdate{
							Title:        issueTitle,
							State:        "closed",
							State_Reason: "completed",
							Assignee:     "jonathon-chew",
							Number:       issue.Number,
							Id:           issue.Id,
						}); err != nil {
							log.Println(aphrodite.ReturnError(fmt.Sprintf("Error closing GitHub issue #%d: %v", issue.Number, err)))
						}

						foundNewTODO = true
					}
				} */
			}
		}
	}

	if !foundNewTODO {
		fmt.Println("No new todo found in any file in this directory")
	}
}

/*
Add a new task into the .tasks folder
*/
func (t *Task) AddTask() error {
	if !fileExists() {
		return fmt.Errorf("Unable to find the folder")
	}

	err := findTasks()
	if err != nil {
		return err
	}

	err = os.Chdir(folder)
	if err != nil {
		return err
	}

	err = os.Mkdir(t.Number, os.FileMode(os.O_CREATE))
	if err != nil {
		return err
	}

	err = os.WriteFile(t.Number, []byte(t.Details), os.FileMode(os.O_CREATE))
	if err != nil {
		return err
	}

	return nil
}

/*
Get a task that has been recorded and not done
*/
func (t *Task) GetTasks() error {
	if !fileExists() {
		return fmt.Errorf("Unable to find the folder")
	}

	err := findTasks()
	if err != nil {
		return err
	}

	err = os.Chdir(folder)
	if err != nil {
		return err
	}

	err = os.Mkdir(t.Number, os.FileMode(os.O_CREATE))
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(t.Number)
	if err != nil {
		return err
	}

	fmt.Printf("%s", string(contents))

	return nil
}

/*
Remove a task that has been completed
*/
func (t *Task) RemoveTasks(folder string) error {
	if !fileExists() {
		return fmt.Errorf("Unable to find the folder")
	}

	err := findTasks()
	if err != nil {
		return err
	}

	err = os.Remove(t.Number)
	if err != nil {
		return err
	}

	return nil
}

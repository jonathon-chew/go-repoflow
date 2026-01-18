package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	utils "github.com/jonathon-chew/go-repoflow/internal/Utils"
	"github.com/jonathon-chew/go-repoflow/pkg/git"
)

func CLI(CommandLineArguments []string) error {
	// aphrodite.PrintColour("Cyan", "I have found additional command line arguments, switching to CLI mode\n")

	var NoIssues error = errors.New("no GitHub issues found")

	for index, command := range CommandLineArguments {
		switch command {
		case "--commit-calendar", "--cc", "-cc":
			var option string
			if len(CommandLineArguments) > index+1 {
				option = CommandLineArguments[index+1]
				log.Print(option)
			}
			git.MakeCommitMap(option)
			return nil

		case "--check", "-c":
			entries := utils.MakeDirectoryList(utils.FindFilesInCurrentDirectory())

			for _, entry := range entries {
				ErrCheckingForUpdate := git.CheckForGitUpdate(entry)
				if ErrCheckingForUpdate != nil {
					return ErrCheckingForUpdate
				}
			}

		case "--clone", "-cl":
			git.CloneAllPublicRepos()
			return nil

		case "--get", "-get", "-g":
			returned, err := git.ListGithubIssues(true)
			if err != nil && errors.Is(err, NoIssues) {
				aphrodite.PrintWarning("no GitHub issues found")
				return nil
			}

			if err != nil {
				return err
			}

			var closedFlag, openFlag bool = false, true
			// Check for extra flags
			if len(os.Args) > 2 {
				for _, extraCommand := range os.Args[2:] {
					switch extraCommand {
					case "--closed", "-closed", "-c":
						closedFlag = true
					case "--all", "-all", "-a":
						openFlag = false
					}
				}
			}

			for index, issue := range returned {
				if closedFlag && issue.State == "closed" {
					fmt.Printf("%d The issue title is:\n%s\nThe body is: %s\nThe status is: %s\n\n", index+1, strings.TrimSpace(issue.Title), issue.Body, aphrodite.ReturnWarning(issue.State))
					fmt.Printf("______________\n")
					continue
				}

				if openFlag && issue.State == "open" {
					fmt.Printf("%d The issue title is:\n%s\nThe body is: %s\nThe status is: %s\n\n", index+1, strings.TrimSpace(issue.Title), issue.Body, aphrodite.ReturnInfo(issue.State))
					fmt.Printf("______________\n")
					continue
				}

				if !closedFlag && !openFlag {
					fmt.Printf("%d The issue title is:\n%s\nThe body is: %s\nThe status is: %s\n\n", index+1, strings.TrimSpace(issue.Title), issue.Body, issue.State)
					fmt.Printf("______________\n")
				}
			}

			return nil

		case "--set", "-set", "-s":
			var IssueTitle, IssueBody string
			if CommandLineArguments[index+1] == "title" || CommandLineArguments[index+1] == "--title" || CommandLineArguments[index+1] == "-title" || CommandLineArguments[index+1] == "-t" {
				IssueTitle = CommandLineArguments[index+2]
			} else {
				return errors.New("could not find a title flag proceeding the set command")
			}

			if CommandLineArguments[index+3] == "body" || CommandLineArguments[index+3] == "--body" || CommandLineArguments[index+3] == "-body" || CommandLineArguments[index+3] == "-b" {
				IssueBody = CommandLineArguments[index+4]
			} else {
				return errors.New("could not find a body flag proceeding the set command")
			}

			makeError := git.MakeGithubIssue(IssueTitle, IssueBody)
			if makeError != nil {
				fmt.Println(makeError)
				return makeError
			}

			return nil

		case "--version", "-version", "-v":
			fmt.Printf("v0.7.1\n")

		case "--help", "-help", "-h":

			aphrodite.PrintBold("Cyan", "No Arguments\n")
			aphrodite.PrintColour("Green", "You can run with no arguments to check all the files in the current directory for any undocumented todos and upload them to github\n\n")

			aphrodite.PrintBold("Cyan", "Get issues\n")
			aphrodite.PrintColour("Green", "You can pass in a get flag which will List the github issues, this can be supplimented with --open and --closed to filter to show only issues with those flags\n\n")

			aphrodite.PrintBold("Cyan", "Set issues\n")
			aphrodite.PrintColour("Green", "If you pass in the set flag, please pass in the title flag and body flag (in that order) to make a new issue with the relevent Title and Body\n\n")

			aphrodite.PrintBold("Cyan", "Version\n")
			aphrodite.PrintColour("Green", "Version Number can be passed in with the version flag\n\n")

			aphrodite.PrintBold("cyan", "Tags\n")
			aphrodite.PrintColour("Green", "Returns the latest tag following the format v[number].[number].[number]\n\n")

			aphrodite.PrintBold("cyan", "Increment Tag\n")
			aphrodite.PrintColour("Green", "Finds the biggest version number in the format format v[number].[number].[number] and adds 1 to the major / minor / patch numbers\n\n")

			aphrodite.PrintBold("cyan", "Open Issues\n")
			aphrodite.PrintColour("Green", "Open the github page on the issues page to manage from there\n\n")

			aphrodite.PrintBold("cyan", "Open Pull\n")
			aphrodite.PrintColour("Green", "Open the github page on the pull request page to manage from there\n\n")

			aphrodite.PrintBold("cyan", "Check\n")
			aphrodite.PrintColour("Green", "Check all folders 1 level deep to see if there are any updates required to push/pull\n\n")

			aphrodite.PrintBold("cyan", "Commit Calendar\n")
			aphrodite.PrintColour("Green", "Print to the terminal the git history activity for the last year!\n\n")

			aphrodite.PrintBold("cyan", "Clone\n")
			aphrodite.PrintColour("Green", "Clone all public repos into a temporary directory\n\n")

		case "--tags", "-tags", "-t", "--tag", "-tag":
			version, ErrGetLatestTag := git.GetLatestTag()
			if ErrGetLatestTag != nil {
				return ErrGetLatestTag
			}
			fmt.Println(version)

		case "--increment-tag", "-increment-tag", "-i", "--incrementtag", "-incrementtag":
			var argument string
			if len(CommandLineArguments) > index+1 {
				argument = CommandLineArguments[index+1]
			} else {
				argument = ""
			}
			ErrMakingNewTag := git.NewGitTag(argument)
			if ErrMakingNewTag != nil {
				return ErrMakingNewTag
			}

		case "--open", "-open", "-o":
			ErrOpeningRemoteOrigin := git.OpenRemoteOrigin("")
			if ErrOpeningRemoteOrigin != nil {
				return ErrOpeningRemoteOrigin
			}

		case "--open-issues", "-open-issues", "-oi":
			ErrOpeningRemoteOrigin := git.OpenRemoteOrigin("issues")
			if ErrOpeningRemoteOrigin != nil {
				return ErrOpeningRemoteOrigin
			}

		case "--open-pull", "-open-pull", "-op":
			ErrOpeningRemoteOrigin := git.OpenRemoteOrigin("pull")
			if ErrOpeningRemoteOrigin != nil {
				return ErrOpeningRemoteOrigin
			}

		}
	}

	return nil
}

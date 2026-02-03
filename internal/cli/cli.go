package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	utils "github.com/jonathon-chew/go-repoflow/internal/Utils"
	"github.com/jonathon-chew/go-repoflow/internal/git"
	"golang.org/x/term"
)

func ReducedMessage(issue git.GithubIssueResponse, index, width int) error {

	var message string
	minimumWidth := len(fmt.Sprintf("%d Title:  Body:  state: \n", index+1))

	switch issue.State {
	default:
		if width > len(issue.Title)+len(issue.Status)+minimumWidth {
			message = fmt.Sprintf("%d Title: %s state: %s\n", index+1, strings.TrimSpace(issue.Title), aphrodite.ReturnInfo(issue.State))
		} else { //  len(issue.Title) > width
			message = fmt.Sprintf("%d Title: %s\n", index+1, strings.TrimSpace(issue.Title))
		}
	case "open":
		message = fmt.Sprintf("%d Title: %s Body: %s state: %s\n", index+1, strings.TrimSpace(issue.Title), issue.Body, aphrodite.ReturnInfo(issue.State))
		if width >= len(message) {
			fmt.Print(message)
		} else {
			if minimumWidth >= width { // error
				return fmt.Errorf("[ERROR]: minimum width %d required is greater than the width of the terminal %d", minimumWidth, width)
			} else {
				// Check title len to check reduction of body could solve issue
				if width > len(issue.Title)+len(issue.Status)+minimumWidth {
					message = fmt.Sprintf("%d Title: %s state: %s\n", index+1, strings.TrimSpace(issue.Title), aphrodite.ReturnInfo(issue.State))
				} else { //  len(issue.Title) > width
					message = fmt.Sprintf("%d Title: %s\n", index+1, strings.TrimSpace(issue.Title))
				}
				// Check body length if just title?!
				fmt.Print(message)
			}
		}
		// fmt.Printf("______________\n")
	case "closed":
		message = fmt.Sprintf("%d Title: %s Body: %s state: %s\n", index+1, strings.TrimSpace(issue.Title), issue.Body, aphrodite.ReturnWarning(issue.State))

		if width >= len(message) {
			fmt.Print(message)
		} else {
			if minimumWidth >= width { // error
				return fmt.Errorf("[ERROR]: minimum width %d required is greater than the width of the terminal %d", minimumWidth, width)
			} else {
				// Check title len to check reduction of body could solve issue
				if width > len(issue.Title)+len(issue.Status)+minimumWidth {
					message = fmt.Sprintf("%d Title: %s state: %s\n", index+1, strings.TrimSpace(issue.Title), aphrodite.ReturnWarning(issue.State))
				} else { //  len(issue.Title) > width
					message = fmt.Sprintf("%d Title: %s\n", index+1, strings.TrimSpace(issue.Title))
				}
				// Check body length if just title?!
				fmt.Print(message)
			}
		}
		// fmt.Printf("______________\n")
	}

	return nil
}

func CLI(CommandLineArguments []string) error {
	// aphrodite.PrintColour("Cyan", "I have found additional command line arguments, switching to CLI mode\n")

	var NoIssues error = errors.New("no GitHub issues found")

	for index, command := range CommandLineArguments {
		switch command {
		default:
			if command != "minor" && command != "major" && command != "patch" {
				aphrodite.PrintError(command + " is not recognised")
			}
		case "--repo-stats", "-rs":
			RepoStats, ErrGettingRepoStats := git.GetRepoStats()
			if ErrGettingRepoStats != nil {
				return ErrGettingRepoStats
			}

			aphrodite.PrintInfo("Fork count: " + strconv.Itoa(RepoStats.Forks_count) + " Open Issue Count: " + strconv.Itoa(RepoStats.Open_issues_count) + " Stargazer's Count: " + strconv.Itoa(RepoStats.Stargazers_count) + " Watchers Count: " + strconv.Itoa(RepoStats.Watchers_count) + "\n")

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

		case "--get", "-get", "-g", "--list", "-list", "-l":
			returned, err := git.ListGithubIssues(true)
			if err != nil && errors.Is(err, NoIssues) {
				aphrodite.PrintWarning("no GitHub issues found")
				return nil
			}

			if err != nil {
				return err
			}

			var closedFlag, openFlag, oneLineFlag bool = false, true, false
			var width int
			var ErrGettingTerminalDetails error
			// Check for extra flags
			if len(os.Args) > 2 {
				for _, extraCommand := range os.Args[2:] {
					switch extraCommand {
					case "--closed", "-closed", "-c":
						closedFlag = true
						openFlag = false
					case "--all", "-all", "-a":
						openFlag = false
					case "--oneline", "--one-line", "-ol":
						oneLineFlag = true
						width, _, ErrGettingTerminalDetails = term.GetSize(int(os.Stdout.Fd()))
						if ErrGettingTerminalDetails != nil {
							return ErrGettingTerminalDetails
						}
					}
				}
			}

			for index, issue := range returned {
				if oneLineFlag == false {
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
				} else {
					if closedFlag && issue.State == "closed" {
						ErrMakingReducedMessage := ReducedMessage(issue, index, width)
						if ErrMakingReducedMessage != nil {
							return ErrMakingReducedMessage
						}
						continue
					}

					if openFlag && issue.State == "open" {
						ErrMakingReducedMessage := ReducedMessage(issue, index, width)
						if ErrMakingReducedMessage != nil {
							return ErrMakingReducedMessage
						}
						continue
					}

					if !closedFlag && !openFlag {
						ErrMakingReducedMessage := ReducedMessage(issue, index, width)
						if ErrMakingReducedMessage != nil {
							return ErrMakingReducedMessage
						}
						continue
					}
				}
			}

			return nil

		case "--set", "-set", "-s":
			var IssueTitle, IssueBody string
			var IssueLabel []string

			if CommandLineArguments[index+1] == "title" || CommandLineArguments[index+1] == "--title" || CommandLineArguments[index+1] == "-title" || CommandLineArguments[index+1] == "-t" {
				IssueTitle = CommandLineArguments[index+2]
			} else {
				return errors.New("could not find a title flag proceeding the set command")
			}

			if len(CommandLineArguments) < index+4 {
				return errors.New("could not find a body flag or text proceeding the set command")
			} else {
				if CommandLineArguments[index+3] == "body" || CommandLineArguments[index+3] == "--body" || CommandLineArguments[index+3] == "-body" || CommandLineArguments[index+3] == "-b" {
					IssueBody = CommandLineArguments[index+4]
				} else {
					return errors.New("could not find a body flag proceeding the set command")
				}
			}

			if len(CommandLineArguments) < index+6 {
				IssueLabel = []string{}
			} else {
				if CommandLineArguments[index+5] == "label" || CommandLineArguments[index+5] == "--label" || CommandLineArguments[index+5] == "-label" || CommandLineArguments[index+5] == "-l" {
					IssueLabel = append(IssueLabel, CommandLineArguments[index+6])
				} else {
					return errors.New("could not find a tag flag proceeding the set command")
				}
			}

			makeError := git.MakeGithubIssue(IssueTitle, IssueBody, IssueLabel)
			if makeError != nil {
				fmt.Println(makeError)
				return makeError
			}

			return nil

		case "--version", "-version", "-v":
			fmt.Printf("v0.7.14\n")

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

			aphrodite.PrintBold("cyan", "RepoStats")
			aphrodite.PrintColour("Green", "Get the repo stats from github Forks, Open Issues, Stargazer's, Watchers\n\n")

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

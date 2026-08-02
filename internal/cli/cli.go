package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/go-repoflow/internal/Utils"
	"github.com/jonathon-chew/go-repoflow/pkg/git"
	"github.com/jonathon-chew/go-repoflow/pkg/git/git_utils"

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

	// for index, command := range CommandLineArguments {
	for index := 0; index < len(CommandLineArguments); index++ {
		command := CommandLineArguments[index]

		switch command {
		default:
			aphrodite.PrintError(command + " is not recognised")
		case "--debug", "-debug", "-d":
			continue
		case "--repo-stats", "-rs":
			RepoStats, ErrGettingRepoStats := git.GetRepoStats()
			if ErrGettingRepoStats != nil {
				return ErrGettingRepoStats
			}

			/* 			aphrodite.PrintInfo("Fork count: " + strconv.Itoa(RepoStats.Forks_count) + " \nOpen Issue Count: " + strconv.Itoa(RepoStats.Open_issues_count) + " \nStargazer's Count: " + strconv.Itoa(RepoStats.Stargazers_count) + "\nWatchers Count: " + strconv.Itoa(RepoStats.Watchers_count) + "\n") */

			aphrodite.PrintInfo("Fork count: " + strconv.Itoa(RepoStats.Forks_count) + " \n" +
				"Open Issue Count: " + strconv.Itoa(RepoStats.Open_issues_count) + " \n" +
				"Stargazer's Count: " + strconv.Itoa(RepoStats.Stargazers_count) + "\n" +
				"Watchers Count: " + strconv.Itoa(RepoStats.Watchers_count) + "\n")

		case "--commit-calendar", "--cc", "-cc":
			var option string
			if len(CommandLineArguments) > index+1 {
				option = CommandLineArguments[index+1]
				log.Print(option)
			}
			git.MakeCommitMap(option)
			return nil

		case "--check", "-c":
			entries := git_utils.MakeDirectoryList(git_utils.FindFilesInCurrentDirectory())

			for _, entry := range entries {
				ErrCheckingForUpdate := git.CheckForGitUpdate(entry)
				if ErrCheckingForUpdate != nil {
					return ErrCheckingForUpdate
				}
			}

		case "--clone", "-cl":
			git.CloneAllPublicRepos()
			return nil

		case "--doctor":
			Utils.Doctor()
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

			// Flags controlling how issues/TODOs are displayed or processed.
			// By default we show open issues from GitHub (modifyFile = true).
			// Passing --local switches behaviour to operate only on local TODOs
			// without modifying files or contacting GitHub.
			var closedFlag, openFlag, oneLineFlag, modifyFile bool = false, true, false, true
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
					case "--local", "-local", "-l":
						modifyFile = false
					}
				}
			}

			if !modifyFile {
				git.ProcessTodosInRepo(modifyFile)
				index++
				continue
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
				index += 2
			} else {
				return errors.New("could not find a title flag proceeding the set command")
			}

			if len(CommandLineArguments) < index+4 {
				return errors.New("could not find a body flag or text proceeding the set command")
			} else {
				if CommandLineArguments[index+3] == "body" || CommandLineArguments[index+3] == "--body" || CommandLineArguments[index+3] == "-body" || CommandLineArguments[index+3] == "-b" {
					IssueBody = CommandLineArguments[index+4]
					index += 2
				} else {
					return errors.New("could not find a body flag proceeding the set command")
				}
			}

			if len(CommandLineArguments) < index+6 {
				IssueLabel = []string{}
			} else {
				if CommandLineArguments[index+5] == "label" || CommandLineArguments[index+5] == "--label" || CommandLineArguments[index+5] == "-label" || CommandLineArguments[index+5] == "-l" {
					IssueLabel = append(IssueLabel, CommandLineArguments[index+6])
					index += 2
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
			help()
			return nil

		case "--tags", "-tags", "-t", "--tag", "-tag":
			version, ErrGetLatestTag := git.GetLatestTag(false)
			if ErrGetLatestTag != nil {
				return ErrGetLatestTag
			}
			fmt.Println(version)

		case "--Change-log", "-changelog", "--changelog", "--log", "--change":
			git.MakeChangeLog(".")

		case "--increment-tag", "-increment-tag", "-i", "--incrementtag", "-incrementtag":
			var argument string
			var force bool
			if len(CommandLineArguments) > index+1 {
				argument = CommandLineArguments[index+1]
			} else {
				argument = ""
			}

			if len(CommandLineArguments) > index+2 {
				if strings.ToLower(CommandLineArguments[index+2]) == "true" || strings.ToLower(CommandLineArguments[index+2]) == "t" || strings.ToLower(CommandLineArguments[index+2]) == "force" || strings.ToLower(CommandLineArguments[index+2]) == "f" {
					force = true
					index += 2
				}
			}

			ErrMakingNewTag := git.NewGitTag(argument, force)
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

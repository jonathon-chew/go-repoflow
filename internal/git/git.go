package git

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	utils "github.com/jonathon-chew/go-repoflow/internal/Utils"
)

var HTTPStatusResponseMeanings = map[string]string{
	"201": "Created",
	"400": "Bad Request",
	"401": "Unauthorized",
	"403": "Forbidden",
	"404": "Resource not found",
	"410": "Gone",
	"422": "Validation failed, or the endpoint has been spammed.",
	"503": "Service unavailable",
}

type Credentials struct {
	Owner string
	Repo  string
	Token string
}

// type CommitMap map[string]int

// UTILS
func GetRemoteOrigin() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		return "", err
	}

	return out.String(), nil
}

func FindGitFolder() bool {

	directoryList := utils.MakeDirectoryList(utils.FindFilesInCurrentDirectory())

	// Look in the directories for a git folder
	if !slices.Contains(directoryList, ".git") {
		fmt.Println("[ERROR]: No git folder found")
		return false // recursively look?
	}

	return true
}

func OpenRemoteOrigin(place string) error {
	url, ErrGetRemote := GetRemoteOrigin()
	if ErrGetRemote != nil {
		return ErrGetRemote
	}

	url = strings.TrimSpace(url)

	if strings.Contains(url, "github.com") && place != "" {
		switch place {
		case "pull":
			url = url + "/pulls"
		case "issues":
			url = url + "/issues"
		}
	} else if place != "" {
		return fmt.Errorf("[ERROR]: only github.com has been implimented so far")
	}

	cmd := exec.Command("open", url)

	ErrRun := cmd.Run()
	if ErrRun != nil {
		fmt.Printf("Error: %s\n", ErrRun)
		return ErrRun
	}

	return nil
}

// GIT TAG
func getTags() (string, error) {
	cmd := exec.Command("git", "tag")

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		return "", err
	}

	versions := out.String()

	if len(versions) == 0 {
		return "", nil
	}

	return versions, nil
}

func GetLatestTag() (string, error) {

	if !FindGitFolder() {
		return "", fmt.Errorf("[Error]: Unable to find a git folder in the current directory")
	}

	versions, err := getTags()
	if err != nil {
		return "", fmt.Errorf("[Error]: Unable to successfully get the tags\n ")
	}

	if versions == "" {
		// There was nothing back from Get Tags therefore we should make one
		return "", nil
	}

	versionList := strings.Split(versions, "\n")

	// If the list is only 1 item long it's the biggest, so early return
	if len(versionList) == 1 {
		return versionList[0], nil
	}

	var biggestMajor, biggestMinor, biggestPatch int
	var latestVersion string

	for _, version := range versionList {
		version = strings.TrimSpace(version)
		if version == "" {
			continue
		}

		if len(version) < 4 {
			continue
		}

		if !strings.Contains(version, ".") || !strings.HasPrefix(version, "v") {
			fmt.Printf("[WARNING]: Skipping looking at tag %s, as doesn't follow the convention v.[0-9].[0-9].[0-9]\n", version)
			continue
		}

		versionParts := strings.Split(version[1:], ".")
		if len(versionParts) < 3 {
			continue
		}

		major, ErrMajorConv := strconv.Atoi(versionParts[0])
		minor, ErrMinorConv := strconv.Atoi(versionParts[1])
		patch, ErrPatchConv := strconv.Atoi(versionParts[2])

		if ErrMajorConv != nil || ErrMinorConv != nil || ErrPatchConv != nil {
			fmt.Printf("[WARNING]: Skipping tag %s due to conversion error\n", version)
			continue
		}

		// Check if this version is greater than the current latest
		if major > biggestMajor ||
			(major == biggestMajor && minor > biggestMinor) ||
			(major == biggestMajor && minor == biggestMinor && patch > biggestPatch) {
			biggestMajor = major
			biggestMinor = minor
			biggestPatch = patch
			latestVersion = version
		}
	}

	if latestVersion == "" {
		return "", nil
	}

	return latestVersion, nil
}

func makeTag(newTag string) error {
	cmd := exec.Command("git", "tag", newTag, "-m", "Release Version: "+strings.ReplaceAll(newTag, "v", ""))

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", stderr.String())
		return err
	}

	aphrodite.PrintInfo(fmt.Sprintf("New latest tag:%s\n", newTag))

	aphrodite.PrintBold("Cyan", "Do you want to push the new tag to git?\n")

	var userChoicePushToGit string
	_, ErrGettingUserChioce := fmt.Scan(&userChoicePushToGit)
	if ErrGettingUserChioce != nil {
		return ErrGettingUserChioce
	}

	if userChoicePushToGit == "y" || userChoicePushToGit == "Y" || userChoicePushToGit == "yes" || userChoicePushToGit == "Yes" || userChoicePushToGit == "YES" {
		aphrodite.PrintInfo("Pushing to remote git respository.\n")
		// git push --tags --force-with-lease=false
		tagPushCmd := exec.Command("git", "push", "--tags", "--force-with-lease=false")
		ErrPushingTags := tagPushCmd.Run()
		if ErrPushingTags != nil {
			return ErrPushingTags
		}
		aphrodite.PrintInfo("Successfully pushed.\n")
	}

	return nil
}

func NewGitTag(argument string) error {
	version, ErrGetLatestTag := GetLatestTag()
	if ErrGetLatestTag != nil {
		return ErrGetLatestTag
	}

	if version == "" {
		ErrMakingTag := makeTag("v0.1.0")
		if ErrMakingTag != nil {
			return ErrMakingTag
		}
		return nil
	}

	fmt.Println("Current latest tag: ", version)

	if argument != "major" && argument != "minor" && argument != "patch" {
		var userChoiceVersionUpdate string

		fmt.Printf("Do you want to increase the major, minor or patch of the tag?\n")

		_, ErrUserInput := fmt.Scanln(&userChoiceVersionUpdate)
		if ErrUserInput != nil {
			return ErrUserInput
		}
		if userChoiceVersionUpdate != "major" && userChoiceVersionUpdate != "minor" && userChoiceVersionUpdate != "patch" {
			return fmt.Errorf("[ERROR]: user input was not major, minor or patch")
		} else {
			argument = userChoiceVersionUpdate
		}
	}

	major, ErrMajorConv := strconv.Atoi(strings.Split(version[1:], ".")[0])
	if ErrMajorConv != nil {
		return ErrMajorConv
	}

	minor, ErrMinorConv := strconv.Atoi(strings.Split(version[1:], ".")[1])
	if ErrMinorConv != nil {
		return ErrMinorConv
	}

	patch, ErrPatchConv := strconv.Atoi(strings.Split(version[1:], ".")[2])
	if ErrPatchConv != nil {
		return ErrPatchConv
	}
	var newTag string

	switch argument {
	case "major":
		newMajor := major + 1
		newTag = fmt.Sprintf("v%d.%d.%d", newMajor, 0, 0)
	case "minor":
		newMinor := minor + 1
		newTag = fmt.Sprintf("v%d.%d.%d", major, newMinor, 0)
	case "patch":
		newPatch := patch + 1
		newTag = fmt.Sprintf("v%d.%d.%d", major, minor, newPatch)
	default:
		return errors.New(argument + " was not recognised as a valid command")
	}

	ErrMakingTag := makeTag(newTag)
	if ErrMakingTag != nil {
		return ErrMakingTag
	}

	return nil
}

// MAKE A GIT REQUEST
func genericGitRequest() (Credentials, error) {
	remoteOrigin, err := GetRemoteOrigin()
	var credentials Credentials
	if err != nil {
		fmt.Printf("Unable to get the remote origin\n")
		return credentials, err
	}

	if strings.Contains(remoteOrigin, "github") {

		gitUrl := strings.ReplaceAll(remoteOrigin, ".git", "")
		gitDetails := strings.Split(strings.ReplaceAll(gitUrl, "https://github.com/", ""), "/")

		credentials.Owner = gitDetails[0]
		credentials.Repo = strings.Replace(gitDetails[1], "\n", "", -1)
		credentials.Token = os.Getenv("GH_PERSONAL_TOKEN")

		if credentials.Token == "" {
			_, VarExists := os.LookupEnv("GH_PERSONAL_TOKEN")
			if VarExists {
				return credentials, errors.New("GH_PERSONAL_TOKEN is empty")
			} else {
				return credentials, errors.New("no GH_PERSONAL_TOKEN in the environment")
			}
		}

		return credentials, nil

	} else if strings.Contains(remoteOrigin, "gitlab") {

		gitUrl := strings.ReplaceAll(remoteOrigin, ".git", "")
		gitDetails := strings.Split(strings.ReplaceAll(gitUrl, "https://gitlab.", ""), "/")

		credentials.Owner = gitDetails[0] // check this still applies for gitlab - as i'm not sure it does, this might need to be a git call
		credentials.Repo = strings.Replace(gitDetails[1], "\n", "", -1)
		credentials.Token = os.Getenv("GL_PERSONAL_TOKEN")

		if credentials.Token == "" {
			_, VarExists := os.LookupEnv("GL_PERSONAL_TOKEN")
			if VarExists {
				return credentials, errors.New("GL_PERSONAL_TOKEN is empty")
			} else {
				return credentials, errors.New("no GL_PERSONAL_TOKEN in the environment")
			}
		}

		return credentials, nil

	} else {
		return credentials, fmt.Errorf("the remote origin is not github/gitlab, and the ability to create issues for %s is not currently implimented", remoteOrigin)
	}
}

// Entry is the folder that you would like to check if their is an update to git in it.
// Only does it in the root directory, if recusively going into folders it won't return false positives
// The only time would be a submodule
func CheckForGitUpdate(entry string) error {

	currentWorkingDirectory, ErrGettingWorkingDirectory := os.Getwd()
	if ErrGettingWorkingDirectory != nil {
		return fmt.Errorf("getwd: %w", ErrGettingWorkingDirectory)
	}

	dirPath := filepath.Join(currentWorkingDirectory, entry)

	// Fast repo check: is there a .git directory / file?
	if _, err := os.Stat(filepath.Join(dirPath, ".git")); err != nil {
		return errors.New("not a git folder") // not a git repo (or not accessible)
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dirPath

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git status failed in %s: %s", dirPath, stderr.String())
	}

	if out.Len() > 0 {
		fmt.Printf("%s has a git update\n%s\n", entry, out.String())
	}

	aheadCmd := exec.Command("git", "rev-list", "--count", "@{u}..HEAD")
	aheadCmd.Dir = dirPath

	var aheadOut bytes.Buffer
	aheadCmd.Stdout = &aheadOut

	if err := aheadCmd.Run(); err == nil {
		if strings.TrimSpace(aheadOut.String()) != "0" {
			fmt.Printf("%s has commits to push\n", entry)
		}
	}

	return nil
}

func MakeCommitMap(option string) {

	root := "." // You can make this configurable
	repos := utils.FindGitRepos(root)
	var totalCount int

	totalCommits := make(utils.CommitMap)
	for _, repo := range repos {
		// fmt.Println("Scanning:", repo)
		commits := getCommitDates(repo)
		for date, count := range commits {
			totalCommits[date] += count
			totalCount += count
		}
	}

	aphrodite.PrintInfo("Total Count: " + strconv.Itoa(totalCount) + "\n")
	utils.RenderDateGraph(totalCommits, option)
}

func getCommitDates(repo string) utils.CommitMap {
	cmd := exec.Command("git", "log", "--pretty=format:%ad", "--date=short")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error reading commits from", repo, err)
		return nil
	}
	commits := make(utils.CommitMap)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		date := scanner.Text()
		commits[date]++
	}
	return commits
}

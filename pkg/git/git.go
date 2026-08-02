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
	git_utils "github.com/jonathon-chew/go-repoflow/pkg/git/git_utils"
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

	directoryList := git_utils.MakeDirectoryList(git_utils.FindFilesInCurrentDirectory())

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

func GetLatestTag(checkForGitFolder bool) (string, error) {

	if !FindGitFolder() && checkForGitFolder {
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

func makeTag(newTag string, force bool) error {
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

	var userChoicePushToGit string

	if force {
		userChoicePushToGit = "y"
	} else {
		aphrodite.PrintBold("Cyan", "Do you want to push the new tag to git?\n")
		_, ErrGettingUserChioce := fmt.Scan(&userChoicePushToGit)
		if ErrGettingUserChioce != nil {
			return ErrGettingUserChioce
		}
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

func NewGitTag(argument string, force bool) error {
	version, ErrGetLatestTag := GetLatestTag(false)
	if ErrGetLatestTag != nil {
		return ErrGetLatestTag
	}

	if version == "" {
		ErrMakingTag := makeTag("v0.1.0", force)
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

	ErrMakingTag := makeTag(newTag, force)
	if ErrMakingTag != nil {
		return ErrMakingTag
	}

	return nil
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
		aphrodite.PrintError(fmt.Sprintf("%s not a git folder", dirPath)) // not a git repo (or not accessible)
		return nil
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
	repos := git_utils.FindGitRepos(root)
	var totalCount int

	totalCommits := make(git_utils.CommitMap)
	for _, repo := range repos {
		// fmt.Println("Scanning:", repo)
		commits := getCommitDates(repo)
		for date, count := range commits {
			totalCommits[date] += count
			totalCount += count
		}
	}

	aphrodite.PrintInfo("Total Count: " + strconv.Itoa(totalCount) + "\n")
	git_utils.RenderDateGraph(totalCommits, option)
}

func getCommitDates(repo string) git_utils.CommitMap {
	cmd := exec.Command("git", "log", "--pretty=format:%ad", "--date=short")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error reading commits from", repo, err)
		return nil
	}
	commits := make(git_utils.CommitMap)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		date := scanner.Text()
		commits[date]++
	}
	return commits
}

func MakeChangeLog(repo string) {
	// COMMAND := "git log --pretty=format:(%H) %d | %s"
	cmd := exec.Command("git", "log", "--pretty=format:(%H) %D | %s")

	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error reading commits from", repo, err)
		return
	}

	type LineEntry struct {
		CommitHash   string
		TagValue     string
		Tag          string
		CommitString string
	}

	commitHistory := []LineEntry{}
	/* previousTag, err := GetLatestTag(false)
	if err != nil {
		fmt.Print("Unbale to get the latest tag")
		return
	} */

	previousTag := "LATEST"

	tag := previousTag
	var tagList []string

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		commit := scanner.Text()
		splitCommitEntry := strings.SplitN(commit, "|", 2)

		metaData := strings.SplitN(splitCommitEntry[0], " ", 3)

		hash, tagValue := "", ""
		if len(metaData) == 3 {
			hash = metaData[0]
			tagValue = metaData[1]
			tag = metaData[2]

			if tag != "" && strings.Contains(tag, "v") && strings.Contains(tag, ".") {
				tagList = append(tagList, tag)
				previousTag = tag
			}
		} else if len(metaData) == 1 {
			hash = metaData[0]
			tagValue = ""
			tag = previousTag
		}

		commitHistory = append(commitHistory, LineEntry{
			CommitHash:   hash,
			TagValue:     tagValue,
			Tag:          previousTag,
			CommitString: splitCommitEntry[1],
		})
	}

	fileP, err := os.Create("./CHANGELOG.md")

	var content strings.Builder
	content.WriteString("# CHANGELOG")

	// previousTag = tagList[0]
	var eachTag strings.Builder
	var news, updates, refactors, deletes, miscs []string //  []string

	for i, eachEntry := range commitHistory {

		if i == 0 {
			previousTag = eachEntry.Tag
		}

		if eachEntry.Tag != previousTag {
			fmt.Fprintf(&eachTag, "\n\n\t## %s\n", eachEntry.Tag)

			if len(news) > 0 {
				fmt.Fprintf(&eachTag, "\n\t\t### NEW\n\t\t%s", strings.Join(news, "\n\t\t"))
			}

			if len(updates) > 0 {
				fmt.Fprintf(&eachTag, "\n\t\t### UPDATES\n\t\t%s", strings.Join(updates, "\n\t\t"))
			}

			if len(refactors) > 0 {
				fmt.Fprintf(&eachTag, "\n\t\t### REFACTORS\n\t\t%s", strings.Join(refactors, "\n\t\t"))
			}

			if len(deletes) > 0 {
				fmt.Fprintf(&eachTag, "\n\t\t### DELETES\n\t\t%s", strings.Join(deletes, "\n\t\t"))
			}

			if len(miscs) > 0 {
				fmt.Fprintf(&eachTag, "\n\t\t### MISC\n\t\t%s", strings.Join(miscs, "\n\t\t"))
			}

			content.WriteString(eachTag.String())
			news = []string{}
			updates = []string{}
			refactors = []string{}
			deletes = []string{}
			miscs = []string{}
			previousTag = eachEntry.Tag

			eachTag.Reset()
		}

		// fmt.Printf("each entry tag = %v previous tag = %v\n", eachEntry.Tag, previousTag)

		var commitType, commitMessage string
		if strings.Contains(eachEntry.CommitString, ":") {
			temp_commit := strings.SplitN(eachEntry.CommitString, ":", 2)
			commitType = temp_commit[0]
			commitMessage = temp_commit[1]
		} else {
			commitType = ""
			commitMessage = eachEntry.CommitString
		}

		// fmt.Printf("committype[0] is %v\n", commitType[0])

		switch strings.TrimSpace(commitType) {
		case "new":
			news = append(news, "1."+commitMessage)
		case "update":
			updates = append(updates, "1."+commitMessage)
		case "refactor":
			refactors = append(refactors, "1."+commitMessage)
		case "delete":
			deletes = append(deletes, "1."+commitMessage)
		default:
			miscs = append(miscs, "1."+commitMessage)
		}
	}

	os.WriteFile(fileP.Name(), []byte(content.String()), os.ModeAppend)
}

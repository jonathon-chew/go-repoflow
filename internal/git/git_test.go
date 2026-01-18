package git

import (
	"fmt"
	"os"
	"strings"
	"testing"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func TestRemoteURL(t *testing.T) {
	t.Logf("Testing GetRemoteOrigin")
	url, err := GetRemoteOrigin()
	if err != nil {
		t.Fatal(aphrodite.ReturnError(fmt.Sprintf("Failed to get remote origin: %v", err)))
	}

	t.Logf("Remote URL: %s", url)
}

func TestFindGitFolder(t *testing.T) {
	t.Log("Testing GetGitTag")

	os.Chdir("../../")

	if FindGitFolder() {
		t.Log("Successfully found a git folder")
	} else {
		t.Error(aphrodite.ReturnError("Did not find a git folder"))
	}
}

func TestGitTag(t *testing.T) {
	t.Log("Testing GetGitTag")

	returnString, ErrGettingTags := getTags()
	if ErrGettingTags != nil {
		t.Error(aphrodite.ReturnError(fmt.Sprintf("%s", ErrGettingTags)))
	}

	t.Log(returnString)
}

func TestHelpVersionMatchesLatestGitTag(t *testing.T) {
	t.Log("Testing whether the help function version matches the latest git tag")

	_, ErrFindingFile := os.Stat("./internal/cli/cli.go")
	if ErrFindingFile != nil {
		t.Error(aphrodite.ReturnError("[Error]: Could not find the file"))
	}

	fileContentsBytes, err := os.ReadFile("./internal/cli/cli.go")
	if err != nil {
		t.Error(aphrodite.ReturnError(fmt.Sprintf("There was an error opening the file %s", err)))
		return
	}

	splitFile := strings.Split(string(fileContentsBytes), "\n")

	for index, line := range splitFile {
		if strings.TrimSpace(line) == `case "--version", "-version", "-v":` {

			helpVersionLine := strings.TrimSpace(splitFile[index+1])
			helpTag := strings.Replace(helpVersionLine, `fmt.Printf("`, "", 1)
			helpTag = strings.Replace(helpTag, `\n")`, "", 1)

			actualTag, err := GetLatestTag()
			if err != nil {
				t.Error("Unable to get latest tag to compare")
			}

			lineShouldBe := fmt.Sprintf(`fmt.Printf("%s\n")`, actualTag)

			if helpVersionLine == lineShouldBe {
				t.Log("Versions match")
				return
			} else {
				t.Errorf("Versions don't match Actual Tag: %s Version: %s", actualTag, helpTag)
				return
			}
		}
	}

	t.Error("Unable to find a version line in the CMD")
}

func TestLatestGitTag(t *testing.T) {
	t.Log("Testing GetLatestGitTag")

	returnString, ErrGettingTags := GetLatestTag()
	if ErrGettingTags != nil {
		t.Error(ErrGettingTags)
	}

	t.Log(returnString)
}

func TestListIssues(t *testing.T) {
	t.Logf("Testing ListGithubIssues")
	returned, err := ListGithubIssues(false) // false is NOT passed from the CLI so will always report if it connected to github
	if err != nil {
		t.Fatal(aphrodite.ReturnError(fmt.Sprintf("Failed to get remote origin: %v", err)))
	} else {
		t.Logf("Found issues %d", len(returned))
	}

	// for _, i := range returned{
	// 	t.Logf("From the test file: The title from the API was: %v\n", i.Title)
	// }

}

func TestGenericGit(t *testing.T) {
	t.Logf("Testing GetRemoteOrigin")
	GitCredentials, err := genericGitRequest()
	if err != nil {
		t.Fatal(aphrodite.ReturnError(fmt.Sprintf("Failed to get Git data: %v", err)))
	}

	t.Logf("Owner: %s, Repo: %s, Token: %s", GitCredentials.Owner, GitCredentials.Repo, GitCredentials.Token)
}

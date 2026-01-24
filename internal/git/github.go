package git

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	utils "github.com/jonathon-chew/go-repoflow/internal/Utils"
)

type Limit struct {
	Limit     int `json:"limit"`
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"`
}

type RateLimit struct {
	Core      Limit `json:"core"`
	Resources struct {
		Search                      Limit `json:"Search"`
		Graphql                     Limit `json:"graphql"`
		Integration_manifest        Limit `json:"integration_manifest"`
		Source_import               Limit `json:"source_import"`
		Code_scanning_upload        Limit `json:"code_scanning_upload"`
		Actions_runner_registration Limit `json:"actions_runner_registration"`
		Scim                        Limit `json:"scim"`
		Dependency_snapshots        Limit `json:"dependency_snapshots"`
		Code_search                 Limit `json:"code_search"`
		Code_scanning_autofix       Limit `json:"code_scanning_autofix"`
	}
	Rate Limit `json:"rate"`
}

// GITHUB STRUCTS
type Github_Assignee struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Github_Issue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Milestone int      `json:"milestone,omitempty"`
	Label     []string `json:"labels,omitempty"`
	Assignees string   `json:"assignees,omitempty"`
}

type Github_Label struct {
}

type GithubIssueResponse struct {
	Url            string `json:"url"`
	Repository_url string `json:"repository_url"`
	Labels_url     string `json:"labels_url"`
	Comments_url   string `json:"comments_url"`
	Events_url     string `json:"events_url"`
	Id             int    `json:"id"`
	Node_id        string `json:"node_id"`
	Number         int    `json:"number"`
	Title          string `json:"title"`
	User           struct {
		Login          string `json:"login"`
		Id             int    `json:"id"`
		Repos_url      string `json:"repos_url"`
		Events_url     string `json:"events_url"`
		Type           string `json:"type"`
		User_view_type string `json:"user_view_type"`
		Site_admin     bool   `json:"site_admin"`
	} `json:"user"`
	Labels             []Github_Label    `json:"labels"`
	State              string            `json:"state"`
	State_Reason       string            `json:"state_reason"`
	Locked             bool              `json:"locked"`
	Assignee           Github_Assignee   `json:"assignee"`
	Assignees          []Github_Assignee `json:"assignees"`
	Comments           int               `json:"comments"`
	Created_at         string            `json:"created_at"`
	Updated_at         string            `json:"updated_at"`
	Author_association string            `json:"author_association"`
	Active_lock_reason string            `json:"active_lock_reason"`
	Body               string            `json:"body"`
	Message            string            `json:"message"`
	Status             string            `json:"status"`
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"html_url"`
	Star        int    `json:"stargazers_count"`
}

type User struct {
	Public_repos int `json:"public_repos"`
}

type RepoInformation struct {
	Forks_count       int `json:"forks_count"`
	Forks             int `json:"forks"`
	Stargazers_count  int `json:"stargazers_count"`
	Watchers_count    int `json:"watchers_count"`
	Open_issues_count int `json:"open_issues_count"`
}

func conntactGithub[T any](websiteUrl string, token string) (T, error) {

	var v T

	request, err := http.NewRequest("GET", websiteUrl, nil)
	if err != nil {
		return v, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := http.Client{}

	req, err := client.Do(request)
	if err != nil {
		return v, err
	}

	defer req.Body.Close()

	responseBody, err := io.ReadAll(req.Body)
	if err != nil {
		return v, err
	}

	// fmt.Printf("Repsonse Body: %s\n\n", string(responseBody))

	if err := json.Unmarshal(responseBody, &v); err != nil {
		return v, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if req.StatusCode != http.StatusOK {
		return v, fmt.Errorf("GitHub API error: %s", req.Status)
	}

	return v, nil
}

func GetRateLimit() (RateLimit, error) {

	var rateLimit RateLimit

	GitCredentials, err := getGitCredentials()
	if err != nil {
		return rateLimit, err
	}

	RateLimit, ErrContactingGithub := conntactGithub[RateLimit]("https://api.github.com/rate_limit", GitCredentials.Token)
	if ErrContactingGithub != nil {
		return rateLimit, ErrContactingGithub
	}

	return RateLimit, nil
}

func GetRepoStats() (RepoInformation, error) {

	var repoInformation RepoInformation

	GitCredentials, err := getGitCredentials()
	if err != nil {
		return repoInformation, err
	}

	RepoInformation, ErrContactingGithub := conntactGithub[RepoInformation](fmt.Sprintf("https://api.github.com/repos/%s/%s", GitCredentials.Owner, GitCredentials.Repo), GitCredentials.Token)
	if ErrContactingGithub != nil {
		return RepoInformation, ErrContactingGithub
	}

	return RepoInformation, nil
}

// LIST GIT ISSUES
func ListGithubIssues(passedFromCLI bool) ([]GithubIssueResponse, error) {
	// "https://api.github.com/repos/%s/%s/issues?state=all"
	var ResponseInstance []GithubIssueResponse

	GitCredentials, err := getGitCredentials()
	if err != nil {
		return ResponseInstance, err
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all", GitCredentials.Owner, GitCredentials.Repo), nil)
	if err != nil {
		return ResponseInstance, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", GitCredentials.Token))

	client := http.Client{}

	req, err := client.Do(request)
	if err != nil {
		return ResponseInstance, err
	}

	defer req.Body.Close()

	if !passedFromCLI {
		fmt.Printf("The response was: %s, %s\n\n", req.Status, HTTPStatusResponseMeanings[req.Status])
	}

	responseBody, err := io.ReadAll(req.Body)
	if err != nil {
		return ResponseInstance, err
	}

	// fmt.Printf("Repsonse Body: %s\n\n", string(responseBody))

	if err := json.Unmarshal(responseBody, &ResponseInstance); err != nil {
		return ResponseInstance, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if len(ResponseInstance) == 0 {
		// return ResponseInstance, errors.New("no GitHub issues found")
		return ResponseInstance, nil
	}

	if req.StatusCode != http.StatusOK {
		return ResponseInstance, fmt.Errorf("GitHub API error: %s", req.Status)
	}

	// fmt.Printf("ResponseInstance: %v\n\n", ResponseInstance)

	// for _, response := range ResponseInstance {
	// 	fmt.Println("The title for the response is: ", strings.TrimSpace(response.Title), " with ID: ", response.Id)
	// }

	return ResponseInstance, nil
}

func MakeGithubIssue(TITLE, BODY string) error {

	// Get the credentials required
	GithubCredentials, err := getGitCredentials()
	if err != nil {
		return err
	}

	// Create the issue using a struct
	issue := Github_Issue{
		Title: strings.TrimSpace(TITLE),
		Body:  BODY,
	}

	// Convert the struct into JSON using the tags and Marshal
	jsonData, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	// Convert the JSON into bytes
	requestBody := bytes.NewBuffer(jsonData)

	// Make the request
	request, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", GithubCredentials.Owner, GithubCredentials.Repo), io.Reader(requestBody))
	if err != nil {
		fmt.Printf("Error making the HTTP request %s\n", err)
		return err
	}

	// Set the required headers
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GithubCredentials.Token))

	// Make a new client
	client := http.Client{}

	// Complete the request - Client.Do because the http.NewRequest handles the method
	req, err := client.Do(request)
	if err != nil {
		return err
	}

	if req.StatusCode != 200 && req.StatusCode != 201 {
		fmt.Println(req.Body)
		return fmt.Errorf("the response was not positive, %d", req.StatusCode)
	}

	fmt.Printf("The response was: %s, %s\n", req.Status, HTTPStatusResponseMeanings[req.Status])

	return nil
}

// Get the github credentials based on the env variable for github, and the parsing of hte git remote
func getGitCredentials() (Credentials, error) {

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

	} else {
		return credentials, fmt.Errorf("the remote origin is not github/gitlab, and the ability to create issues for %s is not currently implimented", remoteOrigin)
	}
}

// REMOVE GIT ISSUES
// (#2) TODO: Add the ability to remove to dos which have been closed on github
func RemoveLineDueToGithubIssue(line string, foundGithubIssues []GithubIssueResponse) (bool, error) {

	// Loop through the issues and compare to the line
	for _, issue := range foundGithubIssues {
		if strings.Contains(strings.TrimSpace(line), issue.Title) {
			err := CloseGithubIssue(&issue)
			if err != nil {
				return true, err // trying this out - as first half the of the function was "completed" successfully but the second half wasn't!
			}
			return true, nil
		}
	}

	// If the loop didn't find anything return false and no error!
	return false, nil
}

// (#3) TODO: Add the ability to close issues on github which have been removed from the code base

func CloseGithubIssue(closeIssue *GithubIssueResponse) error {

	// Put together the JSON message required to close an issue
	closeIssue.State = "closed"
	closeIssue.State_Reason = "completed"

	// Get the credentials
	GithubCredentials, err := getGitCredentials()
	if err != nil {
		return err
	}

	// Convert the struct into JSON using the tags and Marshal
	jsonData, err := json.Marshal(closeIssue)
	if err != nil {
		return err
	}

	// Convert the JSON into bytes
	requestBody := bytes.NewBuffer(jsonData)

	// Write the request
	request, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", GithubCredentials.Owner, GithubCredentials.Repo, closeIssue.Number), requestBody)
	if err != nil {
		return err
	}

	// Set the required headers
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", GithubCredentials.Token))

	client := http.Client{}

	// Make the request
	closeGithubIssueResponse, clientErr := client.Do(request)
	if clientErr != nil {
		fmt.Printf("The response from github was: %s\n", HTTPStatusResponseMeanings[closeGithubIssueResponse.Status])
		return clientErr
	}

	fmt.Printf("The response from github was: %s\n", HTTPStatusResponseMeanings[closeGithubIssueResponse.Status])

	// Return if error?
	return nil

}

func CloneAllPublicRepos() {

	userName, ErrGettingUserName := utils.GetUserInput([]byte("What is the name of the user/org you would like to clone? \n"))
	if ErrGettingUserName != nil {
		return
	}

	confirmPrompt, ErrGettingConfirmedPrompt := utils.GetUserInput([]byte("We're going to get everything from: " + userName + " y/Y? \n"))
	if ErrGettingConfirmedPrompt != nil {
		return
	}

	if confirmPrompt != "y" && confirmPrompt != "Y" {
		os.Stdin.Write([]byte("You've elected not to carry on"))
		return
	}

	var UserUrl string = "https://api.github.com/users/" + userName
	userReq, err := http.Get(UserUrl)
	if err != nil {
		log.Fatal(err)
	}
	var userDetails User
	if err := json.NewDecoder(userReq.Body).Decode(&userDetails); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	userReq.Body.Close()

	if userDetails.Public_repos > 50 {
		userReponse, ErrGettingConfirmLargeDownload := utils.GetUserInput([]byte("There are " + strconv.Itoa(userDetails.Public_repos) + " repos to clone - are you sure? y/Y\n"))
		if ErrGettingConfirmLargeDownload != nil {

			return
		}

		if userReponse != "y" && userReponse != "Y" {
			log.Fatal("Too many repositories, user has elected to stop")
			return
		}
	}

	var RepoURL string = "https://api.github.com/users/" + userName + "/repos"
	repoReq, err := http.Get(RepoURL)
	if err != nil {
		log.Fatal(err)
	}

	defer repoReq.Body.Close()

	var repos []Repo
	if err := json.NewDecoder(repoReq.Body).Decode(&repos); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Write header
	os.Stdin.Write([]byte("# GitHub Repositories"))

	utils.NewDirectory()
	ErrMovingDirectory := os.Chdir(utils.TemporaryDirectory)
	if ErrMovingDirectory != nil {
		log.Fatal(ErrMovingDirectory)
		return
	}

	// Write out each repo as it's processed
	for _, repo := range repos {
		if repo.Name != userName {
			os.Stdin.Write([]byte("Name: " + repo.Name + "\n"))
			os.Stdin.Write([]byte("Description: " + repo.Description + "\n"))
			os.Stdin.Write([]byte("URL: " + repo.Url + "\n\n"))

			cmd := exec.Command("git", "clone", "--depth", "1", repo.Url)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				log.Printf("Error: %s\n", stderr.String())
				return
			}

		}
	}
}

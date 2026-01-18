package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Create_Gitlab_Issue struct {
	Title                                   string   `json:"title"`
	Created_at                              string   `json:"created_at"`
	Merge_request_to_resolve_discussions_of int      `json:"merge_request_to_resolve_discussions_of"`
	Discussion_to_resolve                   string   `json:"discussion_to_resolve"`
	Iid                                     int      `json:"iid"`
	Description                             string   `json:"description"`
	Assignee_ids                            []int    `json:"assignee_ids"`
	Assignee_id                             int      `json:"assignee_id"`
	Milestone_id                            int      `json:"milestone_id"`
	Labels                                  []string `json:"labels"`
	Add_labels                              []string `json:"add_labels"`
	Remove_labels                           []string `json:"remove_labels"`
	Due_date                                string   `json:"due_date"`
	Confidential                            bool     `json:"confidential"`
	Discussion_locked                       bool     `json:"discussion_locked"`
	Issue_type                              string   `json:"issue_type"`
	Weight                                  int      `json:"weight"`
	Epic_id                                 int      `json:"epic_id"`
	Epic_iid                                int      `json:"epic_iid"`
}

type Gitlab_Milestone struct {
	Id          string `json:"id"`
	Iid         string `json:"iid"`
	Project_id  string `json:"project_id"`
	Group_id    string `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Due_date    string `json:"due_date"`
	Start_date  string `json:"start_date"`
	Expired     string `json:"expired"`
	Web_url     string `json:"web_url"`
}

type Gitlab_Custom_attributes struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Gitlab_Assignees struct {
	Id                int                        `json:"id"`
	Username          string                     `json:"username"`
	Public_email      string                     `json:"public_email"`
	Name              string                     `json:"name"`
	State             string                     `json:"state"`
	Locked            bool                       `json:"locked"`
	Avatar_url        string                     `json:"avatar_url"`
	Avatar_path       string                     `json:"avatar_path"`
	Custom_attributes []Gitlab_Custom_attributes `json:"custom_attributes"`
	Web_url           string                     `json:"web_url"`
}

type Gitlab_Author struct {
	Id                int                        `json:"id"`
	Username          string                     `json:"username"`
	Public_email      string                     `json:"public_email"`
	Name              string                     `json:"name"`
	State             string                     `json:"state"`
	Locked            bool                       `json:"locked"`
	Avatar_url        string                     `json:"avatar_url"`
	Avatar_path       string                     `json:"avatar_path"`
	Custom_attributes []Gitlab_Custom_attributes `json:"custom_attributes"`
	Web_url           string                     `json:"web_url"`
}

type Gitlab_Closed_by struct {
	Id                int                        `json:"id"`
	Username          string                     `json:"username"`
	Public_email      string                     `json:"public_email"`
	Name              string                     `json:"name"`
	State             string                     `json:"state"`
	Locked            bool                       `json:"locked"`
	Avatar_url        string                     `json:"avatar_url"`
	Avatar_path       string                     `json:"avatar_path"`
	Custom_attributes []Gitlab_Custom_attributes `json:"custom_attributes"`
	Web_url           string                     `json:"web_url"`
}

type Gitlab_Iteration struct {
	Id          string `json:"id"`
	Iid         string `json:"iid"`
	Sequence    string `json:"sequence"`
	Group_id    string `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Start_date  string `json:"start_date"`
	Due_date    string `json:"due_date"`
	Web_url     string `json:"web_url"`
}

type Create_Gitlab_Issue_Response struct {
	Id                   int              `json:"id"`
	Iid                  int              `json:"iid"`
	Project_id           int              `json:"project_id"`
	Title                string           `json:"title"`
	Description          string           `json:"description"`
	State                string           `json:"state"`
	Created_at           string           `json:"created_at"`
	Updated_at           string           `json:"updated_at"`
	Closed_at            string           `json:"closed_at"`
	Closed_by            Gitlab_Closed_by `json:"closed_by"`
	Labels               string           `json:"labels"`
	Milestone            Gitlab_Milestone `json:"milestone"`
	Assignees            Gitlab_Assignees `json:"assignees"`
	Author               Gitlab_Author    `json:"author"`
	Type                 string           `json:"type"`
	Assignee             Gitlab_Assignees `json:"assignee"`
	User_notes_count     string           `json:"user_notes_count"`
	Merge_requests_count string           `json:"merge_requests_count"`
	Upvotes              string           `json:"upvotes"`
	Downvotes            string           `json:"downvotes"`
	Due_date             string           `json:"due_date"`
	Confidential         bool             `json:"confidential"`
	Discussion_locked    bool             `json:"discussion_locked"`
	Issue_type           string           `json:"issue_type"`
	Web_url              string           `json:"web_url"`
	Time_stats           struct {
		Time_estimate          int    `json:"time_estimate"`
		Total_time_spent       int    `json:"total_time_spent"`
		Human_time_estimate    string `json:"human_time_estimate"`
		Human_total_time_spent string `json:"human_total_time_spent"`
	} `json:"time_stats"`
	Task_completion_status string `json:"task_completion_status"`
	Weight                 string `json:"weight"`
	Blocking_issues_count  string `json:"blocking_issues_count"`
	Has_tasks              string `json:"has_tasks"`
	Task_status            string `json:"task_status"`
	Links                  struct {
		Self                   string `json:"self"`
		Notes                  string `json:"notes"`
		Award_emoji            string `json:"award_emoji"`
		Project                string `json:"project"`
		Closed_as_duplicate_of string `json:"closed_as_duplicate_of"`
	} `json:"_links"`
	References struct {
		Short    string `json:"short"`
		Relative string `json:"relative"`
		Full     string `json:"full"`
	} `json:"references"`
	Severity              string `json:"severity"`
	Subscribed            string `json:"subscribed"`
	Moved_to_id           string `json:"moved_to_id"`
	Imported              string `json:"imported"`
	Imported_from         string `json:"imported_from"`
	Service_desk_reply_to string `json:"service_desk_reply_to"`
	Epic_iid              string `json:"epic_iid"`
	Epic                  struct {
		Id                       string `json:"id"`
		Iid                      string `json:"iid"`
		Title                    string `json:"title"`
		Url                      string `json:"url"`
		Group_id                 string `json:"group_id"`
		Human_readable_end_date  string `json:"human_readable_end_date"`
		Human_readable_timestamp string `json:"human_readable_timestamp"`
	} `json:"epic"`
	Iteration     Gitlab_Iteration `json:"iteration"`
	Health_status string           `json:"health_status"`
}

type Get_Gitlab_Issue_Response struct {
	Id                   int              `json:"id"`
	Iid                  int              `json:"iid"`
	Project_id           int              `json:"project_id"`
	Title                string           `json:"title"`
	Description          string           `json:"description"`
	State                string           `json:"state"`
	Created_at           string           `json:"created_at"`
	Updated_at           string           `json:"updated_at"`
	Closed_at            string           `json:"closed_at"`
	Closed_by            Gitlab_Closed_by `json:"closed_by"`
	Labels               string           `json:"labels"`
	Milestone            Gitlab_Milestone `json:"milestone"`
	Assignees            Gitlab_Assignees `json:"assignees"`
	Author               Gitlab_Author    `json:"author"`
	Type                 string           `json:"type"`
	Assignee             Gitlab_Assignees `json:"assignee"`
	User_notes_count     string           `json:"user_notes_count"`
	Merge_requests_count string           `json:"merge_requests_count"`
	Upvotes              string           `json:"upvotes"`
	Downvotes            string           `json:"downvotes"`
	Due_date             string           `json:"due_date"`
	Confidential         bool             `json:"confidential"`
	Discussion_locked    bool             `json:"discussion_locked"`
	Issue_type           string           `json:"issue_type"`
	Web_url              string           `json:"web_url"`
	Time_stats           struct {
		Time_estimate          int    `json:"time_estimate"`
		Total_time_spent       int    `json:"total_time_spent"`
		Human_time_estimate    string `json:"human_time_estimate"`
		Human_total_time_spent string `json:"human_total_time_spent"`
	} `json:"time_stats"`
	Task_completion_status string `json:"task_completion_status"`
	Weight                 string `json:"weight"`
	Blocking_issues_count  string `json:"blocking_issues_count"`
	Has_tasks              string `json:"has_tasks"`
	Task_status            string `json:"task_status"`
	Links                  struct {
		Self                   string `json:"self"`
		Notes                  string `json:"notes"`
		Award_emoji            string `json:"award_emoji"`
		Project                string `json:"project"`
		Closed_as_duplicate_of string `json:"closed_as_duplicate_of"`
	} `json:"links"`
	References struct {
		Short    string `json:"short"`
		Relative string `json:"relative"`
		Full     string `json:"full"`
	} `json:"references"`
	Severity              string `json:"severity"`
	Subscribed            string `json:"subscribed"`
	Moved_to_id           string `json:"moved_to_id"`
	Imported              string `json:"imported"`
	Imported_from         string `json:"imported_from"`
	Service_desk_reply_to string `json:"service_desk_reply_to"`
	Epic_iid              string `json:"epic_iid"`
	Epic                  struct {
		Id                       string `json:"id"`
		Iid                      string `json:"iid"`
		Title                    string `json:"title"`
		Url                      string `json:"url"`
		Group_id                 string `json:"group_id"`
		Human_readable_end_date  string `json:"human_readable_end_date"`
		Human_readable_timestamp string `json:"human_readable_timestamp"`
	} `json:"epic"`
	Iteration     Gitlab_Iteration `json:"iteration"`
	Health_status string           `json:"health_status"`
}

func (c Create_Gitlab_Issue) Create(title, description string) error {

	c.Title = title
	c.Description = description
	return nil
}

func Make_GitLab_Issue(title, description string) error {
	var newGitlabRequest Create_Gitlab_Issue

	newGitlabRequest.Title = title
	newGitlabRequest.Description = description

	GitlabCredentials, err := genericGitRequest()
	if err != nil {
		return err
	}

	// Convert the struct into JSON using the tags and Marshal
	jsonData, err := json.Marshal(newGitlabRequest)
	if err != nil {
		return err
	}

	// Convert the JSON into bytes
	requestBody := bytes.NewBuffer(jsonData)

	var GitlabID string

	// Make the request
	// /api/v4/projects/{id}/issues
	request, err := http.NewRequest("POST", fmt.Sprintf("https://gitlab"+""+"/api/v4/projects/%s/issues", GitlabID), io.Reader(requestBody))
	if err != nil {
		fmt.Printf("Error making the HTTP request %s\n", err)
		return err
	}

	// Set the required headers
	// "PRIVATE-TOKEN: <your_access_token>"
	request.Header.Set("PRIVATE-TOKEN: ", GitlabCredentials.Token)

	// Make a new client
	client := http.Client{}

	// Complete the request - Client.Do because the http.NewRequest handles the method
	req, err := client.Do(request)
	if err != nil {
		return err
	}

	if req.StatusCode != http.StatusOK && req.StatusCode != http.StatusCreated {
		fmt.Println(req.Body)
		return fmt.Errorf("the response was not positive, %d", req.StatusCode)
	}

	fmt.Printf("The response was: %s, %s\n", req.Status, HTTPStatusResponseMeanings[req.Status])

	return nil

}

func Get_Gitlab_Issues(passedFromCLI bool) {

}

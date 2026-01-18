# go-repoflow (Go)
<p align="center">
<img width="400" src="doc/images/go-repoflow.png" alt="Mascot/Icon" title="Mascot/Icon" />
</p>

A developer-focused CLI that bridges local git workflows and GitHub automation. 

## 🚀 Features

- Finds all the TODO lines in the current folder 
- Finds all the open issues in your github - using git remote 
- Checks to see whether or not the issue is in github 
    - If it is not on GitHub in will add a issue number to the start of the todo line
    - If it is on GitHub it will ignore the issue 
- Visualize commit activity across all git repositories in subdirectories, aggregated into a single terminal calendar view.
- Tag managment, create, list, and increment semantic version tags with minimal friction.
- Instantly open the remote repository in your browser (GitHub supported) for pull requests and issue URLs.
- Clone all public repositories for a given GitHub user or organization into a temporary workspace.
- Scan all subdirectories (one level deep) and report repositories with uncommitted or unpushed changes.

## 🛠️ Prerequisites

- [Go](https://golang.org/dl/) installed (version 1.16+ recommended)
- A github token for the repository with permission to read / edit issues 
    - [Github Documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)

        1. Verify your email address, if it hasn't been verified yet.
        2. In the upper-right corner of any page on GitHub, click your profile photo, then click settings.
        3. In the left sidebar, click Developer settings.
        4. In the left sidebar, under Personal access tokens, click Fine-grained tokens.
        5. Click Generate new token.
        6. Under Token name, enter a name for the token.
        7. Under Expiration, select an expiration for the token. Infinite lifetimes are allowed but may be blocked by a maximum lifetime policy set by your organization or enterprise owner. For more information, See Enforcing a maximum lifetime policy for personal access tokens.
            i. Optionally, under Description, add a note to describe the purpose of the token.
            ii. Optionally, if the resource owner is an organization that requires approval for fine-grained personal access tokens, below the resource owner, in the box, enter a justification for the request.
        9. Under Repository access, select which repositories you want the token to access. You should choose the minimal repository access that meets your needs. Tokens always include read-only access to all public repositories on GitHub.
        10. If you selected Only select repositories in the previous step, under the Selected repositories dropdown, select the repositories that you want the token to access.
        11. Under Permissions, select which permissions to grant the token. Depending on which resource owner and which repository access you specified, there are repository, organization, and account permissions. You should choose the minimal permissions necessary for your needs.

## 📁 Setup

1. Clone this repository:

   ```bash
   git clone https://github.com/jonathon-chew/go-repoflow.git
   cd go-repoflow 
   ```

2. Compile the script:


   ```bash
    go build -o repoflow ./cmd/repoflow/
   ```

3. Install the script:

    ```bash
   go install`
   ```

OR

1. Go install

    ```bash
    go install github.com/jonathon-chew/go-repoflow/cmd/rf@latest
    ```

## 📂 Output

This will make Github issues for you automatically and edit your codebase - just the todo line, to save the number of the issue for easily finding which issue is the right issue.

## 🧠 Notes

This is inspired by the project here: https://github.com/tsoding/snitch

## 📜 License

This project is licensed under the MIT License. See the LICENSE file for details.

### 🖌️ Attribution

The Go Gopher was originally designed by [Renee French](https://reneefrench.blogspot.com/).  
Used under the [Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/).  

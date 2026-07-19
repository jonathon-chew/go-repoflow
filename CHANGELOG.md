# CHANGELOG

	## v0.10.0 

		### NEW
		1. first implimentation of a written out change log
		1. first implimnetation of writing it out to a file
		### UPDATES
		1. tidying up some formatting of the change log
		1. testing changes - placement of change log for github parsing and how it handles it
		### MISC
		1. Merge pull request #76 from jonathon-chew/changelog

	## v0.9.2 

		### NEW
		1. inital implimentation of parsing git log history to make a change log
		### UPDATES
		1. having a file or folder that is not git enabled no longer stops the check of the directory
		1. README clean up

	## v0.9.1 

		### NEW
		1. starting to impliment DEBUG mode; update: currently inconsistent when removing issues from file or from github.
		### UPDATES
		1. upgrading Aphrodite so only outputs ANSI if output is a terminal
		1.remove MORE debug logs

	## v0.9.0 

		### DELETES
		1. debug logs

	## v0.8.2 

		### NEW
		1. issue update struct as assignee required when sending request, closes issues on github if line not in repo - edited todos will close the old and create a new one

	## v0.8.1 

		### UPDATES
		1. document new  subcommands; fix: default  behaviour; test: add dry-run coverage for ProcessTodosInRepo; refactor: avoid GitHub calls in dry-run mode

	## v0.8.0 

		### NEW
		1. added the ability to just list local todos - WITHOUT updating github; refactor: modify file logic out of main into it's own path;

	## v0.7.24 

		### UPDATES
		1. now looks for all files and folders from the .git root

	## v0.7.23 

		### NEW
		1. default message and incriment loop on subarguemnts refactor: help command

	## v0.7.22 

		### UPDATES
		1. CLI with no error now returns exit status 0 not 1

	## v0.7.21 

		### REFACTORS
		1. main function to return an int if failed

	## v0.7.20, tag: v0.7.19 

		### UPDATES
		1. remove true as a possible false command, as this is doing what it is supposed to!

	## v0.7.18 

		### NEW
		1. force subcommand to force tag update and push without asking

	## v0.7.17 

		### UPDATES
		1. .gitignore
		1. moving git from internal to pkg as could be used by other programmes
		### REFACTORS
		1. Utils into git pkg so it can be a package and not just internal tools
		### MISC
		1. removing my versions of docs from the repo

	## v0.7.16 

		### NEW
		1. help file to contain  all the help information, updated with the name of the command and the first version of the case statement, followed by the pertinent information about it
		### UPDATES
		1. sort the output to be in order

	## v0.7.15 

		### NEW
		1. added the ability to add tags from the command line

	## v0.7.14 

		### UPDATES
		1. adding the repostats options to the help menu as it's an option but not included in the help menu as options

	## v0.7.13 

		### UPDATES
		1. the readme image

	## v0.7.12 

		### UPDATES
		1. the readme image

	## v0.7.11 

		### MISC
		1. Updating the README to have gifs, and not break down the documentation just link to it

	## v0.7.10 

		### NEW
		1. oneline command with basic implimentation

	## v0.7.9 

		### REFACTORS
		1. contact github into it's own function taking in a generic struct

	## v0.7.8 

		### NEW
		1. function for basic repo stats, to be built on later; update: move genericGitRequest to github as this is a github specific request, even thoguh it has gitlab in currently, this will be extracted out or to an 'external' git function file, for generic external git commands

	## v0.7.7 

		### NEW
		1. default action for unrecongised commands
		1. flag for sementically similar action
		### UPDATES
		1. version number

	## v0.7.6 

		### UPDATES
		1. migrating the repository and adjusting tests to match the new structure

	## v0.7.5 

		### REFACTORS
		1. making the folder structure match common patterns
		### DELETES
		1. remove file created accidentally

	## v0.7.4 

		### MISC
		1. new lines in help command where they should be; update: README.md file to list the latest functionality

	## v0.7.0 

		### NEW
		1. Add a total commit comunter
		### UPDATES
		1. renaming the packages and folder to trying and fix the collision issue
		### MISC
		1. rename utils directory to Utils to match import paths (case-sensitive fix)
		1. correct GetLatestTag to return single latest version instead of all versions
		1. second part of test to rename utils folder on github
		1. renameing hte utils folder to force the rename to lower case

	## v0.6.1 

		### NEW
		1. adding a new feature to clone public github repsoitories
		1. adding the ability to clone all public repositories into a tmp directory

	## v0.6.0 

		### UPDATES
		1. months line up better at the end of the lines, and most recent month at the bottom

	## v0.5.4 

		### NEW
		1. feat: adding the ability to see git commit history for the last year in the terminal
		### DELETES
		1. remove examples folder

	## v0.5.3 

		### MISC
		1. logic out of cmd to it's own function in git.go for other uses

	## v0.5.2 

		### UPDATES
		1. help menu to include the new check command

	## v0.5.1 

		### NEW
		1. added the ability to track where local repos are ahead of branches
		1. adding the ability to check all sub-folders 1 level deep for any git updates that need pushing/pulling

	## v0.5.0 

		### UPDATES
		1. adding 401 error code to known issues

	## v0.4.1 

		### NEW
		1. auto make new tag when none are present; refactor: Make new tag into it's own function
		### MISC
		1. Merge pull request #66 from jonathon-chew/GitLab

	## v0.4.0 

		### UPDATES
		1. pulling changes from main
		### DELETES
		1. extra line when calling command line functions as to make it more useful and integrate with other commands like when getting latest tag

	## v0.3.1 

		### NEW
		1. default command if the user entry is not patch/minor/major

	## v0.3.0 

		### NEW
		1. message with release version on incrimenting the tag
		1. github specific functions move to their own file, as well as structus; refactor: gitlab/github structs to reflect their useage to stop collisions
		1. Adding gitlab implimetation

	## v0.2.6 

		### UPDATES
		1. getting issues ONLY gets open by default, all and closed are options

	## v0.2.5 

		### UPDATES
		1. remove unneeded line in stdout print

	## v0.2.4 

		### UPDATES
		1. version info in cmd
		1. removing coverage file from version control
		### REFACTORS
		1. only push the newest tag, which should reduce errors
		1. find .git folder

	## v0.2.3 

		### UPDATES
		1. version number
		### MISC
		1. checking for a .git folder needs to be the exact name in slices.Contains not substring

	## v0.2.2 

		### NEW
		1. feat to check the version number in the help file matches the actual latest tag. Using the CICD.sh will capture this and stop new versions being made!
		### UPDATES
		1. git status check in CICD so it works
		1. adding line by line comment strings
		1. help menu to include new lines from the title

	## v0.2.1 

		### NEW
		1. adding a new feat to the CICD pipe line allowing for dry run - checking, building and testing but not updating github
		### MISC
		1. issue when the git folder was found in the directory list it was reporting that it wasn't

	## v0.2.0 

		### UPDATES
		1. add new lines into helpful print outs confirming actions

	## v0.1.7 

		### NEW
		1. added the ability to open to a specific page

	## v0.1.6 

		### UPDATES
		1. git feedback to be colour coded

	## v0.1.5 

		### MISC
		1. git push tags option - this was creating the exec command but not running it!

	## v0.1.4 

		### UPDATES
		1. incrimenting git tag feature to prompt if they want to immediately push this to git remote

	## v0.1.3 

		### NEW
		1. CICD script, work in progress, update the version, run the tests, build the exe, combine doc/scripts and scripts; update: git_test move to git folder

	## v0.1.2 

		### UPDATES
		1. move open logic from CMD to git
		### MISC
		1. Merge pull request #62 from jonathon-chew/Version_Update

	## v0.1.1 

		### NEW
		1. feat to open the remote URL quickly for pull requests

	## v0.1.0 

		### MISC
		1. update the logic to work with no declared input, ask and continue

	## v0.0.4 

		### NEW
		1. feat incriment git tags for the expected structure v[0-9].[0-9].[0-9]
		1. line seperators for better reading when calling to see github issues
		1. script for finding CLI commands for ease of updating help CLI flag; delete: completed todo
		1. adding image to RAEDME.md
		1. parameter for CLI into listgithubissues to only print certain logs at certain times
		### UPDATES
		1. flags to allow -b and -t for body and title respectively
		1. README.md to inlcude credit for the origional gopher design
		### REFACTORS
		1. the help flag response to better fit with the formatting of other projects, and updated script to find flags which I need to talk about in the help command
		1. different parts of the application into it's own sections, cmd and git
		### MISC
		1. Update the help function to include the new features

	## v0.0.3 

		### NEW
		1. extra arguments for getting issues, adding in colour printing

	## v0.0.2 

		### NEW
		1. adding the helper CLI tool for support
		1. feat to say if it ran but didn't find a TODO; fix: fmt.Printf without a newLine at the end, new lines to break up the text of the found todo, the title and the body
		### REFACTORS
		1. CLI moved to it's own file, minor logic updating in the git calls to match better practice

	## v0.0.1 

		### NEW
		1. added the ability to only re-write the file if the file has been updates,, premptively adding the logic in main.go for updating old issues / removing them; refactor: Triming the Title so it doesn't have so much space,
		### REFACTORS
		1. turning Owner, Token and Repo into a struct for better readability
		### MISC
		1. cleaning up code
		1. cleaning up files
		1. updating TODO
		1. replace badly named variable with more appropriate
		1. Generic todo added
		1. Adding the ability to check version
		1. Work from the airport
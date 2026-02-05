package cli

import aphrodite "github.com/jonathon-chew/Aphrodite"

func help() {
	aphrodite.PrintBold("Cyan", "No Arguments\n")
	aphrodite.PrintColour("Green", "You can run with no arguments to check all the files in the current directory for any undocumented todos and upload them to github\n\n")

	aphrodite.PrintBold("cyan", "--all\n")
	aphrodite.PrintBold("cyan", "All\n")
	aphrodite.PrintColour("Green", "A sub command after get to specify get all issues regardless of state\n\n")

	aphrodite.PrintBold("cyan", "--check\n")
	aphrodite.PrintBold("cyan", "Check\n")
	aphrodite.PrintColour("Green", "Check all folders 1 level deep to see if there are any updates required to push/pull\n\n")

	aphrodite.PrintBold("cyan", "--clone\n")
	aphrodite.PrintBold("cyan", "Clone\n")
	aphrodite.PrintColour("Green", "Clone all public repos into a temporary directory\n\n")

	aphrodite.PrintBold("cyan", "--closed\n")
	aphrodite.PrintBold("cyan", "Closed\n")
	aphrodite.PrintColour("Green", "A sub command after get to specify get ONLY closed issues\n\n")

	aphrodite.PrintBold("cyan", "--commit-calendar\n")
	aphrodite.PrintBold("cyan", "Commit Calendar\n")
	aphrodite.PrintColour("Green", "Print to the terminal the git history activity for the last year!\n\n")

	aphrodite.PrintBold("cyan", "--get\n")
	aphrodite.PrintBold("Cyan", "Get issues\n")
	aphrodite.PrintColour("Green", "You can pass in a get flag which will List the github issues, this can be supplimented with --open and --closed to filter to show only issues with those flags\n\n\n")

	aphrodite.PrintBold("cyan", "--increment-tag\n")
	aphrodite.PrintBold("cyan", "Increment Tag\n")
	aphrodite.PrintColour("Green", "Finds the biggest version number in the format format v[number].[number].[number] and adds 1 to the major / minor / patch numbers\n\n")

	aphrodite.PrintBold("cyan", "--oneline\n")
	aphrodite.PrintBold("cyan", "One Line\n")
	aphrodite.PrintColour("Green", "Get back all issues in a single line, this will be determined by the title, body and status of the issue and work out the maximum information that can be given without relying on word wrap\n\n")

	aphrodite.PrintBold("cyan", "--open-issues\n")
	aphrodite.PrintBold("cyan", "Open Issues\n")
	aphrodite.PrintColour("Green", "Open the github page on the issues page to manage from there\n\n")

	aphrodite.PrintBold("cyan", "--open-pull\n")
	aphrodite.PrintBold("cyan", "Open Pull\n")
	aphrodite.PrintColour("Green", "Open the github page on the pull request page to manage from there\n\n")

	aphrodite.PrintBold("cyan", "--open\n")
	aphrodite.PrintBold("cyan", "Open\n")
	aphrodite.PrintColour("Green", "A sub command after get to specify get ONLY open issues, this is the default behaviour but for the sake of verbosity and doesn't return an error\n\n")

	aphrodite.PrintBold("cyan", "--repo-stats\n")
	aphrodite.PrintBold("cyan", "RepoStats\n")
	aphrodite.PrintColour("Green", "Get the repo stats from github Forks, Open Issues, Stargazer's, Watchers\n\n")

	aphrodite.PrintBold("cyan", "--set\n")
	aphrodite.PrintBold("Cyan", "Set issues\n")
	aphrodite.PrintColour("Green", "If you pass in the set flag, please pass in the title flag and body flag (in that order) to make a new issue with the relevent Title and Body\n\n")

	aphrodite.PrintBold("cyan", "--tags\n")
	aphrodite.PrintBold("cyan", "Tags\n")
	aphrodite.PrintColour("Green", "Returns the latest tag following the format v[number].[number].[number]\n\n")

	aphrodite.PrintBold("cyan", "--version\n")
	aphrodite.PrintBold("Cyan", "Version\n")
	aphrodite.PrintColour("Green", "Version Number can be passed in with the version flag\n\n")

}

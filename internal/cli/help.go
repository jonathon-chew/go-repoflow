package cli

import aphrodite "github.com/jonathon-chew/Aphrodite"

func help() {
	aphrodite.PrintBold("Cyan", "No Arguments\n")
	aphrodite.PrintColour("Green", "You can run with no arguments to check all the files in the current directory for any undocumented todos and upload them to github\n\n")

	aphrodite.PrintBold("cyan", "--check\n")
	aphrodite.PrintBold("cyan", "Check\n")
	aphrodite.PrintColour("Green", "Check all folders 1 level deep to see if there are any updates required to push/pull\n\n")

	aphrodite.PrintBold("cyan", "--Change-log\n")
	aphrodite.PrintBold("cyan", "Change Log\n")
	aphrodite.PrintColour("Green", "Write a change log for all changes and all versions\n\n")

	aphrodite.PrintBold("cyan", "--clone\n")
	aphrodite.PrintBold("cyan", "Clone\n")
	aphrodite.PrintColour("Green", "Clone all public repos into a temporary directory\n\n")

	aphrodite.PrintBold("cyan", "--commit-calendar\n")
	aphrodite.PrintBold("cyan", "Commit Calendar\n")
	aphrodite.PrintColour("Green", "Print to the terminal the git history activity for the last year!\n\n")
	aphrodite.PrintBold("cyan", "\tSubcommands:\n")
	aphrodite.PrintColour("Green", "\tnon-ansi, html, markdown, md\n\n")

	aphrodite.PrintBold("cyan", "--docs\n")
	aphrodite.PrintBold("cyan", "Documentation\n")
	aphrodite.PrintColour("Green", "Check if the right documentation are present!\n\n")

	aphrodite.PrintBold("cyan", "--doctor\n")
	aphrodite.PrintBold("cyan", "Doctor\n")
	aphrodite.PrintColour("Green", "Check if the right tools are in place!\n\n")

	aphrodite.PrintBold("cyan", "--get\n")
	aphrodite.PrintBold("Cyan", "Get issues\n")
	aphrodite.PrintColour("Green", "List GitHub issues for the current repo, or inspect local TODOs depending on subcommands.\n\n")
	aphrodite.PrintBold("cyan", "\tSubcommands:\n")
	aphrodite.PrintBold("cyan", "\t--all\n")
	aphrodite.PrintColour("Green", "\tReturn all issues regardless of state (open and closed)\n\n")
	aphrodite.PrintBold("cyan", "\t--closed\n")
	aphrodite.PrintColour("Green", "\tReturn only closed issues\n\n")
	aphrodite.PrintBold("cyan", "\t--open\n")
	aphrodite.PrintColour("Green", "\tReturn only open issues (default behaviour)\n\n")
	aphrodite.PrintBold("cyan", "\t--oneline / --one-line / -ol\n")
	aphrodite.PrintColour("Green", "\tPrint each issue on a single line, truncating to fit the terminal width\n\n")
	aphrodite.PrintBold("cyan", "\t--local\n")
	aphrodite.PrintColour("Green", "\tInstead of listing GitHub issues, scan local files for TODO comments without modifying files or creating/closing GitHub issues\n\n")

	aphrodite.PrintBold("cyan", "--increment-tag\n")
	aphrodite.PrintBold("cyan", "Increment Tag\n")
	aphrodite.PrintColour("Green", "Finds the biggest version number in the format format v[number].[number].[number] and adds 1 to the major / minor / patch numbers\n\n")

	aphrodite.PrintBold("cyan", "--open-issues\n")
	aphrodite.PrintBold("cyan", "Open Issues\n")
	aphrodite.PrintColour("Green", "Open the github page on the issues page to manage from there\n\n")

	aphrodite.PrintBold("cyan", "--open-pull\n")
	aphrodite.PrintBold("cyan", "Open Pull\n")
	aphrodite.PrintColour("Green", "Open the github page on the pull request page to manage from there\n\n")

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

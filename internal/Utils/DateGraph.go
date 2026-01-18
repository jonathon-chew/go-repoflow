package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

type CommitMap map[string]int

// Step 3: Render basic ASCII heatmap
func RenderDateGraph(commits CommitMap, option string) {
	now := time.Now()
	start := now.AddDate(0, 0, -365)

	// Iterate months from start..now (inclusive)
	first := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	last := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	const barWidth = 31 // 31 days max
	const gap = 1       // space before month label

	for m := first; !m.After(last); m = m.AddDate(0, 1, 0) {
		monthStart := m
		monthEnd := m.AddDate(0, 1, 0).Add(-24 * time.Hour)

		// Clamp month range to [start, now]
		from := monthStart
		if from.Before(start) {
			from = start
		}
		to := monthEnd
		if to.After(now) {
			to = now
		}
		if from.After(to) {
			continue
		}

		var b strings.Builder
		for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
			key := day.Format("2006-01-02")
			switch option {
			case "non-ansii":
				b.WriteString(heatCharNonAnsii(commits[key]))
			case "html":
				b.WriteString(heatCharHTML(commits[key]))
			case "md", "markdown":
				b.WriteString(heatCharMD(commits[key]))
			default:
				b.WriteString(heatChar(commits[key])) // must return 1 visible char/rune}
			}
		}

		// Pad to 31 chars (not bytes)
		line := b.String()

		// Measure *visible* width (without ANSI escapes)
		curWidth := visibleWidth(line)

		if curWidth < barWidth {
			line += strings.Repeat(" ", barWidth-curWidth)
		}

		if option == "markdown" || option == "md" {
			fmt.Println(line + strings.Repeat(" ", gap) + monthStart.Month().String() + "  \n\n")
		} else {
			fmt.Println(line + strings.Repeat(" ", gap) + monthStart.Month().String())
		}
	}
}

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleWidth(s string) int {
	// Remove ANSI escape codes, then count runes
	clean := ansiRE.ReplaceAllString(s, "")
	return utf8.RuneCountInString(clean)
}

func heatChar(count int) string {
	switch {
	case count == 0: // 0 commits in a day
		returnString, err := aphrodite.ReturnColour("Black", "_")
		if err != nil {
			return "§"
		}
		return returnString
	case count < 2: // 1-2 commits in a day
		returnString, err := aphrodite.ReturnColour("Red", "+")
		if err != nil {
			return "+"
		}
		return returnString
	case count < 5: //2-4 commits in a day
		returnString, err := aphrodite.ReturnBold("Green", "+")
		if err != nil {
			return "+"
		}
		return returnString
	case count < 10: //5-9 commits in a day
		returnString, err := aphrodite.ReturnHighIntensity("Yellow", "+")
		if err != nil {
			return "+"
		}
		return returnString
	default: // 10 or more commits in a day
		returnString, err := aphrodite.ReturnHighIntensityBackgrounds("Purple", "+")
		if err != nil {
			return "+"
		}
		return returnString
	}
}

func heatCharNonAnsii(count int) string {
	switch {
	case count == 0: // 0 commits in a day
		return "_"
	case count < 2: // 1-2 commits in a day
		return "+"
	case count < 5: //2-4 commits in a day
		return "*"
	case count < 10: //5-9 commits in a day
		return "&"
	default: // 10 or more commits in a day
		return "x"
	}
}

func heatCharHTML(count int) string {
	switch {
	case count == 0: // 0 commits in a day
		return wrapString("_", "<style='color: black;'", "</style>")
	case count < 2: // 1-2 commits in a day
		return wrapString("+", "<style='color: red'>", "</style>")
	case count < 5: //2-4 commits in a day
		return wrapString("*", "<b>", "</b>")
	case count < 10: //5-9 commits in a day
		return wrapString("&", "<i>", "</i>")
	default: // 10 or more commits in a day
		return wrapString("x", "<b><i>", "</b></i>")
	}
}

func heatCharMD(count int) string {
	switch {
	case count == 0: // 0 commits in a day
		return "_"
	case count < 2: // 1-2 commits in a day
		return "+"
	case count < 5: //2-4 commits in a day
		return wrapString("+", "**", "**")
	case count < 10: //5-9 commits in a day
		return wrapString("&", "*", "*")
	default: // 10 or more commits in a day
		return wrapString("X", "**", "**")
	}
}

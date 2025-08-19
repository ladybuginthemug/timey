package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "time"
	"strconv"

)


// loadRoutines loads routines from a markdown file at a given path.
func loadRoutines(path string) ([]Routine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routines []Routine
	scanner := bufio.NewScanner(file)

	var currentRoutine *Routine
	routineTitleRE := regexp.MustCompile(`^\d+\.\s+(.*)$`)
	timeRE := regexp.MustCompile(`^- Time:\s*(.*)$`)
	todoRE := regexp.MustCompile(`^-\s*\[([ x])\]\s*(.*)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if routineTitleRE.MatchString(line) {
			if currentRoutine != nil {
				routines = append(routines, *currentRoutine)
			}
			title := routineTitleRE.FindStringSubmatch(line)[1]
			currentRoutine = &Routine{Title: title}
		} else if timeRE.MatchString(line) && currentRoutine != nil {
			currentRoutine.Time = timeRE.FindStringSubmatch(line)[1]
		} else if todoRE.MatchString(line) && currentRoutine != nil {
			matches := todoRE.FindStringSubmatch(line)
			checked := matches[1] == "x"
			itemText := matches[2]
			currentRoutine.Checklist = append(currentRoutine.Checklist, ChecklistItem{
				Text:     itemText,
				Complete: checked,
			})
		}
	}

	if currentRoutine != nil {
		routines = append(routines, *currentRoutine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return routines, nil
}

// parseDuration parses a duration string like "5min" or "1h".
func parseDuration(s string) (time.Duration, error) {
    // First extract any numbers from the string
    var numStr string
    for _, char := range s {
        if char >= '0' && char <= '9' {
            numStr += string(char)
        }
    }

    // If no numbers found, return error
    if numStr == "" {
        return 0, fmt.Errorf("no numeric value found in duration: %s", s)
    }

    // Convert number string to int
    num, err := strconv.Atoi(numStr)
    if err != nil {
        return 0, err
    }

    // Convert string to lowercase and remove spaces for unit checking
    s = strings.ToLower(strings.TrimSpace(s))

    // Check for different time units
    switch {
    case strings.Contains(s, "hour") || strings.Contains(s, "hr"), strings.Contains(s, "h"):
        return time.Duration(num) * time.Hour, nil
    case strings.Contains(s, "minute") || strings.Contains(s, "min"), strings.Contains(s, "m"):
        return time.Duration(num) * time.Minute, nil
    case strings.Contains(s, "second") || strings.Contains(s, "sec"), strings.Contains(s, "s"):
        return time.Duration(num) * time.Second, nil
    default:
        // If no unit specified, assume minutes
        return time.Duration(num) * time.Minute, nil
    }
}

func (m *model) currentRoutine() Routine {
	if m.current >= len(m.routines) {
		return Routine{Title: "none", Time: "1m"}
	}
	return m.routines[m.current]
}

func (m *model) currentDuration() time.Duration {
	dur, err := parseDuration(m.currentRoutine().Time)
	if err != nil {
		return time.Minute
	}
	return dur
}
package main

import (
    "bufio"
    "os"
    "strings"
    "time"
	"fmt"
    "regexp"
	tea "github.com/charmbracelet/bubbletea"

)


func loadEventsCmd() tea.Cmd {
    return func() tea.Msg {
        events, err := loadEvents("events/events.md")
        return eventsLoadedMsg{events: events, err: err}
    }
}

// parseDate handles various date/time formats and returns a valid time.Time object.
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	layouts := []string{
		"2 January 2006 15:04",
		"2 January 2006",
		"2 January 15:04",
		"2 January",
		"2 Jan 2006 15:04",
		"2 Jan 15:04",
		"2 Jan 2006",
		"2 Jan",
	}

	currentYear := time.Now().Year()

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			if t.Year() == 0 {
				t = t.AddDate(currentYear, 0, 0)
			}
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse time from string: %s", dateStr)
}

func loadEvents(path string) ([]Event, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("events", os.ModePerm)
			return []Event{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var events []Event
	var currentEvent Event
	scanner := bufio.NewScanner(file)

	eventNameRegex := regexp.MustCompile(`^\d+\.\s+Event Name:\s+(.*)$`)
	timeRegex := regexp.MustCompile(`-\s+Time:\s+(.*)$`)
	repeatRegex := regexp.MustCompile(`-\s+Repeat:\s+(.*)$`)
	codePhraseRegex := regexp.MustCompile(`-\s+Code Phrase:\s*(.*)$`) // Fixed regex for optional space

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if match := eventNameRegex.FindStringSubmatch(line); len(match) > 1 {
			if currentEvent.Name != "" {
				events = append(events, currentEvent)
			}
			currentEvent = Event{Name: match[1]}
		} else if match := timeRegex.FindStringSubmatch(line); len(match) > 1 {
			timeStr := strings.TrimSpace(match[1])
			t, err := parseDate(timeStr)
			if err != nil {
				return nil, err
			}
			currentEvent.DateTime = t
		} else if match := repeatRegex.FindStringSubmatch(line); len(match) > 1 {
			repeat := strings.TrimSpace(match[1])
			currentEvent.Repeat = repeat
		} else if match := codePhraseRegex.FindStringSubmatch(line); len(match) > 1 {
			codePhrase := strings.TrimSpace(match[1])
			currentEvent.CodePhrase = codePhrase
		}
	}
	if currentEvent.Name != "" {
		events = append(events, currentEvent)
	}

	return events, scanner.Err()
}

func getNextOccurrence(e Event, now time.Time) time.Time {
	nextTime := e.DateTime
	if e.Repeat == "" {
		return e.DateTime
	}
	for nextTime.Before(now) {
		switch strings.ToLower(e.Repeat) {
		case "daily":
			nextTime = nextTime.AddDate(0, 0, 1)
		case "weekly":
			nextTime = nextTime.AddDate(0, 0, 7)
		case "monthly":
			// Adding a month can be tricky, so we'll use a library function if possible.
			// For simplicity here, we'll just add 30 days.
			nextTime = nextTime.AddDate(0, 1, 0)
		case "yearly":
			nextTime = nextTime.AddDate(1, 0, 0)
		default:
			return e.DateTime
		}
	}
	return nextTime
}

// saveEventToFile appends a new event to the events.md file.
func saveEventToFile(event Event) error {
	f, err := os.OpenFile("events/events.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	lastNumber := 0
	content, err := os.ReadFile("events/events.md")
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		eventNameRegex := regexp.MustCompile(`^(\d+)\.\s+Event Name:`)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if match := eventNameRegex.FindStringSubmatch(line); len(match) > 1 {
				num := 0
				fmt.Sscanf(match[1], "%d", &num)
				if num > lastNumber {
					lastNumber = num
				}
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n%d. Event Name: %s\n", lastNumber+1, event.Name))
	sb.WriteString(fmt.Sprintf("- Time: %s\n", event.DateTime.Format("2 January 2006 15:04")))
	if event.Repeat != "" {
		sb.WriteString(fmt.Sprintf("- Repeat: %s\n", event.Repeat))
	}
	sb.WriteString(fmt.Sprintf("- Code Phrase: %s\n", event.CodePhrase))

	_, err = f.WriteString(sb.String())
	return err
}
func formatTimeLeft(duration time.Duration) string {
    days := int(duration.Hours() / 24)
    hours := int(duration.Hours()) % 24
    minutes := int(duration.Minutes()) % 60

    if days > 0 {
        return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
    }
    return fmt.Sprintf("%dh %dm", hours, minutes)
}



package main

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/glamour"
)


// saveLog creates a markdown file with a summary of the session.
func (m *model) saveLog() {
    if m.routineFileName == "" {
        return // Do not save a log if no file was selected
    }

    // Create logging directory if it doesn't exist
    if err := os.MkdirAll("logging", 0755); err != nil {
        fmt.Printf("Error creating logging directory: %v\n", err)
        return
    }

    // Generate filename based on current date
    filename := fmt.Sprintf("logging/%s.md", time.Now().Format("Mon, 2 Jan 2006"))
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Error opening log file for writing: %v\n", err)
        return
    }
    defer file.Close()

    var logContent strings.Builder
    logContent.WriteString("\n---\n\n") // Separator for new sessions
    logContent.WriteString(fmt.Sprintf("## Session - %s\n", time.Now().Format("15:04:05")))

    workMap := make(map[string]time.Duration)
    for _, s := range m.sessions {
        workMap[s.RoutineTitle] += s.Elapsed
    }

    for _, routine := range m.routines {
        if dur, ok := workMap[routine.Title]; ok {
            logContent.WriteString(fmt.Sprintf("### %s\n", strings.Title(routine.Title)))
            logContent.WriteString(fmt.Sprintf("Time Spent: %s\n", dur.Truncate(time.Second)))
            logContent.WriteString("Checklist:\n")
            for _, item := range routine.Checklist {
                status := "[ ]"
                if item.Complete {
                    status = "[x]"
                }
                logContent.WriteString(fmt.Sprintf("- %s %s\n", status, item.Text))
            }
            logContent.WriteString("\n")
        }
    }

    if m.totalPaused > 0 {
        logContent.WriteString(fmt.Sprintf("Total Paused: %s\n", m.totalPaused.Truncate(time.Second)))
    }

    file.WriteString(logContent.String())
}



// generateSummaryMarkdown generates a markdown string of the session summary.
func (m model) generateSummaryMarkdown() string {
	workMap := make(map[string]time.Duration)
	for _, s := range m.sessions {
		workMap[s.RoutineTitle] += s.Elapsed
	}

	var out strings.Builder
	for _, routine := range m.routines {
		if dur, ok := workMap[routine.Title]; ok {
			out.WriteString(fmt.Sprintf("### %s\n", strings.Title(routine.Title)))
			out.WriteString(fmt.Sprintf("Time Spent: %s\n", dur.Truncate(time.Second)))
			out.WriteString("\nChecklist:\n")
			for _, item := range routine.Checklist {
				status := "[ ]"
				if item.Complete {
					status = "[x]"
				}
				out.WriteString(fmt.Sprintf("- %s %s\n", status, item.Text))
			}
			out.WriteString("\n")
		}
	}

	if m.totalPaused > 0 {
		out.WriteString(fmt.Sprintf("Total Paused: %s\n", m.totalPaused.Truncate(time.Second)))
	}

	return out.String()
}

func (m model) stopSession() model {
    if m.state == stateRunning {
        m.elapsed += time.Since(m.startTime)
        m.sessions = append(m.sessions, Session{
            RoutineTitle: m.currentRoutine().Title,
            Elapsed:      m.elapsed,
        })
    }
    if m.state == statePausing {
        m.totalPaused += time.Since(m.pauseStart)
    }

    m.state = stateStopped

    // Render the markdown summary once when the session stops
    // Initialize the viewport
    vp := viewport.New(m.width, 20)
    vp.Style = summaryViewportStyle
    
    const glamourGutter = 2
    glamourRenderWidth := m.width - vp.Style.GetHorizontalFrameSize() - glamourGutter

    renderer, err := glamour.NewTermRenderer(
        glamour.WithAutoStyle(),
        glamour.WithWordWrap(glamourRenderWidth),
    )
    if err != nil {
        vp.SetContent("Error rendering summary with glamour: " + err.Error())
    } else {
        summary, err := renderer.Render(m.generateSummaryMarkdown())
        if err != nil {
            vp.SetContent("Error rendering summary with glamour: " + err.Error())
        } else {
            vp.SetContent(summary)
        }
    }

    m.viewport = vp

    // Save the log after generating the summary
    m.saveLog()

    return m
}

// calculateTotalElapsed sums up the time from all recorded sessions.
func (m model) calculateTotalElapsed() time.Duration {
	var total time.Duration
	for _, s := range m.sessions {
		total += s.Elapsed
	}
	return total
}

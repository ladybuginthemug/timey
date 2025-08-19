package main

import (
    "fmt"
    "strings"
    "time"

    "github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
    switch m.state {
    case stateQuotes, stateCountdown:
        if m.state == stateQuotes {
            return renderQuotesView(m)
        }
        return renderCountdownView(m)

    case stateFilePicker:
        return renderFilePickerView(m)

    case stateRoutineView:
        return renderRoutineView(m)

    case statePausing:
        return renderPausingView(m)

    case stateReadyToStart:
        return renderReadyToStartView()

    case stateStopped:
        return renderStoppedView(m)

    case stateRunning, statePaused:
        return renderRunningView(m)

    case stateAddRoutine:
        return renderAddRoutineView(m)

    case stateAddEvent:
        return renderAddEventView(m)
        
    default:
        return "Unknown state"
    }
}


func renderQuotesView(m model) string {
	quote := getRandomQuote(m.quotes)
    
    return lipgloss.NewStyle().
        Padding(5).
        Align(lipgloss.Center).
        Render(
            lipgloss.NewStyle().
                Bold(true).
                Foreground(lipgloss.Color("69")).
                Render(quote.Text) +
                "\n" +
                lipgloss.NewStyle().
                    Foreground(lipgloss.Color("241")).
                    Render("- " + quote.Author) +
                "\n\n" +
                controlsStyle.Render("\nhelp • ← → : change view • l: list routines • a: add routine • e: add event • q: quit"),
        )
}

func renderCountdownView(m model) string {
    h := int(m.countdownRemaining.Hours())
    mn := int(m.countdownRemaining.Minutes()) % 60
    s := int(m.countdownRemaining.Seconds()) % 60
    timeStr := fmt.Sprintf("%02d:%02d:%02d", h, mn, s)

    // Build events section
    if m.loading {
    return fmt.Sprintf("Loading events... %s", m.countdownSpinner.View())
    }
    if m.err != nil {
        return fmt.Sprintf("Error: %v\n", m.err)
    }
    var eventsStr strings.Builder

    eventsStr.WriteString("\n\n")
    eventsStr.WriteString("Events ")
	eventsStr.WriteString("\n\n")

	if len(m.events) == 0 {
		eventsStr.WriteString("No events found. Press 'a' to add one.\n")
	} else {
		for _, event := range m.events {
			now := time.Now()
			nextOccurrence := getNextOccurrence(event, now)
			duration := time.Until(nextOccurrence)

			displayName := event.Name
			if event.CodePhrase != "" {
				displayName = event.CodePhrase
			}

			var countdown string
			days := int(duration.Hours() / 24)
			hours := int(duration.Hours()) % 24
			minutes := int(duration.Minutes()) % 60
			seconds := int(duration.Seconds()) % 60

			if days > 0 {
				countdown = fmt.Sprintf("%d d ", days)
			} else if hours > 0 {
				countdown = fmt.Sprintf("%d h %d m ", hours, minutes)
			} else if minutes > 0 {
				countdown = fmt.Sprintf("%d m %d s", minutes, seconds)
			} else if seconds > 0 {
				countdown = fmt.Sprintf("%d s", seconds)
			} else {
				countdown = fmt.Sprintf("%d d ago", days)
			}

			eventLine := fmt.Sprintf("%s: %s", eventNameStyle.Render(displayName), countdown)
			eventsStr.WriteString("• " + eventLine + "\n")
		}
	}


    return fmt.Sprintf("\n%s \n     %s time left today: %s\n%s %s\n",
        eventSeparatorStyle.Render(m.countdownGreetText),
        m.countdownSpinner.View(),
        styled.Render(timeStr),
        eventsStr.String(),
        controlsStyle.Render("\nhelp • ← → : switch view • l: list routines • a: add routine • e: add event • q: quit"),)
}

func renderFilePickerView(m model) string {
    listWidth := m.width / 2
    contentWidth := m.width - listWidth

    if m.viewport.Width == 0 || m.viewport.Height == 0 {
        m.viewport.Width = contentWidth - m.viewport.Style.GetHorizontalFrameSize()
        m.viewport.Height = m.fileList.Height() - m.viewport.Style.GetVerticalFrameSize()
    }

    leftPane := lipgloss.NewStyle().Width(listWidth).Render(m.fileList.View())
    rightPaneContent := m.viewport.View() + summaryHelpStyle("\nhelp • ↑/↓ : scroll • enter: select • q: back\n")
    rightPane := lipgloss.NewStyle().Width(contentWidth).Render(rightPaneContent)

    return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func renderRoutineView(m model) string {
    listWidth := m.width / 2
    contentWidth := m.width - listWidth

    if m.viewport.Width == 0 || m.viewport.Height == 0 {
        m.viewport.Width = contentWidth - m.viewport.Style.GetHorizontalFrameSize()
        m.viewport.Height = m.fileList.Height() - m.viewport.Style.GetVerticalFrameSize()
    }

    leftPane := lipgloss.NewStyle().Width(listWidth).Render(m.fileList.View())
    rightPaneContent := m.viewport.View() + summaryHelpStyle("\n help • ↑/↓ : scroll • enter: start routine • q: back to list\n")
    rightPane := lipgloss.NewStyle().Width(contentWidth).Render(rightPaneContent)

    return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func renderPausingView(m model) string {
    return fmt.Sprintf("\n\n%s Paused: %s\n\n%s",
        m.spinner.View(),
        pauseTimerStyle.Render(time.Since(m.pauseStart).Truncate(time.Second).String()),
        pausedControlsStyle.Render("s: resume • q: quit"))
}

func renderReadyToStartView() string {
    return "\n\nReady to start the next one?\n\ny: yes  •  n: no\n"
}

func renderStoppedView(m model) string {
    return m.viewport.View() + summaryHelpStyle("\nhelp • ↑/↓: scroll • q: menu\n")
}

func renderRunningView(m model) string {
    var b strings.Builder
    r := m.currentRoutine()
    dur := m.currentDuration()

    b.WriteString(routineTitleStyle.Render(r.Title) + "\n")

    currentRoutineElapsed := getCurrentRoutineElapsed(m)
    percent := getProgressPercentage(currentRoutineElapsed, dur)
    b.WriteString(renderProgressBar(m, percent) + "\n")

    if m.state == statePaused {
        b.WriteString(fmt.Sprintf("Paused: %s\n", time.Since(m.pauseStart).Truncate(time.Second)))
    } else {
        b.WriteString("\n")
    }

    b.WriteString(fmt.Sprintf("Progress: %s / %s\n\n",
        currentRoutineElapsed.Truncate(time.Second), dur))

    renderChecklist(&b, r.Checklist, m.selectedTodo)
    b.WriteString(controlsStyle.Render("\nhelp • s: start • p: pause • n: next • b: back • ↑/↓: select • space: toggle todo • q: quit\n"))

    return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func renderAddRoutineView(m model) string {
    var s strings.Builder
    s.WriteString(m.viewport.View())
    s.WriteString("\n\n")
    s.WriteString(m.textInput.View())

    if m.builderStage == stageDone {
        s.WriteString(focusedStyle.Render("\n Routine saved to file: " + m.filename + " ]\n"))
    }

    return s.String()
}


func getCurrentRoutineElapsed(m model) time.Duration {
    if m.state == stateRunning {
        return m.elapsed + time.Since(m.startTime)
    }
    return m.elapsed
}

func getProgressPercentage(elapsed, duration time.Duration) float64 {
    percent := float64(elapsed) / float64(duration)
    if percent > 1.0 {
        return 1.0
    }
    return percent
}

func renderProgressBar(m model, percent float64) string {
    progressBar := m.progress.ViewAs(percent)
    return fmt.Sprintf("%s  %.0f%%", progressBar, percent*100)
}

func renderChecklist(b *strings.Builder, checklist []ChecklistItem, selectedTodo int) {
    for i, item := range checklist {
        checked := " "
        if item.Complete {
            checked = "x"
        }
        itemString := fmt.Sprintf("[ %s ] %s", checked, item.Text)
        if i == selectedTodo {
            b.WriteString(focusedChecklistStyle.Render(itemString))
        } else {
            b.WriteString(checklistStyle.Render(itemString))
        }
        b.WriteString("\n")
    }
}


func renderAddEventView(m model) string {
    var s strings.Builder
    s.WriteString(m.eventViewport.View())
    s.WriteString("\n\n")
    s.WriteString(m.eventTextInput.View())
    if m.eventBuilderStage == eventStageDone {
        s.WriteString(focusedStyle.Render("\n Event saved! Press 'enter' to add another or 'q' to quit.\n"))
    }
    return s.String()
}
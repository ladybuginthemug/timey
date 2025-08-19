package main

import (

	"fmt"

	"os"
	"path/filepath"

	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/progress"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/glamour"
    "github.com/charmbracelet/lipgloss"

)


// Init initializes the application. It returns a command to be executed.
// This method is required by the tea.Model interface.
func (m model) Init() tea.Cmd {
	if m.state == stateAddRoutine {
		return textinput.Blink
	}
	// If starting in countdown state, initiate countdown commands
	if m.state == stateCountdown {
		return tea.Batch(m.countdownSpinner.Tick, tea.Every(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}))
	}
	if m.state == stateAddEvent {
			return textinput.Blink

	}
	return nil
}


// newAppModel initializes the entire application model.
func newAppModel() (model, error) {
	// File picker setup
	files, err := os.ReadDir("routines")
	if err != nil {
		// If routines directory doesn't exist, create it.
		if os.IsNotExist(err) {
			os.MkdirAll("routines", os.ModePerm)
		} else {
			return model{}, fmt.Errorf("could not read 'routines' directory: %w", err)
		}
	}
	var items []list.Item
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			displayName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			displayName = strings.ReplaceAll(displayName, "_", " ")
			items = append(items, fileItem{
				fileName:    file.Name(),
				displayName: displayName,
			})
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a Routine File"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowHelp(false)

	p := progress.New(progress.WithDefaultGradient())
	s := spinner.New(spinner.WithSpinner(spinner.Jump))
	s.Style = spinnerStyle

	// Viewport for routine content view and summary
	vp := viewport.New(0, 0)
	vp.Style = summaryViewportStyle
	vp.SetContent("")
	
	// Routine builder text input setup
	ti := textinput.New()
	ti.Placeholder = "Routine Title (then Enter)"
	ti.Prompt = focusedStyle.Render(ti.Placeholder) + " "
	ti.Cursor.Style = focusedStyle

	// Glamour renderer for routine builder viewport
	builderRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(30), // Default width, will be updated by WindowSizeMsg
	)
	if err != nil {
		return model{}, err
	}

	// Countdown setup
	now := time.Now()
	sCountdown := spinner.New()
	sCountdown.Spinner = spinner.Dot
	sCountdown.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Load quotes
	quotes, err := loadQuotes("quotes/quotes.md")
	if err != nil {
		// If quotes directory doesn't exist, create it
		if os.IsNotExist(err) {
			os.MkdirAll("quotes", os.ModePerm)
			quotes = []Quote{} // Empty quotes list
		} else {
			return model{}, fmt.Errorf("could not read quotes: %w", err)
		}
	}

	// Load events
    events, err := loadEvents("events/events.md")
    if err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll("events", os.ModePerm)
            events = []Event{}
        } else {
            return model{}, fmt.Errorf("could not read events: %w", err)
        }
    }

	evp := viewport.New(0, 0)
	evp.Style = summaryViewportStyle
    // event builder text input setup
	eti := textinput.New()
	eti.Placeholder = "Event Name"
	eti.Prompt = focusedStyle.Render(ti.Placeholder) + " "
	eti.Cursor.Style = focusedStyle

	// Glamour renderer for event builder viewport
	eventBuilderRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(30), // Default width, will be updated by WindowSizeMsg
	)
	if err != nil {
		return model{}, err
	}

	

	return model{
		state:              stateCountdown, // Start with the countdown
		fileList:           l,
		progress:           p,
		spinner:            s,
		viewport:           vp,
		textInput:          ti,
		builderStage:       stageTitle,
		renderer:           builderRenderer,
		countdownRemaining: timeLeftToday(),
		countdownSpinner:   sCountdown,
		countdownGreetText: greet(now),
		quotes:             quotes,
		events:             events,
		eventViewport:           evp,
        eventTextInput:          eti,
        eventBuilderStage:       eventStageName,
        eventRenderer:           eventBuilderRenderer,   
	}, nil
}


// updatePaneSizes adjusts the dimensions of the file list and viewport based on current window size.
func (m *model) updatePaneSizes() {
    if m.width == 0 || m.height == 0 {
        return
    }

    if m.state == stateFilePicker || m.state == stateRoutineView {
        listWidth := m.width / 2
        contentWidth := m.width - listWidth

        top, right, bottom, left := m.fileList.Styles.Title.GetPadding()
        m.fileList.SetSize(listWidth-left-right, m.height-top-bottom-1)
        
        m.viewport.Width = contentWidth - m.viewport.Style.GetHorizontalFrameSize()
        m.viewport.Height = m.fileList.Height() - m.viewport.Style.GetVerticalFrameSize()
    } else if m.state == stateStopped {
        m.viewport.Width = m.width - m.viewport.Style.GetHorizontalFrameSize()
        m.viewport.Height = m.height - m.viewport.Style.GetVerticalFrameSize()
    } else if m.state == stateAddRoutine {
        // Adjust viewport for routine builder
        builderViewportWidth := m.width - m.viewport.Style.GetHorizontalFrameSize()
        builderViewportHeight := m.height - lipgloss.Height(m.textInput.View()) - lipgloss.Height(separatorStyle.Render("")) - 2 // Account for input and separator
        m.viewport.Width = builderViewportWidth
        m.viewport.Height = builderViewportHeight

        // Adjust renderer word wrap for builder viewport
        m.renderer, _ = glamour.NewTermRenderer(
            glamour.WithAutoStyle(),
            glamour.WithWordWrap(m.viewport.Width - 2), // Account for glamour's internal gutter
        )
	 } else if m.state == stateAddEvent {
        // Adjust viewport for routine builder
        builderViewportWidth := m.width - m.viewport.Style.GetHorizontalFrameSize()
        builderViewportHeight := m.height - lipgloss.Height(m.eventTextInput.View()) - lipgloss.Height(separatorStyle.Render("")) - 2 
        m.eventViewport.Width = builderViewportWidth
        m.eventViewport.Height = builderViewportHeight

        // Adjust renderer word wrap for builder viewport
        m.eventRenderer, _ = glamour.NewTermRenderer(
            glamour.WithAutoStyle(),
            glamour.WithWordWrap(m.viewport.Width - 2), // Account for glamour's internal gutter
        )
    }
    }


func main() {
	model, err := newAppModel()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(&model).Run(); err != nil { // Pass pointer to model
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

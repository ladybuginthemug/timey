package main

import (
    "time"

    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/progress"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/glamour"
)


// fileItem is a custom type to hold both the raw filename and the formatted display name.
type fileItem struct {
	fileName    string
	displayName string
}

// FilterValue is required by the list.Item interface.
func (i fileItem) FilterValue() string { return i.displayName }
func (i fileItem) Title() string       { return i.displayName }
func (i fileItem) Description() string { return "" }

// ChecklistItem represents a single to-do item in a routine.
type ChecklistItem struct {
	Text     string
	Complete bool
}

// Routine holds the data for a single routine, including its title, time, and checklist.
type Routine struct {
	Title     string
	Time      string
	Checklist []ChecklistItem
}

// Session tracks the duration of each routine segment.
type Session struct {
	RoutineTitle string
	Elapsed      time.Duration
}

// A tickMsg is sent on a regular interval to update the timer.
type tickMsg time.Time

// appState defines the different states of the application.
type appState int

const (
	stateCountdown appState = iota // Start with the countdown
	stateQuotes
	stateFilePicker
	stateRoutineView
	stateRunning
	statePaused
	stateStopped
	statePausing
	stateReadyToStart
	stateAddRoutine 
	stateAddEvent
)

// stage represents the current state of the routine builder.
type stage int

const (
	stageTitle stage = iota
	stageHabit
	stageTime
	stageTodo
	stageDone
)

type Quote struct {
    Text string
    Author string
}


type Event struct {
	Name       string
	DateTime   time.Time
	Repeat     string // Optional repeat pattern (e.g., "daily", "weekly", "monthly", "yearly")
	CodePhrase string // Optional code phrase for the event
}	

type eventStage int

const (
	eventStageName eventStage = iota
	eventStageTime
	eventStageRepeat
	eventStageCodePhrase
	eventStageDone
)

type eventsLoadedMsg struct {
	events []Event
	err    error
}

type model struct {
	routines     []Routine
	current      int
	startTime    time.Time
	elapsed      time.Duration
	totalElapsed time.Duration

	pauseStart    time.Time
	totalPaused   time.Duration
	state         appState
	selectedTodo  int

	sessions        []Session
	progress        progress.Model
	spinner         spinner.Model
	width           int
	height          int 
	
	// New fields for file picker functionality
	fileList list.Model
	// holds the name of the selected file.
	routineFileName string 

	// summary view
	//  rendered summary.
	summaryRendered string 
	viewport viewport.Model

	// fields for routine builder functionality
	textInput        textinput.Model
	builderStage     stage 
	routineMarkdown  string
	currentHabit     int
	currentHabitName string
	// filename for the new routine
	filename         string 
	todosCount       int
	// renders the markdown in the routine builder.
	renderer         *glamour.TermRenderer 

	
	countdownRemaining time.Duration
	countdownSpinner   spinner.Model

	// Used to signal exit from countdown, not global quit
	countdownQuitting  bool 
	countdownGreetText string

	quotes []Quote

	events []Event 
	err               error
	loading           bool
	eventViewport     viewport.Model
	eventTextInput    textinput.Model
	eventBuilderStage eventStage
	eventMarkdown     string
	eventFilename     string
	currentEventName  string
	currentEventTime  string
	eventRenderer    *glamour.TermRenderer
	eventRepeat	 	  string // Optional repeat pattern for the event

}



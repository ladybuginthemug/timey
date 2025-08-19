package main

import ("github.com/charmbracelet/lipgloss"
	   "github.com/charmbracelet/bubbles/list"
)


// Styling constants for the UI.
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("205"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4)
	// Additional styles for the routine runner view.
	routineTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).PaddingBottom(1)
	controlsStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingTop(1)
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Used in routine builder

	// New style for checklist items to ensure consistent indentation.
	checklistStyle        = lipgloss.NewStyle().PaddingLeft(2)
	focusedChecklistStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("205"))
	spinnerStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	pauseTimerStyle       = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Align(lipgloss.Center).
				Width(20)
			

	// New styling for the paused window and summary view
	pausedControlsStyle = lipgloss.NewStyle().
				Padding(2).
				Bold(true).
				Foreground(lipgloss.Color("240"))
	
	summaryViewportStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				PaddingRight(2)
	
	summaryHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(1).
				Render
	
	// Styles for routine builder
	blurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingTop(1).PaddingBottom(1)

	// Event styles
	eventNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))  // Use a bright color for event names
	
	eventTimeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("69"))  // Softer color for time remaining

	eventSeparatorStyle = lipgloss.NewStyle().
		Width(40).
		BorderStyle(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center)

	styled = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("230")).
        Background(lipgloss.Color("27")).
        Padding(0, 1)
	
	rstyled = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("27")).
        Background(lipgloss.Color("230")).
        Padding(0, 1)
	)
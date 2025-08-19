package main

import (
    "os"
    "strings"
	"fmt"

)


// createFile sets up the filename and initial markdown for the routine.
func (m *model) createFile(title string) {
	fname := strings.ReplaceAll(title, " ", "_")
	os.MkdirAll("routines", os.ModePerm) // Ensure directory exists
	m.filename = fmt.Sprintf("routines/%s.md", fname)
	m.routineMarkdown = fmt.Sprintf("# %s\n\n", title)
}

// saveHabitHeader adds a new habit to the markdown string.
func (m *model) saveHabitHeader(habit string, timeStr string) {
	m.currentHabit++
	m.routineMarkdown += fmt.Sprintf("%d. %s\n", m.currentHabit, habit)
	m.routineMarkdown += fmt.Sprintf("- Time: %s\n", timeStr)
	m.todosCount = 0
}

// saveTodo adds a new todo item to the current habit in the markdown string.
func (m *model) saveTodo(todo string) {
	m.todosCount++
	m.routineMarkdown += fmt.Sprintf("- [ ] %s\n", todo)
}
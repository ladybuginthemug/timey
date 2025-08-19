package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/bubbles/progress"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = m.width - 4
		m.updatePaneSizes()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			switch m.state {
			case stateCountdown, stateQuotes:
				return m, tea.Quit
			case stateRunning, statePausing, statePaused:
				if m.state == stateRunning {
					m.elapsed += time.Since(m.startTime)
				}
				if m.state == statePausing {
					m.totalPaused += time.Since(m.pauseStart)
				}
				m.sessions = append(m.sessions, Session{
					RoutineTitle: m.currentRoutine().Title,
					Elapsed:      m.elapsed,
				})
				*m = m.stopSession()
				return m, nil
			default:
				m.state = stateCountdown
				m.countdownRemaining = timeLeftToday()
				return m, tea.Batch(m.countdownSpinner.Tick, tick())
			}

		case "l":
			if m.state == stateQuotes || m.state == stateCountdown {
				m.state = stateFilePicker
				m.updatePaneSizes()
				return m, nil
			}

		case "left", "right":
			if m.state == stateQuotes {
				m.state = stateCountdown
				return m, tea.Batch(m.countdownSpinner.Tick, tick())
			} else if m.state == stateCountdown {
				m.state = stateQuotes
				return m, nil
			}

		case "a":
			if m.state == stateQuotes || m.state == stateCountdown {
				m.state = stateAddRoutine
				m.textInput.Focus()
				m.routineMarkdown = ""      // reset markdown
        		m.viewport.SetContent("")   // clear viewport
				m.updatePaneSizes()
				return m, textinput.Blink
			}

		case "s":
			if m.state == statePausing || m.state == statePaused {
				m.totalPaused += time.Since(m.pauseStart)
				m.startTime = time.Now()
				m.state = stateRunning
				return m, tick()
			}

		case "e":
			if m.state == stateQuotes || m.state == stateCountdown {
				m.state = stateAddEvent
				m.eventTextInput.Reset()
				m.eventBuilderStage = eventStageName
				m.eventTextInput.Focus()
				m.eventTextInput.Placeholder = "Event Name:"
				m.eventTextInput.Prompt = focusedStyle.Render(m.eventTextInput.Placeholder) + " "
				m.eventMarkdown = ""
				m.updatePaneSizes()
				return m, textinput.Blink
			}

		case "enter":
			switch m.state {
			case stateFilePicker:
				selectedItem, ok := m.fileList.SelectedItem().(fileItem)
				if !ok {
					return m, nil
				}
				path := fmt.Sprintf("routines/%s", selectedItem.fileName)
				content, err := os.ReadFile(path)
				if err != nil {
					m.viewport.SetContent("Error reading file: " + err.Error())
					return m, nil
				}
				glamourRenderWidth := m.viewport.Width - 2
				renderer, err := glamour.NewTermRenderer(
					glamour.WithAutoStyle(),
					glamour.WithWordWrap(glamourRenderWidth),
				)
				if err != nil {
					m.viewport.SetContent("Error rendering content: " + err.Error())
				} else {
					renderedContent, err := renderer.Render(string(content))
					if err != nil {
						m.viewport.SetContent("Error rendering content: " + err.Error())
					} else {
						m.viewport.SetContent(renderedContent)
					}
				}
				m.viewport.GotoTop()
				routines, err := loadRoutines(path)
				if err != nil {
					return m, nil
				}
				m.routines = routines
				m.routineFileName = selectedItem.fileName
				m.state = stateRoutineView
				m.updatePaneSizes()
				m.elapsed = 0
				m.current = 0
				m.selectedTodo = 0
				m.sessions = []Session{}
				return m, nil

			case stateRoutineView:
				m.state = stateRunning
				m.startTime = time.Now()
				return m, tick()

			case stateReadyToStart:
				m.state = stateRunning
				m.startTime = time.Now()
				return m, tick()

			case stateAddRoutine:
				val := strings.TrimSpace(m.textInput.Value())
				if val == "" {
					return m, nil
				}
				switch m.builderStage {
				case stageTitle:
					m.createFile(val)
					m.textInput.Reset()
					m.textInput.Placeholder = "Habit name (or 'done' to finish):"
					m.textInput.Prompt = focusedStyle.Render(m.textInput.Placeholder) + " "
					m.builderStage = stageHabit
				case stageHabit:
					if strings.ToLower(val) == "done" {
						m.builderStage = stageDone
						_ = os.WriteFile(m.filename, []byte(m.routineMarkdown), 0644)
						m.textInput.Blur()
						m.state = stateCountdown
						m.countdownRemaining = timeLeftToday()
						return m, tea.Batch(m.countdownSpinner.Tick, tick())
					}
					m.currentHabitName = val
					m.textInput.Reset()
					m.textInput.Placeholder = "Time for this habit (e.g. 1h):"
					m.textInput.Prompt = focusedStyle.Render(m.textInput.Placeholder) + " "
					m.builderStage = stageTime
				case stageTime:
					m.saveHabitHeader(m.currentHabitName, val)
					m.currentHabitName = ""
					m.textInput.Reset()
					m.textInput.Placeholder = "Todo item (enter 'done' when no more todos):"
					m.textInput.Prompt = focusedStyle.Render(m.textInput.Placeholder) + " "
					m.builderStage = stageTodo
				case stageTodo:
					if strings.ToLower(val) == "done" {
						m.routineMarkdown += "\n"
						m.textInput.Reset()
						m.textInput.Placeholder = "Next Habit name (or 'done' to finish):"
						m.textInput.Prompt = focusedStyle.Render(m.textInput.Placeholder) + " "
						m.builderStage = stageHabit
					} else {
						m.saveTodo(val)
						m.textInput.Reset()
					}
				}
				return m, nil

			case stateAddEvent:
				val := strings.TrimSpace(m.eventTextInput.Value())
				if val == "" && m.eventBuilderStage != eventStageDone {
					return m, nil
				}
				switch m.eventBuilderStage {
				case eventStageName:
				    m.currentEventName = val
					m.eventMarkdown = fmt.Sprintf("# %s\n", m.currentEventName)
					m.eventTextInput.Reset()
					m.eventTextInput.Placeholder = "Event Time (e.g. 17 August 2025 15:00):"
					m.eventTextInput.Prompt = focusedStyle.Render(m.eventTextInput.Placeholder) + " "
					m.eventBuilderStage = eventStageTime
					return m, textinput.Blink

				case eventStageTime:
					m.currentEventTime = val
					m.eventMarkdown += fmt.Sprintf("- Time: %s\n", m.currentEventTime)
					m.eventTextInput.Reset()
					m.eventTextInput.Placeholder = "Repeat (daily, weekly, monthly, yearly, or 'none'):"
					m.eventTextInput.Prompt = focusedStyle.Render(m.eventTextInput.Placeholder) + " "
					m.eventBuilderStage = eventStageRepeat
					return m, textinput.Blink

				case eventStageRepeat:
					repeat := strings.TrimSpace(val)
					if strings.ToLower(repeat) == "none" {
						repeat = ""
					}
					m.eventMarkdown += fmt.Sprintf("- Repeat: %s\n", repeat)
					m.eventTextInput.Reset()
					m.eventTextInput.Placeholder = "Code Phrase (optional, or 'none'):"
					m.eventTextInput.Prompt = focusedStyle.Render(m.eventTextInput.Placeholder) + " "
					m.eventBuilderStage = eventStageCodePhrase
					m.eventFilename = repeat // Store repeat for later
					return m, textinput.Blink

				case eventStageCodePhrase:
					codePhrase := strings.TrimSpace(val)
					if strings.ToLower(codePhrase) == "none" {
						codePhrase = ""
					}
					m.eventMarkdown += fmt.Sprintf("- Code Phrase: %s\n", codePhrase)

					// Parse event time
					t, err := parseDate(m.currentEventTime)
					if err != nil {
						m.eventMarkdown += fmt.Sprintf("\n\nError parsing time: %s\nTry again.", err.Error())
						m.eventTextInput.Reset()
						m.eventBuilderStage = eventStageTime
						return m, textinput.Blink
					}

					// Save event
					newEvent := Event{
						Name:       m.currentEventName,
						DateTime:   t,
						Repeat:     m.eventFilename,
						CodePhrase: codePhrase,
					}
					if err := saveEventToFile(newEvent); err != nil {
						m.eventMarkdown += fmt.Sprintf("\n\nError saving event: %s", err.Error())
						m.eventTextInput.Reset()
						m.eventBuilderStage = eventStageName
						return m, textinput.Blink
					}

					m.eventMarkdown += "\n\nEvent saved! Press Enter to add another or 'q' to quit."
					m.eventTextInput.Reset()
					m.eventBuilderStage = eventStageDone
					return m, nil

				case eventStageDone:
					// Reset for another event or quit
					m.state = stateCountdown
					m.eventTextInput.Blur()
					return m, loadEventsCmd()
				}
	
			}
		}

		if m.state == stateRunning || m.state == statePaused {
			switch msg.String() {
			case "up":
				if m.selectedTodo > 0 {
					m.selectedTodo--
				}
			case "down":
				if m.selectedTodo < len(m.currentRoutine().Checklist)-1 {
					m.selectedTodo++
				}
			case " ":
				if len(m.currentRoutine().Checklist) > 0 && m.selectedTodo >= 0 && m.selectedTodo < len(m.currentRoutine().Checklist) {
					m.routines[m.current].Checklist[m.selectedTodo].Complete = !m.routines[m.current].Checklist[m.selectedTodo].Complete
				}
			case "p":
				if m.state == stateRunning {
					m.elapsed += time.Since(m.startTime)
					m.pauseStart = time.Now()
					m.state = statePausing
					return m, m.spinner.Tick
				}
			case "n":
				if m.current < len(m.routines)-1 {
					m.elapsed += time.Since(m.startTime)
					m.sessions = append(m.sessions, Session{
						RoutineTitle: m.currentRoutine().Title,
						Elapsed:      m.elapsed,
					})
					m.current++
					m.state = stateReadyToStart
					m.elapsed = 0
					m.selectedTodo = 0
					return m, nil
				}
				*m = m.stopSession()
				return m, nil
			case "b":
				if m.current > 0 {
					m.elapsed += time.Since(m.startTime)
					m.sessions = append(m.sessions, Session{
						RoutineTitle: m.currentRoutine().Title,
						Elapsed:      m.elapsed,
					})
					m.current--
					m.state = stateReadyToStart
					m.elapsed = 0
					m.selectedTodo = 0
					return m, nil
				}
			}
		}

		if m.state == stateReadyToStart {
			switch msg.String() {
			case "y":
				m.state = stateRunning
				m.startTime = time.Now()
				return m, tick()
			case "n":
				m.pauseStart = time.Now()
				m.state = statePausing
				return m, m.spinner.Tick
			}
		}

	case tickMsg:
		if m.state == stateCountdown {
			m.countdownRemaining = timeLeftToday()
			return m, tea.Batch(
				m.countdownSpinner.Tick,
				tea.Every(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }),
			)
		}
		if m.state == stateRunning {
			if m.elapsed+time.Since(m.startTime) >= m.currentDuration() {
				m.elapsed += time.Since(m.startTime)
				m.sessions = append(m.sessions, Session{
					RoutineTitle: m.currentRoutine().Title,
					Elapsed:      m.elapsed,
				})
				m.current++
				if m.current >= len(m.routines) {
					*m = m.stopSession()
					return m, nil
				}
				m.state = stateReadyToStart
				m.elapsed = 0
				m.selectedTodo = 0
				return m, nil
			}
			return m, tick()
		}

	case spinner.TickMsg:
		if m.state == stateCountdown {
			var sp spinner.Model
			sp, cmd = m.countdownSpinner.Update(msg)
			m.countdownSpinner = sp
			return m, cmd
		}
		if m.state == statePausing {
			var sp spinner.Model
			sp, cmd = m.spinner.Update(msg)
			m.spinner = sp
			cmds = append(cmds, cmd)
			cmds = append(cmds, m.spinner.Tick)
		}
	case eventsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.events = msg.events
		return m, tick()
	}

	// Update UI components based on state
	if m.state == stateFilePicker {
		m.fileList, cmd = m.fileList.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateRunning || m.state == statePaused || m.state == stateReadyToStart || m.state == statePausing {
		var p tea.Model
		p, cmd = m.progress.Update(msg)
		m.progress = p.(progress.Model)
		cmds = append(cmds, cmd)
	} else if m.state == stateStopped || m.state == stateRoutineView || m.state == stateAddRoutine {
		if m.state == stateAddRoutine && m.builderStage != stageDone {
			if strings.TrimSpace(m.routineMarkdown) != "" {
			renderedMarkdown, err := m.renderer.Render(m.routineMarkdown)
			if err == nil {
				m.viewport.SetContent(renderedMarkdown)
			} else {
				m.viewport.SetContent(m.routineMarkdown)
			}
		} else {
			m.viewport.SetContent("")  // No artifact on empty
		}
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
		}
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateAddEvent  && m.builderStage != stageDone {
		if strings.TrimSpace(m.eventMarkdown) != "" {
        renderedMarkdown, err := m.renderer.Render(m.eventMarkdown)
        if err == nil {
            m.eventViewport.SetContent(renderedMarkdown)
        } else {
            m.eventViewport.SetContent(m.eventMarkdown)
        }
		} else {
			m.eventViewport.SetContent("")  // No artifact on empty
		}
			m.eventTextInput, cmd = m.eventTextInput.Update(msg)
			cmds = append(cmds, cmd)
			m.eventViewport, cmd = m.eventViewport.Update(msg)
			cmds = append(cmds, cmd)
			}

	return m, tea.Batch(cmds...)
	}

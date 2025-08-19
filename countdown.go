
package main

import (
    "fmt"
    "math"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)


// tick sends a message every second to update the timer.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// greet generates a greeting based on the current hour
func greet(now time.Time) string {
	year := now.Year()
	month := int(now.Month())
	mday := now.Day()
	hour, minute := now.Hour(), now.Minute()
	weekday := int(now.Weekday())
	yday := now.YearDay()

	greeting := ""
	switch {
	case hour < 12:
		greeting += "Good morning!\n"
	case hour < 18:
		greeting += "Good afternoon!\n"
	default:
		greeting += "Good evening!\n"
	}
	greeting += fmt.Sprintf("Today is %s, %d %s %d.\n", now.Weekday(), mday, now.Month(), year)
	greeting += "\n year progress: "
	greeting += yearLeft(year, yday) + "\n"
	greeting += "\n month progress: "
	greeting += monthLeft(mday, month, year) + "\n"
	greeting += "\n week progress: "
	greeting += weekLeft(weekday) + "\n"
	greeting += "\n day progress: "
	greeting += dayLeft(hour, minute) + "\n"

	return greeting
}

// leapYear checks if a year is a leap year
func leapYear(year int) int {
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return 366
	}
	return 365
}

// progressBar returns a formatted progress bar string
func progressBar(start, end, dleft float64, length int, isHours bool) string {
	percentage := math.Floor((start / end) * 100)
	left := int(float64(length) * percentage / 100)
	right := length - left

	unit := "days"
	if isHours {
		unit = "hours"
	}

	return fmt.Sprintf("%.0f %s left\n [%s%s] %.0f%% ",
		dleft,
		unit,
		strings.Repeat("■", left),
		strings.Repeat("□", right),
		percentage,
	)
}

// yearLeft returns the progress bar for the current year
func yearLeft(year, yday int) string {
	daysInYear := leapYear(year)
	left := float64(daysInYear - yday)
	return progressBar(float64(yday), float64(daysInYear), left, 30, false)
}

// monthLeft returns the progress bar for the current month
func monthLeft(mday, month, year int) string {
	now := time.Now()
	loc := now.Location()

	var nextMonth time.Time
	if month == 12 {
		nextMonth = time.Date(year+1, time.January, 1, 0, 0, 0, 0, loc)
	} else {
		nextMonth = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, loc)
	}

	lastDay := nextMonth.AddDate(0, 0, -1).Day()
	left := float64(lastDay - mday)

	return progressBar(float64(mday), float64(lastDay), left, 30, false)
}

// weekLeft returns the progress bar for the current week
func weekLeft(weekday int) string {
	left := float64(7 - weekday)
	return progressBar(float64(weekday), 7, left, 30, false)
}

// dayLeft returns the progress bar for the current day
func dayLeft(hour, minute int) string {
	h := float64(hour)
	m := float64(minute)
	totalHours := h + (m / 60)
	left := math.Floor(24.0 - totalHours)
	return progressBar(totalHours, 24.0, left, 30, true)
}

// timeLeftToday returns time left until end of day
func timeLeftToday() time.Duration {
	now := time.Now()
	eod := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return eod.Sub(now)
}

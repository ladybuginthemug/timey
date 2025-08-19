# timey

A terminal-based daily routine and event planner built with Bubbletea.

## How to build and run

1. Clone this repository:
    ```
    git clone https://github.com/ladybuginthemug/timey.git
    cd timey
    ```

2. Build:
    ```
    go build
    ```

3. Run:
    ```
    ./timey
    ```
    Or, run directly without building:
    ```
    go run .
    ```

## Requirements

- Go 1.23 or newer

## Features

- Routine builder
- Event scheduler
- Countdown timer
- Quotes
- Markdown-based storage
---
## How Routine Structure Works

 Each routine consists of several steps (habits or activities), each with a time estimate and optional checklist items. The file is written in markdown, and the format is critical for the app to read and display it correctly.

### Structure Overview

- The routine begins with a title, prefixed by `#`.
- Each habit starts with a number and description: `N. Description`
- The time for each habit is listed below, prefixed by `- Time:`.
- Each checklist item is a markdown todo: `- [ ] item`
  
```
# Morning Productivity

1. Wake up
- Time: 10 min
- [ ] Drink water
- [ ] Make the bed

2. Exercise
- Time: 30 min
- [ ] Stretching
- [ ] Cardio
```

---
## How Event Structure Works

Each event is listed with an incrementing number, a name, time, repeat pattern (optional), and code phrase (optional). The file is written in markdown and the format is important for the app to read and display upcoming events.

### Structure Overview

- Each event starts with a number and the label `Event Name:`
- The time for the event is listed below, prefixed by `- Time:`
- The repeat pattern is optional and listed as `- Repeat:`(daily, weekly, monthly, yearly)
- The code phrase is optional and listed as `- Code Phrase:`, if you want keep the event as a secret 


```
1. Event Name: Team Standup
- Time: 20 August 2025 09:30
- Repeat: daily
- Code Phrase: 

2. Event Name: Surprise party
- Time: 31 August 2025 23:59
- Repeat: 
- Code Phrase: Secret
```

---


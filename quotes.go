package main

import (
    "bufio"
    "math/rand"
    "os"
    "regexp"
    "strings"
    "time"
)

func loadQuotes(path string) ([]Quote, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var quotes []Quote
    scanner := bufio.NewScanner(file)
    quoteRE := regexp.MustCompile(`^\d+\.\s*"([^"]+)"\s*-\s*(.+)$`)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        if matches := quoteRE.FindStringSubmatch(line); matches != nil {
            quotes = append(quotes, Quote{
                Text:   matches[1],
                Author: matches[2],
            })
        }
    }

    return quotes, scanner.Err()
}

func getRandomQuote(quotes []Quote) Quote {
    if len(quotes) == 0 {
        return Quote{
            Text:   "Time is what we want most, but what we use worst.",
            Author: "William Penn",
        }
    }
    
    rand.Seed(time.Now().UnixNano())
    return quotes[rand.Intn(len(quotes))]
}
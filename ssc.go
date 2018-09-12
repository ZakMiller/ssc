package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ogier/pflag"
	"github.com/pkg/browser"
	"github.com/mmcdole/gofeed"
)

func getTitles(count int) (string, error) {
	const (
		rssFeedURL = "http://slatestarcodex.com/feed/"
	)
	resp, err := http.Get(rssFeedURL)
	if err != nil {
		return "", fmt.Errorf("couldn't retrieve feed")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read feed")
	}
	text := string(b)

	parser := gofeed.NewParser()
	feed, err := parser.ParseString(text)
	if err != nil {
		return "", fmt.Errorf("couldn't parse feed")
	}

	var buffer bytes.Buffer
	if count == 1 {
		buffer.WriteString("The most recent article:\n")
	} else {
		buffer.WriteString(fmt.Sprintf("The %d most recent articles:\n", count))
	}

	for i := 0; i < count; i++ {
		buffer.WriteString(fmt.Sprintf("%v\n", feed.Items[i].Title))
	}

	return buffer.String(), nil
}

func main() {
	count := pflag.Int("count", 1, "The number of titles you want to retrieve. Default 1.")

	pflag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("A command is required. Try 'titles' or 'open'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "titles":
		result, err := getTitles(*count)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	case "open":
		browser.OpenURL("https://slatestarcodex.com")

	default:
		fmt.Println("Unrecognized command. Try 'titles' or 'open'")
	}
}

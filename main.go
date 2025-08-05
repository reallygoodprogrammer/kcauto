package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	
	"github.com/playwright-community/playwright-go"
)

func main() {
	perr := func(err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	var lastEpisodeFile string
	var noLastFile bool

	flag.StringVar(&lastEpisodeFile, "l", "last-episode.txt", "")
	flag.StringVar(&lastEpisodeFile, "last-episode", "last-episode.txt", "")
	flag.BoolVar(&noLastFile, "n", false, "")
	flag.BoolVar(&noLastFile, "no-write", false, "")

	flag.Usage = func() { 
		fmt.Println("usage: kcauto <url-to-episode>")
		fmt.Println("\t-l/--last-episode FILENAME\tset last episode file")
		fmt.Println("\t-n/--no-write\t\t\tdo not create last-episode file")
		fmt.Println("\t-h/--help\t\t\tdisplay this help message")
		os.Exit(0)
	}

	flag.Parse()
	userDataDir := os.Getenv("FIREFOX_PROFILE")

	url := flag.Arg(0)
	if url == "" {
		data, err := os.ReadFile(lastEpisodeFile)
		perr(err)
		url = string(data)
	}

	pw, err := playwright.Run()
	perr(err)
	defer pw.Stop()

	context, err := pw.Firefox.LaunchPersistentContext(userDataDir, playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(false),
	})
	perr(err)
	defer context.Close()

	page, err := context.NewPage()
	perr(err)
	defer page.Close()

	_, err = page.Goto(url)
	perr(err)

	for true {
		if !noLastFile {
			file, err := os.Create(lastEpisodeFile)
			perr(err)
			_, err = file.WriteString(page.URL())
			perr(err)
			file.Close()
		}

		time.Sleep(8 * time.Second)

		videoBox, err := page.QuerySelector("#player_container")
		perr(err)
		videoBox.ScrollIntoViewIfNeeded()

		iframeElement, err := page.QuerySelector("#player_container iframe")
		perr(err)
		iframe, err := iframeElement.ContentFrame()
		perr(err)

		popup, err := iframe.QuerySelector(".jwp-popup")
		perr(err)
		for popup != nil && err == nil {
			time.Sleep(3 * time.Second)
			popup, err = iframe.QuerySelector(".jwp-popup")
			perr(err)
		} 

		page.Click("#player_container")

		time.Sleep(5 * time.Second)

		page.Dblclick("#player_container")

		element, err := iframe.QuerySelector(".jw-text-duration")
		perr(err)
		txtDuration, err := element.TextContent()
		perr(err)

		duration  := parseDuration(txtDuration)
		time.Sleep(duration)
		page.Dblclick("#player_container")

		page.Click(".nexxt")
	}
}

func parseDuration(txt string) time.Duration {
	s := strings.Split(txt, ":")
	minutes, _ := strconv.Atoi(s[0])
	seconds, _ := strconv.Atoi(s[1])

	// Im adding a slight buffer here so I'm not waiting too
	// long between episodes. I'm assuming there are credits
	// so this shouldn't be too much of an issue.
	seconds = seconds - 15

	numeric := (minutes * 60) + seconds
	return (time.Duration(numeric) * time.Second)
}

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

	flag.Parse()
	url := flag.Arg(0)
	userDataDir := os.Getenv("FIREFOX_PROFILE")

	pw, err := playwright.Run()
	perr(err)
	defer pw.Stop()
	/*
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Args: []string{
			"--no-sandbox",
			"--disable-extensions-except=" + extensionPath,
			"--load-extension=" + extensionPath,
		},
	})
	perr(err)
	defer browser.Close()

	page, err := browser.NewPage()
	*/

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
		time.Sleep(5 * time.Second)

		videoBox, err := page.QuerySelector("#player_container")
		perr(err)
		videoBox.ScrollIntoViewIfNeeded()

		iframeElement, err := page.QuerySelector("#player_container iframe")
		perr(err)
		iframe, err := iframeElement.ContentFrame()
		perr(err)

		popup, err := iframe.QuerySelector(".jwp-popup")
		if popup == nil && err == nil {
			page.Click("#player_container")
		} else {
			time.Sleep(11 * time.Second)
			page.Click("#player_container")
		}

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

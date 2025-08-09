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
	var setup bool

	flag.StringVar(&lastEpisodeFile, "l", "last-episode.txt", "")
	flag.StringVar(&lastEpisodeFile, "last-episode", "last-episode.txt", "")
	flag.BoolVar(&noLastFile, "n", false, "")
	flag.BoolVar(&noLastFile, "no-write", false, "")
	flag.BoolVar(&setup, "setup-profile", false, "")

	flag.Usage = func() { 
		fmt.Println("usage: kcauto <url-to-episode>")
		fmt.Println("\t-l/--last-episode FILENAME\tset last episode file")
		fmt.Println("\t-n/--no-write\t\t\tdo not create last-episode file")
		fmt.Println("\t-h/--help\t\t\tdisplay this help message")
		fmt.Println("\t--setup-profile\t\t\topen browser using firefox profile (for ublock setup)")
		os.Exit(0)
	}

	flag.Parse()
	userDataDir := os.Getenv("FIREFOX_PROFILE")

	url := flag.Arg(0)
	if url == "" && !setup {
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

	if setup {
		fmt.Println("use signal interrupt to stop...")
		for true {
			time.Sleep(time.Second * 360)
		}
	}

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

		time.Sleep(10 * time.Second)

		hServerOption, err := page.QuerySelector("option[sv=\"hserver\"][selected=\"\"]")
		perr(err)
		if hServerOption == nil {
			options, err := page.QuerySelectorAll("select#selectServer option")
			perr(err)
			lastOption := options[len(options)-1]
			lastValue, err := lastOption.GetAttribute("value")
			perr(err)
			selectBox, err := page.QuerySelector("select#selectServer")
			perr(err)
			values := []string{lastValue}
			_, err = selectBox.SelectOption(playwright.SelectOptionValues{ 
				Values: &values,
			})
			perr(err)
			continue
		}

		videoBox, err := page.QuerySelector("#player_container")
		perr(err)
		videoBox.ScrollIntoViewIfNeeded()

		time.Sleep(time.Second * 1)
		upgradeBox, err := page.QuerySelector("#upgrade_pop[style=\"display:none\"]")
		perr(err)
		if upgradeBox == nil {
			upgradeBoxClose, err := page.QuerySelector(".pop_close[rel=\"#upgrade_pop\"]")
			perr(err)
			upgradeBoxClose.Click()
		}


		iframeElement, err := page.QuerySelector("#player_container iframe")
		perr(err)
		iframe, err := iframeElement.ContentFrame()
		perr(err)

		popup, err := iframe.QuerySelector(".jwp-popup")
		perr(err)
		for popup != nil {
			popup, err = iframe.QuerySelector(".jwp-popup")
			perr(err)
			time.Sleep(3 * time.Second)
		}
		page.Click("#player_container")

		time.Sleep(6 * time.Second)

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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var defaultURL = "https://github.com/drcongo/spammy-recruiters/" +
	"raw/master/spammers.txt"

var urlFlag = flag.String("url", defaultURL,
	"URL to download spammy recruiters list from",
)

var fileFlag = flag.String("file", "",
	"file to get spammy recruiters from instead of downloading them from URL",
)

var outputFlag = flag.String("o", "",
	"Output file path, prints to STDOUT if empty",
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var spammers []byte
	var err error
	if *fileFlag != "" {
		spammers, err = os.ReadFile(*fileFlag)
	} else {
		spammers, err = spammersFromURL(ctx, *urlFlag)
	}
	if err != nil {
		log.Fatal(err)
	}

	search := formatSpammers(spammers)
	rule, err := renderRule(search)
	if err != nil {
		log.Fatal(err)
	}

	if *outputFlag == "" {
		fmt.Println(string(rule))
	} else {
		err = os.WriteFile(*outputFlag, rule, 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Rule struct {
	Combinator         string `json:"combinator,omitempty"`
	Conditions         string `json:"conditions,omitempty"`
	Discard            bool   `json:"discard"`
	FileIn             string `json:"fileIn,omitempty"`
	MarkFlagged        bool   `json:"markFlagged"`
	MarkRead           bool   `json:"markRead"`
	MarkSpam           bool   `json:"markSpam"`
	Name               string `json:"name,omitempty"`
	PreviousFileInName string `json:"previousFileInName,omitempty"`
	RedirectTo         string `json:"redirectTo,omitempty"`
	Search             string `json:"search,omitempty"`
	ShowNotification   bool   `json:"showNotification"`
	SkipInbox          bool   `json:"skipInbox"`
	SnoozeUntil        string `json:"snoozeUntil,omitempty"`
	Stop               bool   `json:"stop"`
}

func spammersFromURL(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func formatSpammers(spammers []byte) string {
	spammers = bytes.ReplaceAll(spammers, []byte("\n"), []byte{})
	spammers = bytes.ReplaceAll(spammers, []byte("\r"), []byte{})

	parts := bytes.Split(spammers, []byte(" OR "))
	conds := []string{}

	for _, part := range parts {
		cond := strings.TrimSpace(string(part))
		if strings.Contains(cond, " ") {
			cond = "(" + cond + ")"
		}
		conds = append(conds, "from:"+cond)
	}

	return strings.Join(conds, " OR ")
}

func renderRule(search string) ([]byte, error) {
	rule := Rule{
		Combinator:       "any",
		Discard:          false,
		FileIn:           "Recruiter Spam",
		MarkFlagged:      false,
		MarkRead:         false,
		MarkSpam:         false,
		Name:             "Spammy Recruiters",
		Search:           search,
		ShowNotification: false,
		SkipInbox:        true,
		Stop:             false,
	}

	return json.MarshalIndent([]Rule{rule}, "", "  ")
}

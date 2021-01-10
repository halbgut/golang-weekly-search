package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	query      = flag.String("search", "", "The string to look for")
	ignorecase = flag.Bool("ignorecase", false, "Do a case insensitive search")
	verbose    = flag.Bool("verbose", false, "Verbose output")
	single     = flag.Bool("single", false, "Stop after one match has been found")
)

func main() {
	usage, err := run()
	if err != nil {
		fmt.Println(err)
		if usage {
			flag.Usage()
		}
		os.Exit(1)
	}
}

func run() (bool, error) {
	flag.Parse()
	base := "https://golangweekly.com/issues/"
	if len(*query) == 0 {
		return true, fmt.Errorf("Missing -search flag")
	}
	for i := 344; i >= 0; i-- {
		url := fmt.Sprintf("%v%v", base, i)
		start := time.Now()
		ma, err := process(url, *query, *ignorecase)
		if err != nil {
			return false, fmt.Errorf("Failed to fetch issue: %w", err)
		}
		if *single && ma {
			break
		}
		if *verbose {
			fmt.Println("Took this long:", time.Now().Sub(start))
		}
	}
	return false, nil
}

func scheduler() {
}

func process(url string, query string, ic bool) (bool, error) {
	res, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("Request failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Unexpected status code: %v", res.StatusCode)
	}
	defer res.Body.Close()
	node, err := html.Parse(res.Body)
	if err != nil {
		return false, fmt.Errorf("Failed to decode body: %w", err)
	}
	ma := search(node, query, ic)
	if ma != nil {
		fmt.Printf("Found match in %v: %v\n", url, ma.Data)
		return true, nil
	} else {
		return false, nil
	}
}

func search(n *html.Node, str string, ic bool) *html.Node {
	if ic && strings.Contains(n.Data, str) {
		return n
	} else if strings.Contains(strings.ToLower(n.Data), str) {
		return n
	}
	if n.FirstChild != nil {
		ma := search(n.FirstChild, str, ic)
		if ma != nil {
			return ma
		}
	}
	if n.NextSibling != nil {
		ma := search(n.NextSibling, str, ic)
		if ma != nil {
			return ma
		}
	}
	return nil
}

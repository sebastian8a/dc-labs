// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		//log.Print(err)
	}
	return list
}

//!-sema

//!+
type Worklist struct {
	urls  []string
	depth int
}

func main() {
	worklist := make(chan Worklist)
	var n int // number of pending sends to worklist
	var depth = flag.Int("depth", 0, "max depth level")
	var filename = flag.String("results", "results.txt", "result file")
	flag.Parse()
	if *depth == 0 {
		fmt.Print("Missing depth argument\n")
		os.Exit(0)
	}
	newFile, err := os.Create(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	// Start with the command-line arguments.
	n++
	go func() { worklist <- Worklist{urls: flag.Args(), depth: 0} }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		currentDepth := list.depth
		for _, link := range list.urls {
			if currentDepth <= *depth {
				if !seen[link] {
					seen[link] = true
					_, err2 := newFile.WriteString(link + "\n")
					if err2 != nil {
						panic(err2)
					}
					n++
					go func(link string) {
						worklist <- Worklist{urls: crawl(link), depth: currentDepth + 1}
					}(link)
				}
			}
		}
	}
}

//!-

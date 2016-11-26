package main

import (
	"flag"
	"log"
)

var (
	url string
)

func init() {
	flag.StringVar(&url, "url", "", "url of the website")
	flag.Parse()
}

func main() {
	verifyFlags()

	log.Printf("TODO: analyze %q", url)
}

func verifyFlags() {
	if url == "" {
		log.Fatalln("Must specify url")
	}
}

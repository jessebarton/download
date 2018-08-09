package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	arg := os.Args[1]
	name := getURLName(arg)
	err := downloadFile(name, arg)
	if err != nil {
		panic(err)
	}
}

func getURLName(arg string) string {
	url, err := url.Parse(arg)
	if err != nil {
		panic(err)
	}
	path := url.Path

	uriSegments := strings.Split(path, "/")
	var segments int
	for i := range uriSegments {
		segments = i
	}
	return uriSegments[segments]
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Printf("Downloading File: %s", filepath)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

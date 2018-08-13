package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	arg := os.Args[1]
	name := getURLName(arg)

	err := downloadFile(arg, "./")
	if err != nil {
		panic(err)
	}

	f := fileInfo(name)
	fmt.Printf("\nName	Size\n%s	%vk\n ", f.Name(), f.Size())
}

func fileInfo(path string) os.FileInfo {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	return stat
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

func printDownloadPercent(done chan int64, path string, total int64) {
	var stop = false

	for {
		select {
		case <-done:
			stop = true
		default:
			f := fileInfo(path)
			size := f.Size()
			if size == 0 {
				size = 1
			}

			var percent = float64(size) / float64(total) * 100
			fmt.Printf("%.0f%s", percent, "%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}

func downloadFile(url string, dest string) error {
	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	start := time.Now()

	out, err := os.Create(path.String())
	if err != nil {
		fmt.Println(path.String())
		panic(err)
	}
	defer out.Close()

	headResp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
	if err != nil {
		panic(err)
	}

	done := make(chan int64)
	go printDownloadPercent(done, path.String(), int64(size))

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	if err != nil {
		panic(err)
	}

	done <- n

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)

	return nil
}

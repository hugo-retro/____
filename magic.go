package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"menteslibres.net/gosexy/to"
)

type Cache struct {
	cache   int64
	latency int64
}

type EndPoint struct {
	latency    int64
	num_caches int64
	caches     []Cache
}

type Request struct {
  num_requests int64
  video int64
  endpoint int64
}

var num_videos int64
var num_endpoints int64
var num_requests int64
var num_caches int64
var cache_size int64
var videos []int
var endpoints []EndPoint
var requests []Request

func parseCache(lines []string, offset int64, ep EndPoint) {
	//print("parseCache ", offset, " ", lines[offset],"\n")
	parts := strings.Split(lines[offset], " ")
	c := Cache{to.Int64(parts[0]), to.Int64(parts[1])}
	ep.caches = append(ep.caches, c)
}

func parseEndpoint(lines []string, offset int64) int64 {
	//print("parseEndpoint ", offset, " ", lines[offset], " ", "\n")
	parts := strings.Split(lines[offset], " ")
	ep := EndPoint{to.Int64(parts[0]), to.Int64(parts[1]), make([]Cache, 0)}
	endpoints = append(endpoints, ep)

	var x int64

	if ep.num_caches > 0 {
		for x = 0; x < ep.num_caches; x++ {
			parseCache(lines, offset+x+1, ep)
		}
	}

	return ep.num_caches
}

func readFile(filename string) {
	var lines []string

	content, err := ioutil.ReadFile(filename)

	if err != nil {
		print("Failed to open ", filename, err)
		os.Exit(1)
	}

	lines = strings.Split(string(content), "\n")

	nrs := strings.Split(lines[0], " ")
	num_videos = to.Int64(nrs[0])
	num_endpoints = to.Int64(nrs[1])
	num_requests = to.Int64(nrs[2])
	num_caches = to.Int64(nrs[3])
	cache_size = to.Int64(nrs[4])

	var offset int64 = 2
	var x int64

	readVideos(lines)

	// Parse the endpoints and their
	for x = 0; x < num_endpoints; x++ {
		offset += parseEndpoint(lines, offset) + 1
	}
}

func writeOutput(filename string, lines []string) {
	tempFilename := filename + ".out"

	file, err := os.OpenFile(tempFilename, os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		println("Failed to open ", tempFilename, err.Error())
		os.Exit(1)
	}

	for _, line := range lines {
		file.WriteString(line)
		file.WriteString("\n")
	}
}

func magic(lines []string) []string {
	return lines
}

func main() {
	if len(os.Args) == 1 {
		println("Usage: go run magic.go filename")
		os.Exit(1)
	}

	filename := os.Args[1]

	readFile(filename)
}

func readVideos(lines []string) {

	var strs = strings.Split(lines[1], " ")

	fmt.Println(strs[len(strs)-1])

	for i := 0; i < len(strs); i++ {
		var video, err = strconv.Atoi(strings.TrimSpace(strs[i]))

		fmt.Println(video, err)
		if err == nil {
			videos = append(videos, video)
		}
	}
}

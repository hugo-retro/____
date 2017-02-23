package main

import (
  "io/ioutil"
  "strings"
  "os"
  "strconv"
)

var num_videos int;
var num_endpoints int;
var num_requests int;
var num_caches int;
var cache_size int;

func readFile(filename string) ([]string) {
  var lines []string;

  content, err := ioutil.ReadFile(filename)

  if err != nil {
    print("Failed to open ", filename, err)
    os.Exit(1)
  }

  lines = strings.Split(string(content), "\n")

  nrs := strings.Split(lines[0], " ")
  num_videos, err = strconv.Atoi(nrs[0])
  num_endpoints, err = strconv.Atoi(nrs[1])
  num_requests, err = strconv.Atoi(nrs[2])
  num_caches, err = strconv.Atoi(nrs[3])
  cache_size, err = strconv.Atoi(nrs[4])

  return lines[1:len(lines) - 1]
}

func writeOutput(filename string, lines []string) {
  tempFilename := filename + ".out"

  file, err := os.OpenFile(tempFilename, os.O_RDWR|os.O_CREATE, 0600)

  if err != nil {
    println("Failed to open ", tempFilename, err.Error())
    os.Exit(1)
  }

  for _,line := range lines {
    file.WriteString(line)
    file.WriteString("\n")
  }
}

func magic(lines []string) ([]string) {
  return lines
}

func main() {
  if len(os.Args) == 1 {
    println("Usage: go run magic.go filename")
    os.Exit(1)
  }

  filename := os.Args[1]

  writeOutput(filename, magic(readFile(filename)))
}

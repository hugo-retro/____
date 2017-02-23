package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "os"
)

func readFile(filename string) ([]string, error) {
  var lines []string;

  content, err := ioutil.ReadFile(filename)

  if err == nil {
    lines = strings.Split(string(content), "\n")
  }

  return lines[:len(lines) - 1], err
}

func writeOutput(lines []string) {
  for _,line := range lines {
    fmt.Println(line)
  }
}

func magic(lines []string) ([]string) {
  return lines
}

func main() {
  // TODO: Add check here for argv not having at least the filename
  lines, err := readFile(os.Args[1])

  if err != nil {
    fmt.Println("reading file failed! Reason: ", err)
    os.Exit(1)
  }

  writeOutput(magic(lines))
}

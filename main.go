package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	list = flag.Bool("list", false, "list mode")
)

func main() {
	flag.Parse()

	if *list {
		file := getStatePath()
		path, err := filepath.Abs(file)
		if err != nil {
			errAndExit(fmt.Errorf("Error determining directory: %s", err))
		}

		f, err := os.Stat(path)
		if err != nil {
			errAndExit(fmt.Errorf("Error determining directory: %s", err))
		}

		if !f.IsDir() {
			errAndExit(fmt.Errorf("Invalid directory: %s", file))
		}

		s, err := getState(path)
		if err != nil {
			errAndExit(err)
		}

		j, err := s.ToJSON()
		if err != nil {
			errAndExit(err)
		}

		fmt.Println(j)
	}
}

func getStatePath() string {
	var v string

	v = os.Getenv("TF_STATE")
	if v != "" {
		return v
	}

	return "."
}

func errAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

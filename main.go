package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	list    = flag.Bool("list", false, "list mode")
	command = Terraform
)

const (
	Terraform  = "terraform"
	Terragrunt = "terragrunt"
)

func main() {
	flag.Parse()

	v := os.Getenv("TF_TERRAGRUNT")
	if v != "" {
		command = Terragrunt
	}

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

		if s == nil {
			fmt.Println("No state was found")
			os.Exit(1)
		}

		j, err := ToJSON(s)
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

func getState(path string) (State, error) {
	var out bytes.Buffer
	var state State
	terraformVersion := "0.12"

	cmd := exec.Command(command, "state", "pull")
	cmd.Dir = path
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running `%s state pull` in directory %s, %s\n", command, path, err)
	}

	b, err := ioutil.ReadAll(&out)
	if err != nil {
		return nil, fmt.Errorf("Error reading output of `%s state pull`: %s\n", command, err)
	}

	// If there was no output, return nil and no error
	if string(b) == "" {
		return nil, nil
	}

	if string(b[0]) == "o" && string(b[1]) == ":" {
		b = append(b[:0], b[2:]...)
	}

	var tmpState interface{}
	err = json.Unmarshal(b, &tmpState)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling state: %s\n", err)
	}

	if v, ok := tmpState.(map[string]interface{}); ok {
		if v, ok := v["version"].(float64); ok {
			if v == 3 {
				terraformVersion = "0.11"
			}
		}
	}

	switch terraformVersion {
	case "0.11":
		var s StateV011
		err = json.Unmarshal(b, &s)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling state: %s\n", err)
		}
		state = s
	default:
		var s StateV012
		err = json.Unmarshal(b, &s)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling state: %s\n", err)
		}
		state = s
	}

	return state, nil
}

func errAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

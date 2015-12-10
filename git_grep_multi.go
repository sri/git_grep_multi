package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var usage = `Usage: gg search-term <list-of-repos>

Runs 'git grep --break -n <search-term>' on the current git repo or on
multiple repos in the current directory. Optionally, the list of git
repos can be passed in.`

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	search_term := os.Args[1]
	dirs := os.Args[2:]

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if len(dirs) == 0 {
		if err := exec.Command("git", "rev-parse").Run(); err == nil {
			// Currently in a git repo
			dirs = []string{cwd}
		} else {
			// Filter out git repos under the current dir
			list, err := ioutil.ReadDir(cwd)
			if err != nil {
				fmt.Println("error:", err)
				os.Exit(1)
			}
			for _, dir := range list {
				if !dir.IsDir() {
					continue
				}
				full := path.Join(cwd, dir.Name())
				dotgit := path.Join(full, ".git")
				if f, err := os.Stat(dotgit); err == nil && f.IsDir() {
					dirs = append(dirs, full)
				}
			}
		}
	}

	for _, d := range dirs {
		if err := os.Chdir(d); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		cmd := exec.Command("git", "grep", "--break", "-n", search_term)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		if err = os.Chdir(cwd); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

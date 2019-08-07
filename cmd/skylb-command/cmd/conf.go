package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

func getConfDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	dir := filepath.Join(usr.HomeDir, ".skylb-command")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	return dir
}

func getConfFile() string {
	return filepath.Join(getConfDir(), "config")
}

func getHistoryFile() string {
	return filepath.Join(getConfDir(), "history")
}

func createLiner() *liner.State {
	line := liner.NewLiner()
	line.SetCompleter(func(line string) (c []string) {
		for n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(getHistoryFile()); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	return line
}

func saveLiner(liner *liner.State) {
	f, err := os.Create(getHistoryFile())
	if err != nil {
		fmt.Println("Error writing history file: ", err)
		os.Exit(2)
	}
	defer f.Close()

	liner.WriteHistory(f)
}

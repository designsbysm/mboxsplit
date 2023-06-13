package main

import (
	"os"
	"path/filepath"

	"github.com/designsbysm/timber/v2"
)

func main() {
	// setup logger
	timber.New(
		os.Stdout,
		timber.LevelAll,
		"",
		timber.FlagColorful,
	)

	// cli args
	args := os.Args[1:]
	target := ""
	if len(args)-1 >= 0 {
		target = args[len(args)-1]
	}

	// process folders/files
	if err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() && !info.Mode().IsRegular() {
			return nil
		}

		return processFile(path)
	}); err != nil {
		timber.Error(err)
		return
	}
}

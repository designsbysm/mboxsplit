package main

import (
	"os"

	"github.com/designsbysm/timber/v2"
	"github.com/spf13/viper"
)

func config() error {
	// logger
	timber.New(
		os.Stdout,
		timber.LevelAll,
		"",
		timber.FlagColorful,
	)

	// cli
	args := os.Args[1:]
	source := "."
	if len(args)-1 >= 0 {
		source = args[len(args)-1]
	}

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return err
	}

	viper.Set("source", source)

	return nil
}

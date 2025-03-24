package main

import (
	"galvanico/cmd"
	"log/slog"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func main() {
	cmd.Execute()
}

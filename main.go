package main

import (
	"github.com/sirupsen/logrus"

	"github.com/imsilence/kaleidoscope/cmd"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	cmd.Execute()
}

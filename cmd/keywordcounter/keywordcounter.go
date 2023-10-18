package main

import (
	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/internal/keywordcounter"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/logger"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	config := keywordcounter.NewFlagConfig()
	logger.Setup(config.Verbose)
	config.MustValidate()

	keywordcounter.Run(config)
}

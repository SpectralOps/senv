package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/spectralops/senv/pkg"
	"github.com/spf13/viper"
)

var CLI struct {
	Info   bool     `help:"Show information"`
	Create bool     `help:"Create global config"`
	Config bool     `help:"Path to config file"`
	Cmd    []string `arg optional name:"cmd" help:"Command to execute"`
}
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	kong.Parse(&CLI)
	senv := pkg.Senv{
		Cmd: CLI.Cmd,
	}
	senv.Init()

	if CLI.Config {
		configPath := senv.ConfigPath()
		if len(configPath) > 0 {
			fmt.Print(configPath)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if CLI.Info {
		blue := color.New(color.FgBlue).SprintFunc()
		configPath := senv.ConfigPath()
		if len(configPath) > 0 {
			fmt.Printf("%v\n%v\n\n", blue("[config]"), viper.ConfigFileUsed())
		}
		fmt.Printf("%v\n%v\n", blue("[redact]"), strings.Join(senv.RedactList(), "\n"))
		os.Exit(1)
	}

	if CLI.Create {
		senv.CreateConfig()
		os.Exit(0)
	}

	cmd := CLI.Cmd

	if len(cmd) > 0 {
		senv.Exec()
	} else {
		senv.Print()
	}
}

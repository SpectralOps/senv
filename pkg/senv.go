package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type Senv struct {
	Cmd        []string
	ConfigName string
}

func (s *Senv) getConfigName() string {
	if s.ConfigName == "" {
		return ".senv"
	}
	return s.ConfigName
}

func (s *Senv) Init() {
	viper.SetConfigName(s.getConfigName())

	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}

func bail(e error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", e)
	os.Exit(1)
}

func execCmd(cmd string, cmdArgs []string) error {
	command := exec.Command(cmd, cmdArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func (s *Senv) CreateConfig() {
	d, _ := os.UserHomeDir()
	f := path.Join(d, fmt.Sprintf("%v.yaml", s.getConfigName()))
	err := ioutil.WriteFile(f,
		[]byte(`
# Use this to start from a blank list of redactions
# no_defaults: true
#
# Add to this list if you like
# redact:
# - FOO_BAR
`),
		0644,
	)
	if err != nil {
		bail(err)
	}
	fmt.Printf("Wrote config to %v", f)
}

func (s *Senv) ConfigPath() string {
	return viper.ConfigFileUsed()
}

func (s *Senv) RedactList() []string {
	noDefaults := viper.GetBool("no_defaults")
	if noDefaults {
		return viper.GetStringSlice("redact")
	}
	return append(REDACT, viper.GetStringSlice("redact")...)
}

func (s *Senv) Clean() {
	redact := s.RedactList()
	for _, b := range redact {
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			m, _ := regexp.MatchString(b, pair[0])
			if m {
				os.Unsetenv(pair[0])
			}
		}
	}
}

func (s *Senv) Print() {
	s.Clean()
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}
}

func (s *Senv) Exec() {
	s.Clean()
	err := execCmd(s.Cmd[0], s.Cmd[1:])
	if err != nil {
		bail(err)
	}
}

package pkg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
)

func cleanup() string {
	d, _ := os.UserHomeDir()
	f := path.Join(d, "senv.test.yaml")
	os.Remove(f)
	return f
}

func TestRedactWithSettings(t *testing.T) {
	senv := Senv{}
	redact := senv.RedactList()
	assert.Equal(t, REDACT, redact)
}

func TestDefaultRedact(t *testing.T) {
	senv := Senv{}
	redact := senv.RedactList()
	assert.Equal(t, REDACT, redact)
}

func TestConfigInitNotExisting(t *testing.T) {
	cleanup()
	senv := Senv{ConfigName: "senv.test"}
	senv.Init()

	redact := senv.RedactList()
	assert.Equal(t, REDACT, redact)
}

func TestConfigInitWithConfigFile(t *testing.T) {
	cleanup()
	senv := Senv{ConfigName: "senv.test"}
	senv.Init()

	senv.CreateConfig()

	redact := senv.RedactList()
	assert.Equal(t, REDACT, redact)
}

func TestConfigInitWithCustomConfigFile(t *testing.T) {
	f := cleanup()
	ioutil.WriteFile(f,
		[]byte(`
redact:
- FOO_BAR
`),
		0644,
	)
	senv := Senv{ConfigName: "senv.test"}
	senv.Init()

	redact := senv.RedactList()
	assert.Equal(t, append(REDACT, []string{"FOO_BAR"}...), redact)
}

func TestConfigInitWithCustomConfigFileNoDefaults(t *testing.T) {
	f := cleanup()
	ioutil.WriteFile(f,
		[]byte(`
no_defaults: true
redact:
- FOO_BAR
`),
		0644,
	)
	senv := Senv{ConfigName: "senv.test"}
	senv.Init()

	redact := senv.RedactList()
	assert.Equal(t, []string{"FOO_BAR"}, redact)
}

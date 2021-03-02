package env

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/subosito/gotenv"
)

// ReadEnvFile read .env* file
func ReadEnvFile(path string) (map[string]interface{}, error) {
	mappedEnv := make(map[string]interface{})

	var content io.Reader
	var err error

	content, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	env, err := gotenv.StrictParse(buf)
	if err != nil {
		return nil, err
	}
	for k, v := range env {
		mappedEnv[k] = v
	}
	return mappedEnv, nil
}

// WriteEnvFile write .env* file
func WriteEnvFile(path string, data *map[string]interface{}) error {
	lines := []string{}
	for key, val := range *data {
		envName := strings.ToUpper(strings.Replace(key, ".", "_", -1))
		lines = append(lines, fmt.Sprintf("%v=\"%v\"", envName, val))
	}
	out := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(path, []byte(out), 0644); err != nil {
		return err
	}
	return nil

}

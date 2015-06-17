package main

import (
	"strings"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Script struct {
	Path string
	Env []string
}

func NewScript(path string) *Script {
	return &Script{ Path: path, Env: nil }
}

func (s *Script) Exists() bool {
	return s.Path != ""
}

func (s *Script) Run() (string, bool) {

	command := []string{s.Path}
	shebang, err := s.Shebang()

	if err != nil {
		return err.Error(), false
	}

	if shebang != "" {
		command = append(strings.Split(shebang, " "), command...)
	}

	cmd := exec.Command(command[0], command[1:]...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err.Error(), false
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err.Error(), false
	}

	if s.Env != nil {
		cmd.Env = s.Env
	}

	err = cmd.Start()
	if err != nil {
		return err.Error(), false
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}

func (s *Script) Shebang() (string, error) {

	f, err := os.Open(s.Path)
	if (err != nil) {
		return "", err
	}
	defer f.Close()

	shebang := ""
	r := bufio.NewReader(f)
	for i:=0; i<4; i++ {
		if lineBytes, _, err := r.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			break
		} else {
			line := string(lineBytes)
			line = strings.Trim(line," \t")
			if strings.HasPrefix(line, "#!") {
				shebang = strings.Trim(line[2:], " \t")
				break;
			}
		}
	}

	return shebang, nil
}

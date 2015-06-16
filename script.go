package main

import (
	"strings"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Script string

func (s Script) Exists() bool {
	return s != ""
}

func (s Script) Run() (string, bool) {

	command := []string{string(s)}
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

func (s Script) Shebang() (string, error) {

	f, err := os.Open(string(s))
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

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

func (s Script) Run() {

	command := []string{string(s)}
	shebang := s.Shebang()

	if shebang != "" {
		command = append(strings.Split(shebang, " "), command...)
	}

	cmd := exec.Command(command[0], command[1:]...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		mainLog("Post Build Script Failed: \n %s", string(errBuf))
	}
}

func (s Script) Shebang() string {

	f, err := os.Open(string(s))
	if (err != nil) {
		return ""
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

	return shebang
}

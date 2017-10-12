package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	err := os.MkdirAll("coverages", 0755)
	if err != nil {
		log.Fatalln("could not create coverages directory: ", err.Error())
	}
	files, err := ioutil.ReadDir("coverages")
	if err != nil {
		log.Fatalln("could not read coverage directory: ", err.Error())
		return
	}

	var n int

	if len(files) != 0 {
		filename := files[len(files)-1].Name()
		_, err := fmt.Sscanf(filename, "%04d.out", &n)
		if err != nil {
			log.Fatalln("coverage filename was in wrong format: ", err.Error())
		}
		n++
	}
	args := []string{
		fmt.Sprintf("-test.coverprofile=%04d.out", n),
		"-test.outputdir=coverages",
	}
	cmd := exec.Command("testvg", append(args, os.Args[1:]...)...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	stderrStr := stderr.String()
	if err != nil && len(stderrStr) == 0 {
		fmt.Fprint(os.Stdout, stdout.String())
		fmt.Fprint(os.Stderr, stderrStr)
		fmt.Fprintf(os.Stderr, "couldn't run testvg: %v\n", err.Error())
		os.Exit(1)
	}
	stdoutLines := strings.Split(stdout.String(), "\n")
	if len(stdoutLines) >= 3 {
		stdoutLines = stdoutLines[:len(stdoutLines)-3]
	}
	fmt.Fprint(os.Stdout, strings.Join(stdoutLines, "\n")+"\n")
	fmt.Fprint(os.Stderr, stderrStr)

	if err != nil {
		os.Exit(1)
	}
}

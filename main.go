package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/takama/daemon"
)

type Service struct {
	daemon.Daemon
}

var r = regexp.MustCompile("^(1{2,}|2{2,}|3{2,}|4{2,}|5{2,}|6{2,}|7{2,}|8{2,}|9{2,}|0{2,})$")

func (service *Service) Manage() (string, error) {
	usage := "Usage: go-kiriban install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Find kiriban process
	var detectedPid []string
	for {
		processes, _ := ps.Processes()
		for _, p := range processes {
			pid := strconv.Itoa(p.Pid())
			if !contains(detectedPid, pid) && r.MatchString(pid) {
				// Show notification
				command := "tell application \"System Events\" to display dialog \"" + pid + "のプロセスキリ番踏んだよ！(๑˃̵ᴗ˂̵)و\" with title \"go-kiriban\" giving up after 5"
				c := exec.Command("/usr/bin/osascript", "-e", command)
				err := c.Run()
				if err != nil {
					// Noop
				}

				detectedPid = append(detectedPid, pid)
			}
		}

		time.Sleep(5 * time.Second)
	}

	return usage, nil
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func main() {
	srv, err := daemon.New("go-kiriban", "Detecting kiriban process")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		log.Fatal(status, "\nError", err)
	}
	fmt.Println(status)
}

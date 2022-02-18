package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// logfile is a variable that holds the log file path.
// structure of the file: label start end duration
var logfile string = "./.statch.csv"

// statch is a structure that holds each logged item's label, start, end and duration.
type statch struct {
	label    string
	start    time.Time
	end      time.Time
	duration time.Duration
}

// help is a function that prints help.
func help() {
	fmt.Println("TODO:")
}

// start is a function that starts a new insert for given label.
// IDEA: // If any previous insert was still open for given label, then that insert get's a closed value.
func start(label string) {
	srt := time.Now()
	line := fmt.Sprintf("%v,%v,%v,%v\n", label, srt.Unix(), "", "")
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added start time '%v' to label '%v'.\n", srt.String(), label)
}

// logged is a function that CURRENTLY prints the contents of the logfile.
func logged(label string, printing bool) ([]statch, int) {
	var statches []statch
	var fstOpen int = -1
	if _, err := os.Stat(logfile); err != nil {
		fmt.Println("Nothing logged yet")
		return statches, fstOpen
	}
	f, err := os.Open(logfile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
	var i int
	for scanner.Scan() {
		contents := strings.Split(scanner.Text(), ",")
		if contents[0] != label {
			continue
		}
		srt, err := strconv.ParseInt(contents[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		srtTime := time.Unix(srt, 0)
		var endTime time.Time
		var duration time.Duration
		if contents[2] != "" {
			if fstOpen == -1 {
				fstOpen = i
			}
			end, err := strconv.ParseInt(contents[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			endTime = time.Unix(end, 0)
			duration = endTime.Sub(srtTime)
		}
		st := statch{contents[0], srtTime, endTime, duration}
		statches = append(statches, st)
		if printing {
			fmt.Println(st.label, st.start, st.end, st.duration)
		}
		i++
	}
	fmt.Println(statches)
	return statches, fstOpen
}

// TODO:
func end(label string) {
	return
}

func main() {
	labelFlag := flag.String("label", "general", "Label of saved time")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 	var statches []statch
	switch action := os.Args[flag.NFlag()+1]; action {
	case "help":
		help()
	case "start":
		start(*labelFlag)
	case "logged":
		_, _ = logged(*labelFlag, true)
	//case "end":
	//	statches, openInd = logged(*labelFlag, false)
	//	end(*labelFlag, openInd)
	default:
		log.Fatal("Not defined")
	}
}

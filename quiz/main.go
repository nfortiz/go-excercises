package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFileName := flag.String("csv", "assets/problems.csv","a csv file in the format of question,answer")
	flag.Parse()
	file, err := os.Open(*csvFileName)

	timeLimit := flag.Int("limit", 30, "The time limit for the quiz in seconds")

	if err != nil {
		exit(fmt.Sprintf("Failed to load csv file: %s \n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s\n", i+1, p.q)
		answerChn := make(chan string)
		go func () {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChn <- answer
		} ()

		select {
		case <- timer.C:
			fmt.Printf("You scored %d out %d \n", correct, len(problems))
			return
		case answer := <- answerChn:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out %d \n", correct, len(problems))
}

func parseLines(lines[][] string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

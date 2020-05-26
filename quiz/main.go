package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"os"
)

var messages = make(chan string)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a CSV file in the format of 'question,answer'")
	timerSeconds := flag.Int("seconds", 30, "Set in seconds the quiz timer duration")
	randomizeProblems := flag.Bool("randomize", false, "Randomize the order of the quiz questions")
	flag.Parse()
	
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}

	problems := parseLines(lines)

	fmt.Printf("%T\n", problems)

	if *randomizeProblems {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Println(problems)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Timed Pop Quiz!")
	fmt.Printf("%d seconds set for timer\n", *timerSeconds)
  	fmt.Print("Press 'Enter' to start quiz...")
  	reader.ReadBytes('\n')

  	go quizTimer(*timerSeconds)
  	go quiz(problems)

  	for {
  		select {
  		case msg := <- messages:
  			exit(msg)

  		}
  	}

//	quizTimeout := time.Second * time.Duration(timerSeconds)
//	quizTimer := time.NewTimer(quizTimeout)
	//quizTimerC = quizTimer.C

}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}


//[]main.problem

func quiz(quizProblems []problem) {
	correct := 0
	wrong := 0

	width := len(strconv.Itoa(len(quizProblems)))

	for i, p := range quizProblems {
		fmt.Printf("Problem #%-*d: %s = \n", width, i+1, p.q)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == p.a {
			//fmt.Println("Correct")
			correct++
		} else {
			//fmt.Println("WRONG")
			wrong++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(quizProblems))
	os.Exit(0)
}

func quizTimer(quizDuration int) {
	time.Sleep(time.Second * time.Duration(quizDuration))
	messages <- "\n** Quiz timer elapsed! **"
	//fmt.Println("elapsed!")
}
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"os"
)

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

func describe(i interface{}) {
	fmt.Printf("(%+v, %T)\n", i, i)
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func takeQuiz(quizProblems []problem, channel chan problem) {

	// just to have nice padded formating of problems
	width := len(strconv.Itoa(len(quizProblems)))

	var correct string
	var sendAnswer problem
	var answer string

	for i, p := range quizProblems {
		fmt.Printf("Problem #%-*d: %s = \n", width, i+1, p.q)
		fmt.Scanf("%s\n", &answer)
		if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(p.a)) {
			correct = "CORRECT"
		} else {
			correct = "WRONG"
		}
		sendAnswer = problem{
			q: correct,
			a: answer,
		}
		channel <- sendAnswer
	}
	close(channel)
}

func startTimer(quizDuration int, done chan bool) {
	time.Sleep(time.Second * time.Duration(quizDuration))
	done <- true
	close(done)
}

func gradeQuiz(problems []problem, answers []problem) {
	var correctCount int
	// just to have nice padded formating of problems
	fmt.Println("\n******** Quiz Restults ********")
	width := len(strconv.Itoa(len(problems)))
	for i, p := range problems {
		fmt.Printf("Problem #%-*d: %s =\n", width, i+1, p.q)
		fmt.Printf("Correct Answer: %s\n", p.a)
		fmt.Printf("  Given Answer: %s\n", answers[i].a)
		fmt.Printf("Your Answer is: %s\n\n", answers[i].q)
		if answers[i].q == "CORRECT" {
			correctCount++
		}
	}

	fmt.Printf("\nYou scored %d out of %d.\n\n", correctCount, len(problems))
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a CSV file in the format of 'question,answer'")
	timerSeconds := flag.Int("seconds", 30, "Set in seconds the quiz timer duration")
	randomizeProblems := flag.Bool("randomize", false, "Randomize the order of the quiz questions")
	debugOutput := flag.Bool("debug", false, "Output some variables for debugging")
	flag.Parse()
	
	
	file, err := os.Open(*csvFilename)
	codeFileOpenFail := 10
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename), codeFileOpenFail)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	codeFileReadFail := 11
	if err != nil {
		exit("Failed to parse the provided CSV file", codeFileReadFail)
	}
	defer file.Close()

	quizProblems := parseLines(lines)
	if *randomizeProblems {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quizProblems), func(i, j int) {
			quizProblems[i], quizProblems[j] = quizProblems[j], quizProblems[i]
		})
	}

  	var quizAnswers = make([]problem, len(quizProblems))
  	var answerNum int
  	var codeQuizDone int

  	// Pre-Populate quizAnswers with WRONG answers for faster calculation
  	for i := range quizProblems {
  		quizAnswers[i] = problem{
  			q: "WRONG",
  			a: "**NO ANSWER GIVEN**",
  		}
  	}

	if *debugOutput {
		describe(*csvFilename)
		describe(*timerSeconds)
		describe(*randomizeProblems)
		describe(lines)
		describe(quizProblems)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n*** Timed Pop Quiz! ***")
	fmt.Printf("%d seconds set for timer\n", *timerSeconds)
  	fmt.Print("Press 'Enter' to start quiz...\n\n")
  	reader.ReadBytes('\n')
	
  	timerDone := make(chan bool)
	answersChannel := make(chan problem)

  	go startTimer(*timerSeconds, timerDone)
  	go takeQuiz(quizProblems, answersChannel)

  	isRunning := true
  	for isRunning {
  		select {
  		case RecievedAnswer, ok := <- answersChannel:
  			if ok {
	  			quizAnswers[answerNum] = RecievedAnswer
  				answerNum++
  				if *debugOutput {describe(answersChannel)}
  			} else {
  				isRunning = false
  			}
  		case <- timerDone:
  			if *debugOutput {describe(timerDone)}
  			fmt.Println("\n *** Timmer has elapsed! Pencils DOWN!! ***")
  			codeQuizDone = 1
  			isRunning = false
  		}
  	}
	if *debugOutput {describe(quizAnswers)}
  	gradeQuiz(quizProblems, quizAnswers)
 	exit("Quiz is done!", codeQuizDone)
}

var csvFilename = strings.NewReader(`5+5,10
1+1,2
8+3,11
1+2,3
8+6,14
3+1,4
1+4,5
5+1,6
2+3,5
3+3,6
2+4,6
5+2,7
`)


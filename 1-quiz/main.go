package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func readProblemsFromCsv(csvName string) []problem {
	file, err := os.Open(csvName)
	if err != nil {
		log.Fatal("Error reading from file ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	problems, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error parsing csv data ", err)
	}

	problemsList := make([]problem, 0)
	for _, record := range problems {
		problem := problem{
			question: record[0],
			answer:   record[1],
		}
		problemsList = append(problemsList, problem)
	}

	return problemsList
}

func startTimer(reader *bufio.Reader, onTimeout func()) {
	fmt.Println("Press return when you're ready!")
	reader.ReadString('\n')

	timerSeconds := 30
	if len(os.Args) > 1 {
		time, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Timer argument must be an integer ", err)
		}
		timerSeconds = time
	}

	timer := time.NewTimer(time.Duration(timerSeconds) * time.Second)
	fmt.Println(timerSeconds, "second timer started!")
	go func() {
		<-timer.C
		onTimeout()
	}()
}

func askQuestions(problems []problem) {
	reader := bufio.NewReader(os.Stdin)
	score := 0

	onTimeout := func() {
		fmt.Println("Time's up!")
		fmt.Printf("You got %d out of %d!\n", score, len(problems))
		os.Exit(0)
	}
	startTimer(reader, onTimeout)

	for _, problem := range problems {
		fmt.Print("What's ", problem.question, "?:")
		userAnswer, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Failed to read from the terminal ", err)
		}

		if strings.TrimSuffix(userAnswer, "\n") == problem.answer {
			score += 1
			fmt.Println("Correct! Current Score: ", score)
		} else {
			fmt.Println("Wrong! Current Score: ", score)
		}
	}
}

func main() {

	csvName := "problems.csv"
	problems := readProblemsFromCsv(csvName)
	askQuestions(problems)
}

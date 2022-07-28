package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
) 

type Questions struct {
	question string;
	answer   string;
}

func createQuestionsList(data[][]string) []Questions {
	var questions []Questions;
	for _, row := range data {
		questions = append(questions, Questions{question: row[0], answer: row[1]});
	}
	return questions;
}

func askQuestion(question Questions) bool{
	fmt.Printf("Solve %v:\n", question.question);
	var answer string;
	fmt.Scanln(&answer);
	return answer == question.answer;
}

func main() {
	file, err := os.Open("questions.csv");
	if err != nil {
		log.Fatal(err);
	}
	defer file.Close();
	
	csvReader := csv.NewReader(file);
	data, err := csvReader.ReadAll();
	if err != nil {
		log.Fatal(err);
	}
	questionsList := createQuestionsList(data);
	totalQuestions := len(questionsList);
	correctAnswers := 0;
	timeLimit := flag.Int("time", 30, "time limit for each question");
	flag.Parse();

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second);

	for _, question := range questionsList {
		answerCh := make(chan bool);
		go func() {
			answerCh <- askQuestion(question)
		}()

		select { 
		case <-timer.C:
			fmt.Printf("You got %v out of %v questions correct.\n", correctAnswers, totalQuestions);
			return;
		case answer := <-answerCh:
			if answer {
				correctAnswers++;
			}
		}		
	}
}
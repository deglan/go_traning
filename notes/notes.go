package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/notes/note"
	"example.com/notes/todo"
)

type saver interface {
	Save() error
}

type outputable interface {
	saver
	Display()
}

func main() {

	title, content := getNoteData()
	todoText := getUserInput("Add todo: ")
	todo, err := todo.New(todoText)
	if err != nil {
		panic(err)
	}
	todo.Display()
	err = saveData(todo)
	if err != nil {
		return
	}
	note, err := note.New(title, content)
	if err != nil {
		panic(err)
	}
	err = outputData(note)

	if err != nil {
		return
	}

}

func outputData(data outputable) error {
	data.Display()
	return saveData(data)
}

func saveData(data saver) error {
	err := data.Save()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Saved!")
	return nil
}

func getNoteData() (string, string) {
	title := getUserInput("Title: ")
	content := getUserInput("Note content: ")
	return title, content
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	return text
}

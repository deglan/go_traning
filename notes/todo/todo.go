package todo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func New(text string) (Todo, error) {
	if text == "" {
		return Todo{}, errors.New("text cannot be empty")
	}
	return Todo{
		Text: text,
	}, nil
}

type Todo struct {
	Text string `json:"text"`
}

func GetTodoData() string {
	return getUserInput("Todo: ")
}

func (todo Todo) Display() {
	fmt.Printf("Todo: %s \n", todo.Text)
}

func (todo Todo) Save() error {
	fileName := "todo.json"
	toDoJson := ConvertToJson(todo)
	return os.WriteFile(fileName, toDoJson, 0644)
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

func ConvertToJson(todo Todo) []byte {
	json, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}
	return json
}

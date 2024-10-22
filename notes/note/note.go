package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Note struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func New(title, content string) (Note, error) {

	if title == "" || content == "" {
		return Note{}, errors.New("invalid note data")
	}
	return Note{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

func (note Note) Display() {
	fmt.Printf("Your note: %s\n Content: %s \n created at: %s", note.Title, note.Content, note.CreatedAt)
}

func ConvertToJson(note Note) []byte {
	json, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	return json
}

func (note Note) Save() error {
	fileName := strings.ReplaceAll(note.Title, " ", "_")
	fileName = strings.ToLower(fileName) + ".json"
	noteJson := ConvertToJson(note)
	return os.WriteFile(fileName, noteJson, 0644)
}

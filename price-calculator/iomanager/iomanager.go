package iomanager

type IOManager interface {
	ReadLines() ([]string, error)
	WriteResults(data interface{}) error
}

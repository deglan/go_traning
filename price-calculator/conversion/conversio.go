package conversion

import (
	"errors"
	"strconv"
)

func StringsToFloats(strings []string) ([]float64, error) {
	var floats []float64
	for _, stringVal := range strings {
		floatPrice, err := strconv.ParseFloat(stringVal, 64)
		if err != nil {
			return nil, errors.New("faild to cenvert string to float")
		}
		floats = append(floats, floatPrice)
	}
	return floats, nil
}

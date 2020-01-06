package main

import (
	"github.com/wardbradt/jsonbldr"
	"log"
)

func main() {
	sampleMap := map[string]string{
		"a": "apple",
		"b": "banana",
		"c": "clementine",
		"d": "durian",
	}
	builder := jsonbldr.New()
	// Setting the second parameter, omitempty, to true causes empty values to be ignored.
	// Setting the third parameter, rawValues, to false causes each value to be wrapped in double quotes
	//   to make it a valid string.
	builder.AddPairs(sampleMap, true, false)
	builder.CloseObject()
	log.Println(string(builder.Bytes()))
	builder.Reset()

	sampleSlice := []string{"apple", "banana", "clementine", "durian"}
	// Setting the second parameter, omitempty, to true causes empty values to be ignored.
	// Setting the third parameter, rawValues, to false causes each value to be wrapped in double quotes
	//   to make it a valid string.
	builder.AddArray("fruits", sampleSlice, true, false)
	builder.CloseObject()
	log.Println(string(builder.Bytes()))
}

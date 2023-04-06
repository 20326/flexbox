package main

import (
	"fmt"
	"log"

	"github.com/20326/flexbox/config"
)

func main() {
	s := config.New("../testdata")

	data := "{\"foo\": \"haha\"}"
	ext := ".json"

	var actual map[string]interface{}
	errs := s.LoadFile("test.json", &actual).End()
	fmt.Printf("actual: %v", actual)

	errs = s.LoadData([]byte(data), ext, &actual).End()

	errs = s.LoadFile("notfound.json", &actual).End()

	if len(errs) > 0 {
		log.Fatalf("errors: %v", errs)
	} else {
		fmt.Printf("actual: %v", actual)
	}
}

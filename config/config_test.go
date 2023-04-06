package config

import (
	"reflect"
	"testing"
)

// go test -c && ./config.test

func TestLoadYamlFile(t *testing.T) {
	s := New("testdata")

	expected := map[string]interface{}{
		"foo": "bar",
	}

	var actual map[string]interface{}
	errs := s.LoadFile("test.yaml", &actual).End()

	if len(errs) > 0 {
		t.Errorf("Unexpected errors: %v", errs)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestLoadJsonFile(t *testing.T) {
	s := New("testdata")

	expected := map[string]interface{}{
		"foo": "bar",
	}

	var actual map[string]interface{}
	errs := s.LoadFile("test.json", &actual).End()

	if len(errs) > 0 {
		t.Errorf("Unexpected errors: %v", errs)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestLoadData(t *testing.T) {
	s := New("testdata")

	expected := map[string]interface{}{
		"foo": "bar",
	}
	data := "{\"foo\": \"bar\"}"
	ext := ".json"
	var actual map[string]interface{}
	errs := s.LoadData([]byte(data), ext, &actual).End()

	if len(errs) > 0 {
		t.Errorf("Unexpected errors: %v", errs)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

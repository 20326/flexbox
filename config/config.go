package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	Version = "0.1.0"
)

func GetVersion() string {
	return Version
}

// A SuperAgent is a object storing all request data for client.
type SuperAgent struct {
	m sync.Mutex

	Debug  bool
	Env    string
	Path   string
	Errors []error
	logger *log.Logger
}

var (
	Cwd string
)

func init() {
	cur, _ := filepath.Abs(os.Args[0])
	Cwd = filepath.Dir(cur)
}

func New(dir string) *SuperAgent {
	// prefer GOLANG_ENV
	err := godotenv.Load()
	if err != nil {
		fmt.Println("not found .env file")
	}

	var env = "local"
	for _, name := range []string{"GOLANG_ENV", "ENV"} {
		if v := os.Getenv(name); v != "" {
			env = v
			log.Printf("GetEnv %s: %s\n", name, v)
			break
		}
	}

	return &SuperAgent{
		m:      sync.Mutex{},
		Debug:  true,
		Env:    env,
		Path:   dir,
		Errors: nil,
		logger: &log.Logger{},
	}
}

func (s *SuperAgent) SetPath(dir string) *SuperAgent {
	s.Path = dir
	return s
}

func (s *SuperAgent) SetLogger(logger *log.Logger) *SuperAgent {
	s.m.Lock()
	defer s.m.Unlock()

	s.logger = logger
	return s
}

func (s *SuperAgent) addError(err error) *SuperAgent {
	if err != nil {
		if s.Errors == nil {
			s.Errors = []error{err}
		} else {
			s.Errors = append(s.Errors, err)
		}
	}
	return s
}

func (s *SuperAgent) LoadFile(name string, object interface{}) *SuperAgent {
	s.m.Lock()
	defer s.m.Unlock()

	// fixme: https://stackoverflow.com/questions/23847003/golang-tests-and-working-directory

	var filename string
	ext := filepath.Ext(name)
	envName := strings.Join([]string{strings.TrimRight(name, ext), ".", s.Env, ext}, "")

	for _, val := range []string{envName, name} {
		filename = path.Join(s.Path, val)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			break
		}
	}

	log.Printf("ReadFile %s\n", filename)
	data, err := os.ReadFile(filename)

	return s.addError(err).load(data, ext, object)
}

func (s *SuperAgent) LoadData(data []byte, ext string, object interface{}) *SuperAgent {
	s.m.Lock()
	defer s.m.Unlock()

	return s.load(data, ext, object)
}

func (s *SuperAgent) End() []error {
	return s.Errors
}

func (s *SuperAgent) load(data []byte, ext string, object interface{}) *SuperAgent {
	var err error
	switch ext {
	case ".json":
		err = json.Unmarshal(data, object)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, object)
	default:
		err = fmt.Errorf("format of %s not supportted", ext)
	}
	// log.Printf("Load object: %v\n", object)
	return s.addError(err)
}

func SelfDir() string {
	return Cwd
}

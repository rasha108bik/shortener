package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// config — name configuration file.
const config = `cmd/staticlint/config.json`

// Checker interface.
type Checker interface {
	readConfig()
	Run()
}

type myChecker struct {
	Staticcheck []string
}

// NewMyChecker returns a newly initialized myChecker objects that implements the Checker
// interface.
func NewMyChecker() *myChecker {
	return &myChecker{}
}

func (m *myChecker) readConfig() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), config))
	if err != nil {
		panic(err)
	}

	cfg := struct {
		Staticcheck []string
	}{}

	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	m.Staticcheck = cfg.Staticcheck
}

// run statick checker.
func (m *myChecker) Run() {
	m.readConfig()
	mychecks := []*analysis.Analyzer{
		// ErrCheckAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	checks := make(map[string]bool)
	for _, v := range m.Staticcheck {
		checks[v] = true
	}
	// add analyse from staticcheck
	for _, v := range staticcheck.Analyzers {
		if checks[v.Name] {
			mychecks = append(mychecks, v)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}

func main() {
	myChecker := NewMyChecker()
	myChecker.Run()
}

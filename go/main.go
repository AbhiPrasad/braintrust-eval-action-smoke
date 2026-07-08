package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type scoreSummary struct {
	Score        float64 `json:"score"`
	Improvements int     `json:"improvements"`
	Regressions  int     `json:"regressions"`
}

type experimentSummary struct {
	ProjectName    string                  `json:"project_name"`
	ExperimentName string                  `json:"experiment_name"`
	ExperimentURL  string                  `json:"experiment_url"`
	Scores         map[string]scoreSummary `json:"scores"`
}

func main() {
	fmt.Fprintln(os.Stderr, "running Go smoke eval")

	// The eval-action Go runtime runs `go run` and parses JSONL experiment summaries
	// from stdout. This synthetic summary exercises that path without depending on
	// a particular Go eval SDK shape.
	summary := experimentSummary{
		ProjectName:    "Smoke Go Eval Action",
		ExperimentName: "go-smoke",
		ExperimentURL:  "https://www.braintrust.dev/",
		Scores: map[string]scoreSummary{
			"exact_match": {
				Score:        1,
				Improvements: 0,
				Regressions:  0,
			},
		},
	}

	b, err := json.Marshal(summary)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

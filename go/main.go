package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	braintrust "github.com/braintrustdata/braintrust-sdk-go"
	"github.com/braintrustdata/braintrust-sdk-go/eval"
	"go.opentelemetry.io/otel/sdk/trace"
)

const projectName = "Smoke Go Eval Action"

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
	ctx := context.Background()
	tp := trace.NewTracerProvider()
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("failed to shut down tracer provider: %v", err)
		}
	}()

	client, err := braintrust.New(
		tp,
		braintrust.WithProject(projectName),
		braintrust.WithBlockingLogin(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	evaluator := braintrust.NewEvaluator[string, string](client)
	result, err := evaluator.Run(ctx, eval.Opts[string, string]{
		ProjectName: projectName,
		Experiment:  "go-smoke",
		Dataset: eval.NewDataset([]eval.Case[string, string]{
			{Input: "Go", Expected: "Hello Go"},
			{Input: "GitHub Actions", Expected: "Hello GitHub Actions"},
			// Intentional mismatch to exercise eval-action regression reporting.
			{Input: "Braintrust", Expected: "Goodbye Braintrust"},
		}),
		Task: eval.T(func(ctx context.Context, input string) (string, error) {
			return "Hello " + input, nil
		}),
		Scorers: []eval.Scorer[string, string]{
			eval.NewScorer("exact_match", func(ctx context.Context, result eval.TaskResult[string, string]) (eval.Scores, error) {
				if result.Output == result.Expected {
					return eval.S(1), nil
				}
				return eval.S(0), nil
			}),
		},
		Quiet: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	permalink, err := result.Permalink()
	if err != nil {
		log.Fatal(err)
	}

	// eval-action's Go runtime runs `go run` and parses Braintrust-style
	// experiment summaries from stdout as JSONL. The Go SDK creates the real
	// Braintrust experiment above; this line gives eval-action the summary shape
	// it needs to render/update the GitHub PR comment.
	summary := experimentSummary{
		ProjectName:    projectName,
		ExperimentName: result.Name(),
		ExperimentURL:  permalink,
		Scores: map[string]scoreSummary{
			"exact_match": {Score: 2.0 / 3.0, Improvements: 0, Regressions: 1},
		},
	}

	b, err := json.Marshal(summary)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

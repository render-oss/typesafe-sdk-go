# TypeSafe Go SDK

A Go SDK for [TypeSafe Jev](https://typesafe.ai) — fast, focused AI judgments in your Go applications.

## Installation

```bash
go get github.com/render-oss/typesafe-sdk-go
```

## Quick Start

```go
import "github.com/render-oss/typesafe-sdk-go/typesafe"

// Create a client
client := typesafe.New(
	typesafe.WithAPIKey("your-api-key"),
)

// Ask a question
resp, err := client.ClassifyString(ctx,
	"What is the size of my database?",
	map[string]typesafe.Question{
		"intent": typesafe.ChoiceQuestion("What is the user asking about?", map[string]string{
			"databases": "Questions about database size, performance, backups",
			"services": "Questions about web services, scaling",
		}),
		"is_urgent": typesafe.NoulQuestion("Is this an urgent problem?"),
	},
)

// Read answers
intent := resp.Answers["intent"].Choice      // "databases"
urgent := resp.Answers["is_urgent"].Noul     // probability 0-1
```

## Primitives

TypeSafe offers three question types:

### Noul (Binary)
Returns probability of yes (0-1):

```go
typesafe.NoulQuestion("Is this a problem?")

// Or with explicit criteria:
typesafe.NoulQuestionWithCriteria(
	"Is this urgent?",
	"Yes: the system is down or data is at risk",
	"No: it's a feature request or informational",
)
```

### Choice (Select One)
Returns the selected option:

```go
typesafe.ChoiceQuestion("Pick the category", map[string]string{
	"bug": "Something is broken",
	"feature": "Request for new capability",
	"docs": "Documentation question",
})
```

### Score (Rated Scale)
Returns probability-weighted score:

```go
rubric, _ := typesafe.NewRubric(
	"Not at all",
	"Somewhat",
	"Very much",
)
typesafe.ScoreQuestion("How satisfied are you?", rubric)
```

## Advanced

### Structured State
Pass any JSON-serializable state:

```go
state := map[string]interface{}{
	"user_question": "Why is my database slow?",
	"context": []string{
		"Last deploy was 2 hours ago",
		"Database size increased 50% this week",
	},
}

resp, err := client.Classify(ctx, state, questions)
```

### Parallel Questions
Ask multiple independent questions in one request:

```go
resp, err := client.Classify(ctx, "What happened to my service?", map[string]typesafe.Question{
	"primary_intent": typesafe.ChoiceQuestion(...),
	"is_troubleshooting": typesafe.NoulQuestion(...),
	"is_urgent": typesafe.NoulQuestion(...),
	"affects_customers": typesafe.NoulQuestion(...),
})
```

All questions run in parallel; answers don't see each other.

### Custom HTTP Client
```go
client := typesafe.New(
	typesafe.WithAPIKey("key"),
	typesafe.WithHTTPClient(&http.Client{
		Timeout: 10 * time.Second,
	}),
)
```

### Retries

Failed requests are retried twice by default, with exponential backoff and jitter.
Retries apply to network errors, 408, 429, and 5xx; `Retry-After` is honored when
present. `typesafe.WithMaxRetries(0)` disables them.

### Logging and Metrics
```go
client := typesafe.New(
	typesafe.WithAPIKey("key"),
	typesafe.WithObserver(func(a typesafe.Attempt) {
		slog.Info("typesafe", "attempt", a.N, "status", a.StatusCode,
			"ms", a.Duration.Milliseconds(), "retrying", a.WillRetry, "err", a.Err)
	}),
)
```

The observer fires once per HTTP attempt, on the calling goroutine.

## Response Structure

```go
type Answer struct {
	Type          Type                // "noul", "choice", "score"
	Noul          float64             // Probability (Noul)
	Choice        string              // Selected option (Choice)
	Score         float64             // Score value (Score)
	Probabilities map[string]float64  // All options' probabilities
	Confidence    float64             // Confidence in the answer (0-1)
}

type Response struct {
	Model   string            // "jev-latest"
	Answers map[string]Answer // Answers keyed by question ID
	Usage   Usage             // Token consumption
}
```

## Examples

See [`examples/`](examples/) for complete examples

## License

[Apache 2.0](LICENSE)

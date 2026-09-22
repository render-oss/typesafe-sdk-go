package typesafe

import (
	"errors"
	"fmt"
)

// Type represents the kind of TypeSafe question.
type Type string

const (
	TypeNoul   Type = "noul"
	TypeChoice Type = "choice"
	TypeScore  Type = "score"
)

// Request is a TypeSafe Jev API request.
type Request struct {
	State     any                 `json:"state"`
	Model     string              `json:"model"`
	Questions map[string]Question `json:"questions"`
}

// Question defines a single question for TypeSafe to evaluate.
type Question struct {
	Type         Type `json:"type"`
	Instructions any  `json:"instructions"`
	Criteria     any  `json:"criteria,omitempty"`
}

// Response is a TypeSafe Jev API response.
type Response struct {
	Model   string            `json:"model"`
	Answers map[string]Answer `json:"answers"`
	Usage   Usage             `json:"usage"`
}

// Answer is TypeSafe's response to a single question.
type Answer struct {
	Type          Type               `json:"type"`
	Noul          float64            `json:"noul,omitempty"`          // TypeNoul: probability of yes (0-1)
	Choice        string             `json:"choice,omitempty"`        // TypeChoice: selected option
	Score         float64            `json:"score,omitempty"`         // TypeScore: score value
	Probabilities map[string]float64 `json:"probabilities,omitempty"` // All options' probabilities
	Confidence    float64            `json:"confidence,omitempty"`    // Confidence in the answer (0-1)
	Legend        map[string]string  `json:"legend,omitempty"`        // TypeScore: level descriptions
}

// Usage tracks token consumption.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Rubric is an ordered list of score levels (2-10).
type Rubric []string

// NewRubric creates a validated Rubric.
func NewRubric(levels ...string) (Rubric, error) {
	if len(levels) < 2 {
		return nil, errors.New("rubric must have at least 2 levels")
	}
	if len(levels) > 10 {
		return nil, fmt.Errorf("rubric accepts at most 10 levels, got %d", len(levels))
	}
	return Rubric(levels), nil
}

// MustRubric is a convenience constructor that panics if validation fails.
func MustRubric(levels ...string) Rubric {
	r, err := NewRubric(levels...)
	if err != nil {
		panic(err)
	}
	return r
}

// NoulQuestion creates a binary yes/no probability evaluation.
func NoulQuestion(instructions string) Question {
	return Question{
		Type:         TypeNoul,
		Instructions: instructions,
	}
}

// NoulQuestionWithCriteria creates a binary yes/no evaluation with explicit criteria.
func NoulQuestionWithCriteria(instructions, trueCriteria, falseCriteria string) Question {
	return Question{
		Type:         TypeNoul,
		Instructions: instructions,
		Criteria: map[string]string{
			"true":  trueCriteria,
			"false": falseCriteria,
		},
	}
}

// ChoiceQuestion creates a discrete option classification (1+ options).
// criteria maps each option label to its definition/examples.
func ChoiceQuestion(instructions string, criteria map[string]string) Question {
	return Question{
		Type:         TypeChoice,
		Instructions: instructions,
		Criteria:     criteria,
	}
}

// ScoreQuestion creates an ordered rubric rating (2-10 levels).
// rubric is an ordered list describing each level from lowest to highest.
func ScoreQuestion(instructions string, rubric Rubric) Question {
	return Question{
		Type:         TypeScore,
		Instructions: instructions,
		Criteria:     []string(rubric),
	}
}

// StructuredQuestion creates a question with structured instructions and/or criteria.
// This is useful for complex questions that need more than a string.
func StructuredQuestion(qType Type, instructions any, criteria any) Question {
	return Question{
		Type:         qType,
		Instructions: instructions,
		Criteria:     criteria,
	}
}

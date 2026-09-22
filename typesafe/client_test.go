package typesafe

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	c := New(
		WithAPIKey("test-key"),
		WithBaseURL("https://api.typesafe.ai/v1/systemone"),
		WithModel("jev-latest"),
	)

	if c.APIKey != "test-key" {
		t.Errorf("expected APIKey=test-key, got %s", c.APIKey)
	}
	if c.BaseURL != "https://api.typesafe.ai/v1/systemone" {
		t.Errorf("expected BaseURL=https://api.typesafe.ai/v1/systemone, got %s", c.BaseURL)
	}
	if c.Model != "jev-latest" {
		t.Errorf("expected Model=jev-latest, got %s", c.Model)
	}
}

func TestNoulQuestion(t *testing.T) {
	q := NoulQuestion("Is this a yes/no question?")
	if q.Type != TypeNoul {
		t.Errorf("expected TypeNoul, got %s", q.Type)
	}
	if q.Instructions != "Is this a yes/no question?" {
		t.Errorf("expected instructions to be set")
	}
}

func TestChoiceQuestion(t *testing.T) {
	criteria := map[string]string{
		"option1": "First option",
		"option2": "Second option",
	}
	q := ChoiceQuestion("Pick one", criteria)
	if q.Type != TypeChoice {
		t.Errorf("expected TypeChoice, got %s", q.Type)
	}
	if q.Criteria == nil {
		t.Errorf("expected criteria to be set")
	}
}

func TestRubric(t *testing.T) {
	rubric, err := NewRubric("Low", "Medium", "High")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rubric) != 3 {
		t.Errorf("expected 3 levels, got %d", len(rubric))
	}

	// Test validation
	_, err = NewRubric("Only one")
	if err == nil {
		t.Error("expected error for rubric with 1 level")
	}

	_, err = NewRubric(
		"1", "2", "3", "4", "5",
		"6", "7", "8", "9", "10", "11",
	)
	if err == nil {
		t.Error("expected error for rubric with >10 levels")
	}
}

func TestRequest_Marshal(t *testing.T) {
	req := Request{
		State: "test question",
		Model: "jev-latest",
		Questions: map[string]Question{
			"test": NoulQuestion("Is this a test?"),
		},
	}

	// Should not panic when marshaling
	_, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}
}

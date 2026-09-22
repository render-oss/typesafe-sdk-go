package typesafe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestRetryAndObserver(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(Response{Model: "jev-latest"})
	}))
	defer srv.Close()

	var seen []Attempt
	c := New(
		WithBaseURL(srv.URL),
		WithObserver(func(a Attempt) { seen = append(seen, a) }),
	)

	resp, err := c.Classify(context.Background(), "hi", map[string]Question{
		"q": NoulQuestion("Is this a test?"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Model != "jev-latest" {
		t.Errorf("expected model jev-latest, got %s", resp.Model)
	}
	if calls != 3 {
		t.Errorf("expected 3 attempts, got %d", calls)
	}
	if len(seen) != 3 {
		t.Fatalf("expected 3 observed attempts, got %d", len(seen))
	}
	if !seen[0].WillRetry || seen[0].StatusCode != http.StatusServiceUnavailable {
		t.Errorf("first attempt should be a retried 503, got %+v", seen[0])
	}
	if seen[2].Err != nil || seen[2].WillRetry {
		t.Errorf("last attempt should have succeeded, got %+v", seen[2])
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	_, err := New(WithBaseURL(srv.URL)).Classify(context.Background(), "hi", nil)
	if err == nil {
		t.Fatal("expected an error")
	}
	if calls != 1 {
		t.Errorf("expected 1 attempt for a 400, got %d", calls)
	}
}

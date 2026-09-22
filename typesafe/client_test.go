package typesafe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c := New(
		WithAPIKey("test-key"),
		WithBaseURL("https://api.typesafe.ai/v1/systemone"),
		WithModel("jev-latest"),
	)

	require.Equal(t, "test-key", c.APIKey)
	require.Equal(t, "https://api.typesafe.ai/v1/systemone", c.BaseURL)
	require.Equal(t, "jev-latest", c.Model)
}

func TestNoulQuestion(t *testing.T) {
	q := NoulQuestion("Is this a yes/no question?")
	require.Equal(t, TypeNoul, q.Type)
	require.Equal(t, "Is this a yes/no question?", q.Instructions)
}

func TestChoiceQuestion(t *testing.T) {
	criteria := map[string]string{
		"option1": "First option",
		"option2": "Second option",
	}
	q := ChoiceQuestion("Pick one", criteria)
	require.Equal(t, TypeChoice, q.Type)
	require.Equal(t, criteria, q.Criteria)
}

func TestRubric(t *testing.T) {
	rubric, err := NewRubric("Low", "Medium", "High")
	require.NoError(t, err)
	require.Len(t, rubric, 3)

	_, err = NewRubric("Only one")
	require.Error(t, err)

	_, err = NewRubric(
		"1", "2", "3", "4", "5",
		"6", "7", "8", "9", "10", "11",
	)
	require.Error(t, err)
}

func TestRequest_Marshal(t *testing.T) {
	req := Request{
		State: "test question",
		Model: "jev-latest",
		Questions: map[string]Question{
			"test": NoulQuestion("Is this a test?"),
		},
	}

	_, err := json.Marshal(req)
	require.NoError(t, err)
}

func TestRetryAndObserver(t *testing.T) {
	// Encoded up front so the handler has nothing left to assert on: require's
	// FailNow is only valid on the goroutine running the test.
	body, err := json.Marshal(Response{Model: "jev-latest"})
	require.NoError(t, err)

	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	var seen []Attempt
	c := New(
		WithBaseURL(srv.URL),
		WithObserver(func(a Attempt) { seen = append(seen, a) }),
	)

	resp, err := c.Classify(t.Context(), "hi", map[string]Question{
		"q": NoulQuestion("Is this a test?"),
	})
	require.NoError(t, err)
	require.Equal(t, "jev-latest", resp.Model)
	require.Equal(t, 3, calls)
	require.Len(t, seen, 3)

	require.True(t, seen[0].WillRetry)
	require.Equal(t, http.StatusServiceUnavailable, seen[0].StatusCode)
	require.NoError(t, seen[2].Err)
	require.False(t, seen[2].WillRetry)
}

func TestNoRetryOn4xx(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	_, err := New(WithBaseURL(srv.URL)).Classify(t.Context(), "hi", nil)
	require.Error(t, err)
	require.Equal(t, 1, calls)
}

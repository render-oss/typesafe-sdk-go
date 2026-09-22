package main

import (
	"context"
	"fmt"
	"os"

	"github.com/renderinc/typesafe-go/typesafe"
)

// ReviewRequest represents a code review submission
type ReviewRequest struct {
	PR       int
	Title    string
	Body     string
	Author   string
	FileSize int // approximate LOC
}

// prioritizeReview uses TypeSafe to triage code review work
func prioritizeReview(ctx context.Context, client *typesafe.Client, req ReviewRequest) {
	content := fmt.Sprintf("PR Title: %s\n\nDescription: %s", req.Title, req.Body)

	resp, err := client.ClassifyString(ctx, content, map[string]typesafe.Question{
		// Determine urgency for review
		"urgency": typesafe.ScoreQuestion(
			"How urgently should this PR be reviewed?",
			typesafe.MustRubric(
				"Low: refactoring, docs, or non-critical updates",
				"Medium: feature or bug fix, can be reviewed in normal queue",
				"High: blocking other work or fixing critical issue",
				"Critical: production hotfix or blocking deployment",
			),
		),
		// Identify primary change type
		"change_type": typesafe.ChoiceQuestion(
			"What is the primary nature of this change?",
			map[string]string{
				"bugfix":         "Fixes an existing bug or defect",
				"feature":        "Adds new functionality",
				"refactor":       "Improves code structure without behavior change",
				"documentation":  "Documentation, comments, or non-code changes",
				"tests":          "Test additions or improvements",
				"infrastructure": "Build, CI/CD, or deployment changes",
			},
		),
		// Risk assessment
		"risk_level": typesafe.ChoiceQuestion(
			"What is the risk level of this change?",
			map[string]string{
				"low":      "Isolated change, well-tested, easy to revert",
				"medium":   "Affects multiple areas but with safety checks",
				"high":     "Touches core logic or affects many systems",
				"critical": "Could cause data loss or system failure if wrong",
			},
		),
		// Suggest reviewer expertise
		"reviewer_type": typesafe.ChoiceQuestion(
			"What type of reviewer expertise is most important?",
			map[string]string{
				"junior":       "Good learning opportunity for newer developers",
				"general":      "Any experienced developer can review",
				"specialist":   "Needs domain expert (database, security, etc.)",
				"architecture": "Requires architectural knowledge of the system",
			},
		),
	})
	if err != nil {
		fmt.Printf("Error reviewing PR %d: %v\n", req.PR, err)
		return
	}

	urgency := resp.Answers["urgency"].Score
	changeType := resp.Answers["change_type"].Choice
	riskLevel := resp.Answers["risk_level"].Choice
	reviewerType := resp.Answers["reviewer_type"].Choice

	fmt.Printf("\n🔍 PR #%d: %s\n", req.PR, req.Title)
	fmt.Printf("Author: %s | Size: ~%d LOC\n\n", req.Author, req.FileSize)
	fmt.Printf("Type: %s | Risk: %s\n", changeType, riskLevel)
	fmt.Printf("Reviewer: %s expert\n", reviewerType)

	// Urgency-based routing
	switch urgency {
	case 4:
		fmt.Printf("🚨 CRITICAL: Review immediately - blocking or hotfix\n")
	case 3:
		fmt.Printf("⏱️  HIGH: Prioritize in today's queue\n")
	case 2:
		fmt.Printf("📋 MEDIUM: Normal review queue\n")
	default:
		fmt.Printf("✅ LOW: Can be batched or reviewed when bandwidth available\n")
	}

	// Risk/change-type pairing
	if riskLevel == "critical" {
		fmt.Println("⚠️  CRITICAL RISK: Requires second review")
	} else if riskLevel == "high" && changeType != "tests" {
		fmt.Println("⚠️  HIGH RISK: Ensure thorough review before merge")
	}
}

func main() {
	apiKey := os.Getenv("TYPESAFE_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: TYPESAFE_API_KEY environment variable not set")
		os.Exit(1)
	}

	client := typesafe.New(typesafe.WithAPIKey(apiKey))
	ctx := context.Background()

	// Example PRs to triage
	prs := []ReviewRequest{
		{
			PR:       1234,
			Title:    "Fix: database connection leak in transaction handler",
			Body:     "Fixes issue #456 where long-running processes were exhausting connection pool. Added proper cleanup in defer block.",
			Author:   "alice",
			FileSize: 45,
		},
		{
			PR:       1235,
			Title:    "Refactor: extract logging utilities into separate package",
			Body:     "Moving logging code to its own package for better separation of concerns. No behavior change.",
			Author:   "bob",
			FileSize: 120,
		},
		{
			PR:       1236,
			Title:    "Feature: add user email verification flow",
			Body:     "Implements email verification for new signups. Adds new tables, middleware, and email service integration.",
			Author:   "charlie",
			FileSize: 280,
		},
		{
			PR:       1237,
			Title:    "Hotfix: production memory leak in cache eviction",
			Body:     "Emergency fix for production issue causing OOM. Adds missing cleanup in cache eviction logic.",
			Author:   "diana",
			FileSize: 15,
		},
	}

	fmt.Println("📊 Code Review Triage with TypeSafe")
	fmt.Println(string(make([]byte, 60)))

	for _, pr := range prs {
		prioritizeReview(ctx, client, pr)
	}
}

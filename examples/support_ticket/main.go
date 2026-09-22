package main

import (
	"context"
	"fmt"
	"os"

	"github.com/renderinc/typesafe-go/typesafe"
)

// SupportTicket represents a customer support inquiry
type SupportTicket struct {
	ID      string
	Message string
}

// analyzeTicket uses TypeSafe to evaluate a support ticket across multiple dimensions
func analyzeTicket(ctx context.Context, client *typesafe.Client, ticket SupportTicket) {
	resp, err := client.ClassifyString(ctx, ticket.Message, map[string]typesafe.Question{
		// Route to the appropriate department
		"department": typesafe.ChoiceQuestion(
			"Which team should handle this support request?",
			map[string]string{
				"technical":  "Technical issues, bugs, integration failures, error messages",
				"billing":    "Payment problems, invoicing, subscription questions, pricing",
				"sales":      "Feature requests, pre-sales questions, partnership inquiries",
				"general":    "Account questions, documentation, product information",
			},
		),
		// Detect urgency / severity
		"urgency": typesafe.ScoreQuestion(
			"How urgent is this issue?",
			typesafe.MustRubric(
				"Low: cosmetic issue or general question",
				"Medium: feature not working but workaround exists",
				"High: major feature down, workaround is difficult",
				"Critical: system is completely broken, data loss risk",
			),
		),
		// Measure customer sentiment / frustration
		"sentiment": typesafe.NoulQuestionWithCriteria(
			"Is the customer expressing frustration or anger?",
			"Yes: frustrated, angry, using urgent language or exclamation marks",
			"No: neutral, professional, or satisfied tone",
		),
		// Check if customer is at churn risk
		"churn_risk": typesafe.NoulQuestionWithCriteria(
			"Is the customer expressing intention to leave or look for alternatives?",
			"Yes: mentions canceling, switching, or stopping use due to this issue",
			"No: no indication of leaving",
		),
	})

	if err != nil {
		fmt.Printf("Error analyzing ticket %s: %v\n", ticket.ID, err)
		return
	}

	// Extract and display results
	department := resp.Answers["department"].Choice
	urgency := resp.Answers["urgency"].Score
	sentiment := resp.Answers["sentiment"].Noul
	churnRisk := resp.Answers["churn_risk"].Noul

	fmt.Printf("\n=== Ticket %s ===\n", ticket.ID)
	fmt.Printf("Customer: %q\n\n", ticket.Message)
	fmt.Printf("📋 Department: %s\n", department)
	fmt.Printf("⚠️  Urgency: %.0f/4 (%.0f%%)\n", urgency, urgency*25)
	fmt.Printf("😠 Customer Frustrated: %.0f%%\n", sentiment*100)
	fmt.Printf("🚪 Churn Risk: %.0f%%\n\n", churnRisk*100)

	// Recommend handling priority
	if churnRisk > 0.7 || urgency == 4 {
		fmt.Println("🔴 PRIORITY: Escalate immediately — high churn/critical issue")
	} else if urgency >= 3 || sentiment > 0.6 {
		fmt.Println("🟠 HIGH: Assign to senior team member")
	} else if urgency >= 2 {
		fmt.Println("🟡 MEDIUM: Normal queue")
	} else {
		fmt.Println("🟢 LOW: Can be batched or automated")
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

	// Example support tickets
	tickets := []SupportTicket{
		{
			ID: "TICKET-001",
			Message: "Hi, I've been trying to connect my Stripe account for 3 days and the integration keeps failing. I'm losing sales.",
		},
		{
			ID: "TICKET-002",
			Message: "Could you explain how to configure webhooks? I'm reading the docs but want to make sure I understand correctly.",
		},
		{
			ID: "TICKET-003",
			Message: "Your service is completely broken. I've wasted hours on this. If this isn't fixed by tomorrow, we're switching providers.",
		},
		{
			ID: "TICKET-004",
			Message: "What's the difference between the Pro and Enterprise plans?",
		},
	}

	fmt.Println("🎫 Support Ticket Analysis with TypeSafe")
	fmt.Println(string(make([]byte, 60)))

	for _, ticket := range tickets {
		analyzeTicket(ctx, client, ticket)
	}
}

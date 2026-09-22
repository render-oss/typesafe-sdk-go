package main

import (
	"context"
	"fmt"
	"os"

	"github.com/renderinc/typesafe-go/typesafe"
)

// ModerateContent evaluates user-generated content for publication safety
func moderateContent(ctx context.Context, client *typesafe.Client, content string, username string) {
	resp, err := client.ClassifyString(ctx, content, map[string]typesafe.Question{
		// Detect if content is appropriate for public publication
		"is_appropriate": typesafe.NoulQuestionWithCriteria(
			"Is this content appropriate for public publication?",
			"Yes: safe, constructive, follows community guidelines",
			"No: contains harassment, explicit content, hate speech, spam, or abuse",
		),
		// Categorize the type of content
		"content_type": typesafe.ChoiceQuestion(
			"What type of content is this?",
			map[string]string{
				"discussion":  "Thoughtful discussion, questions, or sharing knowledge",
				"feedback":    "Bug reports, feature requests, or constructive criticism",
				"spam":        "Promotional, advertising, or off-topic content",
				"harassment":  "Personal attacks, insults, or abusive language",
				"misinformation": "False claims or misleading information",
			},
		),
		// Check if it needs human review
		"needs_review": typesafe.NoulQuestionWithCriteria(
			"Should this content be reviewed by a human moderator?",
			"Yes: borderline, complex context needed, or subjective decision",
			"No: clearly appropriate or clearly violates policy",
		),
		// Confidence in the judgment
		"confidence": typesafe.ScoreQuestion(
			"How confident is the moderation decision?",
			typesafe.MustRubric(
				"Very uncertain: complex context, judgment call",
				"Somewhat uncertain: likely clear but could have edge cases",
				"Confident: straightforward policy violation or approval",
				"Very confident: obvious violation or clear safe content",
			),
		),
	})

	if err != nil {
		fmt.Printf("Error moderating content: %v\n", err)
		return
	}

	isAppropriate := resp.Answers["is_appropriate"].Noul
	contentType := resp.Answers["content_type"].Choice
	needsReview := resp.Answers["needs_review"].Noul
	confidence := resp.Answers["confidence"].Score

	// Display results
	fmt.Printf("\n📝 %s: %q\n", username, content)
	fmt.Printf("Type: %s | Confidence: %.0f%%\n", contentType, confidence*25)

	if isAppropriate > 0.7 && needsReview < 0.4 {
		fmt.Println("✅ APPROVED: Safe to publish")
	} else if needsReview > 0.6 {
		fmt.Println("🔶 HOLD: Needs human review")
	} else if isAppropriate < 0.3 {
		fmt.Println("❌ REJECTED: Violates community guidelines")
	} else {
		fmt.Println("⚠️  BORDERLINE: Likely needs human review")
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

	// Test cases
	testCases := []struct {
		username string
		content  string
	}{
		{
			"alice",
			"Great article on Go concurrency! I found the channel examples really helpful for my project.",
		},
		{
			"bob",
			"This product is garbage, worst purchase ever!!!",
		},
		{
			"charlie",
			"BUY MY CRYPTO NOW!!! Click here for guaranteed returns!!! 🚀🚀🚀",
		},
		{
			"diana",
			"How do I optimize database queries? I'm seeing slow performance on large tables.",
		},
		{
			"eve",
			"Your service sucks and you should all be ashamed of yourselves.",
		},
	}

	fmt.Println("🛡️  Content Moderation with TypeSafe")
	fmt.Println(string(make([]byte, 60)))

	for _, test := range testCases {
		moderateContent(ctx, client, test.content, test.username)
	}
}

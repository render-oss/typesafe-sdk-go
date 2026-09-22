# TypeSafe Go SDK Examples

Three real-world examples demonstrating TypeSafe for different use cases.

## 1. Support Ticket Routing

**Problem**: Support tickets arrive with varying urgency. Manual routing wastes time.

**Solution**: TypeSafe analyzes each ticket to identify:
- Which department should handle it (technical, billing, sales, general)
- How urgent it is (1-4 scale)
- Whether the customer is frustrated
- Churn risk

**Run it**:
```bash
cd support_ticket
export TYPESAFE_API_KEY=your_key
go run main.go
```

**Output**: Routing recommendations with priority levels 🔴 🟠 🟡 🟢

**Key insight**: Combines multiple signals (urgency + sentiment + risk) to make routing decisions better than any single signal alone.

---

## 2. Content Moderation

**Problem**: Community platforms need to moderate at scale without being error-prone.

**Solution**: TypeSafe evaluates user content to identify:
- Whether it's safe to publish
- What type of content (discussion, feedback, spam, harassment, misinformation)
- Whether it needs human review
- Confidence in the decision

**Run it**:
```bash
cd content_moderation
export TYPESAFE_API_KEY=your_key
go run main.go
```

**Output**: Moderation decisions with confidence flags

**Key insight**: Context matters. "Your service sucks" is harassment; "Your service needs X" is feedback. TypeSafe understands the difference.

---

## 3. Code Review Triage

**Problem**: Code reviews pile up. Urgent hotfixes wait behind refactoring.

**Solution**: TypeSafe analyzes PRs to recommend:
- How urgently it needs review (1-4 scale)
- What kind of change (bugfix, feature, refactor, etc.)
- Risk level (low to critical)
- What expertise is needed (junior, general, specialist, architecture)

**Run it**:
```bash
cd code_review
export TYPESAFE_API_KEY=your_key
go run main.go
```

**Output**: Triage recommendations with reviewer suggestions

**Key insight**: Combines risk and change type to route critical fixes urgently while giving junior devs learning opportunities.

---

## Common Pattern

All three examples use the same pattern:

1. **Define questions** — What do you need to know?
2. **Ask in parallel** — All questions run at once
3. **Compose answers** — Use results to make decisions

Example:
```go
resp, _ := client.Classify(ctx, input, map[string]typesafe.Question{
    "category": typesafe.Choice(...),     // Route to handler
    "urgency": typesafe.Score(...),       // Set priority
    "needs_review": typesafe.Noul(...),   // Human escalation
})

// Compose results into actions
switch resp.Answers["category"].Choice {
case "urgent":
    // handle urgently
}
```

## When to Use TypeSafe

✅ Good fits:
- Classification with multiple dimensions
- Understanding intent or severity
- Routing or triage decisions
- Decision-making that benefits from context

❌ Not ideal for:
- Generating new text (use LLMs instead)
- Code execution or math (use functions)
- Deterministic logic (use if/then rules)

## SDK Documentation

See the main [README.md](../README.md) for:
- Installation
- API reference
- Question types (Noul, Choice, Score)
- Advanced features

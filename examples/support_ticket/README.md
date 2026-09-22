# Support Ticket Routing Example

This example demonstrates using TypeSafe to analyze customer support tickets and route them appropriately.

## The Problem

Support teams receive hundreds of tickets with varying urgency and complexity:
- "I'm losing sales because integration is down"
- "What's your pricing?"
- "I'm very unhappy with your service"

Manual routing wastes time. Automated routing needs to understand context.

## The Solution

Ask TypeSafe **4 parallel questions** about each ticket:

1. **department** (Choice): Which team should handle it?
   - Technical (bugs, integration failures)
   - Billing (payments, invoicing)
   - Sales (features, partnerships)
   - General (account, documentation)

2. **urgency** (Score): How urgent is it? (1-4 scale)
   - Low → can wait
   - Medium → normal queue
   - High → prioritize today
   - Critical → handle immediately

3. **sentiment** (Noul): Is the customer frustrated?
   - Yes → churn risk
   - No → routine handling

4. **churn_risk** (Noul): Might they leave?
   - Yes → escalate
   - No → standard process

## Running the Example

```bash
export TYPESAFE_API_KEY=your_key_here
go run ./examples/support_ticket/main.go
```

## Output

```
🎫 Support Ticket Analysis with TypeSafe

=== Ticket TICKET-001 ===
Customer: "Hi, I've been trying to connect my Stripe account..."

📋 Department: technical
⚠️  Urgency: 4/4 (100%)
😠 Customer Frustrated: 95%
🚪 Churn Risk: 87%

🔴 PRIORITY: Escalate immediately — high churn/critical issue
```

## Routing Decisions

The example combines all answers to recommend handling:

- **🔴 CRITICAL**: Churn risk > 70% OR urgency = 4 → Escalate immediately
- **🟠 HIGH**: Urgency ≥ 3 OR frustrated customer → Senior team member
- **🟡 MEDIUM**: Urgency ≥ 2 → Normal queue
- **🟢 LOW**: Everything else → Can be batched or automated

## Real-World Benefits

✅ **Consistency**: Same rubric for all tickets
✅ **Speed**: Seconds vs. minutes per ticket
✅ **Churn prevention**: Catches frustrated customers at risk
✅ **Team efficiency**: Specialists get complex issues, junior devs handle simple ones
✅ **SLA compliance**: Urgency flagging ensures critical issues are handled quickly

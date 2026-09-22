# Code Review Triage Example

This example demonstrates using TypeSafe to automatically triage code reviews and route them to the right reviewer.

## The Problem

Engineering teams struggle with code review bottlenecks:
- Urgent hotfixes wait behind refactoring
- Junior devs miss learning opportunities
- Specialists get pulled into trivial reviews
- Risk is hard to assess from PR title alone

## The Solution

Analyze PRs to guide reviewer assignment:

1. **urgency** (Score): How quickly should it be reviewed?
   - Low, Medium, High, Critical
2. **change_type** (Choice): What kind of change?
   - Bugfix, feature, refactor, docs, tests, infrastructure
3. **risk_level** (Choice): How risky is it?
   - Low, Medium, High, Critical
4. **reviewer_type** (Choice): What expertise is needed?
   - Junior (learning opportunity), General, Specialist, Architecture

## Running the Example

```bash
export TYPESAFE_API_KEY=your_key_here
go run ./examples/code_review/main.go
```

## Output

```
📊 Code Review Triage with TypeSafe

🔍 PR #1234: Fix: database connection leak in transaction handler
Author: alice | Size: ~45 LOC

Type: bugfix | Risk: high
Reviewer: specialist expert
⏱️  HIGH: Prioritize in today's queue
⚠️  HIGH RISK: Ensure thorough review before merge
```

## Routing Logic

- **🚨 CRITICAL**: Urgency = 4 → Review immediately
- **⏱️ HIGH**: Urgency = 3 → Prioritize today
- **📋 MEDIUM**: Urgency = 2 → Normal queue
- **✅ LOW**: Urgency = 1 → Batch or when bandwidth available

**Risk Escalations:**
- Critical risk → Requires second review
- High risk + not tests → Ensure thorough review

**Reviewer Matching:**
- Junior: Good for learning, features, simple refactors
- Specialist: Security, database, infrastructure changes
- Architecture: Major refactors, design changes

## Real-World Benefits

✅ **Speed**: Faster time-to-merge for low-risk changes
✅ **Efficiency**: Specialists focus on complex reviews
✅ **Learning**: Identify good PRs for junior developers
✅ **Risk management**: High-risk changes get extra scrutiny
✅ **Queue management**: Automatic triage prevents bottlenecks
✅ **SLA compliance**: Critical hotfixes surface automatically

## Integration Ideas

- **Slack notifications**: "PR #1234 needs architecture review urgently"
- **CI pipeline**: Block merge if critical risk detected
- **Dashboard**: Visibility into queue composition and urgency
- **Analytics**: Track review times by type and risk level

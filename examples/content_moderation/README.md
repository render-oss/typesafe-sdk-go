# Content Moderation Example

This example demonstrates using TypeSafe to moderate user-generated content at scale.

## The Problem

Community platforms need to balance safety with efficiency:
- Manual moderation doesn't scale
- Simple keyword filtering is error-prone
- Context matters: criticism ≠ harassment

## The Solution

Evaluate content across multiple dimensions:

1. **is_appropriate** (Noul): Is it safe to publish?
2. **content_type** (Choice): What kind of content?
   - Discussion, feedback, spam, harassment, misinformation
3. **needs_review** (Noul): Should a human look at it?
4. **confidence** (Score): How confident in the decision?

## Running the Example

```bash
export TYPESAFE_API_KEY=your_key_here
go run ./examples/content_moderation/main.go
```

## Output

```
🛡️  Content Moderation with TypeSafe

📝 alice: "Great article on Go concurrency!..."
Type: discussion | Confidence: 100%
✅ APPROVED: Safe to publish

📝 eve: "Your service sucks and you should all be ashamed..."
Type: harassment | Confidence: 100%
❌ REJECTED: Violates community guidelines
```

## Decision Logic

- **✅ APPROVED**: Appropriate > 70% AND needs_review < 40%
- **🔶 HOLD**: Needs human review > 60%
- **❌ REJECTED**: Appropriate < 30%
- **⚠️ BORDERLINE**: Everything else

## Real-World Benefits

✅ **Scale**: Analyze thousands of posts per hour
✅ **Consistency**: Same standards across all moderation
✅ **Context-aware**: Understands nuance and intent
✅ **Confidence flagging**: Routes borderline cases to humans
✅ **Compliance**: Audit trail of moderation decisions

## Privacy Note

Content is analyzed by TypeSafe; consider data residency requirements for sensitive communities.

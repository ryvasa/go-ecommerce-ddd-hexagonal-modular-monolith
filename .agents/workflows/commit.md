---
description: stage, review, and commit changes with conventional commits
---

# Git Commit

Workflow to stage, review diffs, and commit changes using Conventional Commits.

---

## 1. Stage Changes

// turbo

```bash
git add -A && git diff --cached --stat
```

Display a summary of changed files.

---

## 2. Review Diff

// turbo

```bash
git diff --cached
```

Review the changes to be committed. Ensure:

- No files that should not be committed are included (e.g. `.env`, `tmp/`)
- All changes are relevant to the current task
- No debug code or leftover `fmt.Println` statements

---

## 3. Compose Commit Message

Create a commit message following the Conventional Commits format:

```
<type>(<scope>): <short description>.
```

**Allowed types**: `feat`, `fix`, `refactor`, `chore`, `docs`, `style`, `perf`, `test`

**Scope**: name of the affected module (e.g. `user`, `auth`, `cart`, `shared`, `migration`)

**Rules**:
- Description must be lowercase
- Must end with a period (.)
- Maximum 60 characters

Examples:
```
feat(user): a new endpoint to create user.
fix(auth): handle expired token gracefully.
refactor(cart): simplify cart item validation.
chore(migration): add orders table migration.
```

---

## 4. Commit

**[APPROVAL REQUIRED]** — Present the commit message to the user and request confirmation before running:

```bash
git commit -m "<commit message>"
```

---

## 5. Verify Commit

// turbo

```bash
git log -1 --oneline
```

Ensure the commit has been recorded correctly.

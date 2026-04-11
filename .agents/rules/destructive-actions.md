---
trigger: always_on
---

# Destructive Action Rules

## MANDATORY APPROVAL

**ALL destructive actions MUST receive explicit user approval BEFORE execution.**

Any action that removes, deletes, or drops something is considered destructive:

### File System

- Deleting any file or directory
- Overwriting an existing file with significantly different content
- Removing code blocks, functions, types, or interfaces from source files

### Database

- Dropping a table
- Deleting records (even in dev/local)
- Removing fields from existing models or schemas
- Running destructive migrations (down migrations)

### Architecture

- Removing a module from `internal/`
- Deleting domain entities or value objects
- Removing port interfaces or changing their signatures
- Removing Wire bindings from `module.go` or `wire.go`

### Git

- Force pushing (`git push --force`)
- Resetting commits (`git reset --hard`)
- Deleting branches

### Dependencies

- Removing packages from `go.mod`
- Deleting skill, rule, or workflow files from `.agent/`

## Workflow

1. **Identify** the destructive action and list exactly what will be removed
2. **Present** the list to the user with clear explanation of impact
3. **Wait** for explicit approval
4. **Execute** only after approval is received
5. **Never** batch destructive actions silently within a larger task

## Zero Tolerance

- NEVER delete files without listing them first and getting approval
- NEVER drop tables without approval
- NEVER remove domain interfaces/entities without approval
- If in doubt whether an action is destructive, **ask first**

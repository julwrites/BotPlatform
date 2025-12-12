# Task Documentation System Guide

This guide explains how to create, maintain, and update task documentation. It provides a reusable system for tracking implementation work, decisions, and progress.

## Core Philosophy
**"If it's not documented in `docs/tasks/`, it didn't happen."**

## Directory Structure
Tasks are organized by category in `docs/tasks/`:
- `foundation/`: Core architecture and setup
- `infrastructure/`: Services, adapters, platform code
- `domain/`: Business logic, use cases
- `presentation/`: UI, state management
- `features/`: End-to-end feature implementation
- `migration/`: Refactoring, upgrades
- `testing/`: Testing infrastructure
- `review/`: Code reviews and PR analysis

## Task Document Format

We use **YAML Frontmatter** for metadata and **Markdown** for content.

### Frontmatter (Required)
```yaml
---
id: FOUNDATION-20250521-103000   # Auto-generated Timestamp ID
status: pending                  # Current status
title: Initial Project Setup     # Task Title
priority: medium                 # high, medium, low
created: 2025-05-21 10:30:00     # Creation timestamp
category: foundation             # Category
type: task                       # task, story, bug, epic (Optional)
sprint: Sprint 1                 # Iteration identifier (Optional)
estimate: 3                      # Story points / T-shirt size (Optional)
dependencies: [TASK-001, TASK-002] # Dependency list
---
```

### Status Workflow
1. `pending`: Created but not started.
2. `in_progress`: Active development.
3. `review_requested`: Implementation done, awaiting code review.
4. `verified`: Reviewed and approved.
5. `completed`: Merged and finalized.
6. `wip_blocked` / `blocked`: Development halted.
7. `cancelled` / `deferred`: Stopped or postponed.

### Content Template
```markdown
# [Task Title]

[Description of what needs to be done]

### Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Implementation Status
### Completed Work
- ✅ Implemented X (file.py)

### Blockers
[Describe blockers if any]
```

## Tools

Use the `scripts/tasks` wrapper to manage tasks.

```bash
# Create a new task (standard)
./scripts/tasks create foundation "Task Title"

# Create an Agile Story in a Sprint
./scripts/tasks create features "User Login" --type story --sprint "Sprint 1" --estimate 5

# List tasks (can filter by sprint)
./scripts/tasks list
./scripts/tasks list --sprint "Sprint 1"
./scripts/tasks list --archived

# Show task details
./scripts/tasks show [TASK_ID]

# Find the next best task to work on (Smart Agent Mode)
./scripts/tasks next

# Update status
./scripts/tasks update [TASK_ID] in_progress
./scripts/tasks update [TASK_ID] review_requested
./scripts/tasks update [TASK_ID] verified
./scripts/tasks update [TASK_ID] completed
./scripts/tasks complete [TASK_ID]  # Shortcut for completed

# Managing Dependencies
./scripts/tasks link [TASK_ID] [DEPENDENCY_ID]   # Add dependency
./scripts/tasks unlink [TASK_ID] [DEPENDENCY_ID] # Remove dependency
./scripts/tasks index                            # Generate INDEX.yaml

# Visualization & Validation
./scripts/tasks visualize (or graph)             # Generate Mermaid graph
./scripts/tasks validate                         # Check for errors/cycles

# Maintenance
./scripts/tasks migrate                          # Migrate legacy tasks
./scripts/tasks archive [TASK_ID]                # Move to archive
./scripts/tasks delete [TASK_ID]                 # Delete task
./scripts/tasks install-hooks                    # Install git pre-commit hook
./scripts/tasks init                             # Initialize directory structure
```

## Dependency Management

### Linking Tasks
Tasks can have dependencies. Use `link` to establish a relationship.
```bash
./scripts/tasks link TASK-B TASK-A  # Task B depends on Task A
```
This updates the `dependencies` field in the frontmatter.

### The Index
The system maintains a dependency index at `docs/tasks/INDEX.yaml`.
Run `./scripts/tasks index` to regenerate it. This file provides a machine-readable map of all tasks and their relationships.

## Agile Methodology

This system supports Agile/Scrum workflows for LLM-Human collaboration.

### Sprints
- Tag tasks with `sprint: [Name]` to group them into iterations.
- Use `./scripts/tasks list --sprint [Name]` to view the sprint backlog.

### Estimation
- Use `estimate: [Value]` (e.g., Fibonacci numbers 1, 2, 3, 5, 8) to size tasks.

### Auto-Pilot
- The `./scripts/tasks next` command uses an algorithm to determine the optimal next task based on:
    1.  Status (In Progress > Pending)
    2.  Dependencies (Unblocked > Blocked)
    3.  Sprint (Current Sprint > Backlog)
    4.  Priority (High > Low)
    5.  Type (Stories/Bugs > Tasks)

## Agent Integration

Agents (Claude, etc.) use this system to track their work.
- Always check `./scripts/tasks context` or use `./scripts/tasks next` before starting.
- Keep the task file updated with your progress.
- Use `review_requested` when you need human feedback.

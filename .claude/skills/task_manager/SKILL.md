---
name: Task Manager
description: Manage project tasks, status, and dependencies using the Task Documentation System.
---

# Task Manager Skill

This skill allows you to manage tasks within the project using the `scripts/tasks.py` utility. The Task Documentation System is the source of truth for all work.

## Core Commands

### List Tasks
List all tasks, optionally filtering by status or category.
`python3 scripts/tasks.py list [--status STATUS] [--category CATEGORY] [--sprint SPRINT]`

### Create Task
Create a new task documentation file.
`python3 scripts/tasks.py create [category] "Title" --desc "Description" --priority [high|medium|low] --type [task|story|bug]`
*   **Categories**: foundation, infrastructure, domain, presentation, migration, features, testing, review, security, research.

### Update Status
Update the status of a task.
`python3 scripts/tasks.py update [TASK_ID] [status]`
*   **Statuses**: pending, in_progress, review_requested, verified, completed, blocked.

### Show Task Details
Read the content of a task file.
`python3 scripts/tasks.py show [TASK_ID]`

### Context
Get the current working context (active tasks).
`python3 scripts/tasks.py context`

### Next Task
Get a recommendation for the next task to work on.
`python3 scripts/tasks.py next`

### Manage Dependencies
Link or unlink task dependencies.
`python3 scripts/tasks.py link [TASK_ID] [DEPENDENCY_ID]`
`python3 scripts/tasks.py unlink [TASK_ID] [DEPENDENCY_ID]`

## Usage Guidelines
1.  **Always** check `context` before starting work to see what is already in progress.
2.  **Always** create a task for new work using `create`.
3.  **Always** update the task status to `in_progress` when you start working on it.
4.  **Always** update the task status to `completed` when finished.

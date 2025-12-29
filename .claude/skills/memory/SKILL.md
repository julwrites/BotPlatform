---
name: Memory System
description: Store and retrieve long-term project context, architectural decisions, and learnings.
---

# Memory System Skill

This skill allows you to interact with the project's long-term memory system using `scripts/memory.py`. Use this to store information that should persist beyond the current session or task.

## Core Commands

### List Memories
List all stored memories, optionally filtering by tag.
`python3 scripts/memory.py list [--tag TAG] [--limit N]`

### Create Memory
Create a new memory entry. Use this for architectural decisions, important learnings, or project context.
`python3 scripts/memory.py create "Title" "Content" [--tags "tag1, tag2"]`

### Read Memory
Read the full content of a specific memory.
`python3 scripts/memory.py read [FILENAME_OR_SLUG]`

## Usage Guidelines
1.  **Check Memory** before starting complex tasks to see if there are relevant past learnings or decisions.
2.  **Create Memory** when you make a significant architectural decision or learn something that will be useful for future developers (human or AI).
3.  **Tags** help categorize memories (e.g., `architecture`, `setup`, `bugfix`).

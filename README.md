# obsidian-export

A utility written in Go to "export" Obsidian notes.

## Problem

I write my Obsidian notes as "living" documents with entries dated for each iteration
of the meeting. For example, consider the following two files:

**Standup.md**

```markdown
# Monday, May 5th, 2025

## Notes

These are my standup notes.
```

**Journal.md**

```markdown
# Monday, May 5th, 2025

## Notes

This is my journal.
```

## Solution

I like to collect together all of my notes for summarization, but having notes from
the day strewn across so many files is cumbersome. So I vibe-coded this Go utility
using Copilot Chat. This utility exports to a single file like so:

```markdown
# Standup.md

## Notes

These are my standup notes.

---

# Journal.md

## Notes

This is my journal.
```

This isn't specific to Obsidian; my approach to note taking can happen in any kind
of Markdown editor. I named this repository for Obsidian because that's what I use!

## CLI Flags

- `-help`/`--help` - Show the CLI usage
- `-vault-directory` - The Obsidian vault directory to use; defaults to `.`
- `-start-date` - Specify the start date of entries to include in the export
- `-end-date` - Specify the end date of entries to include in the export
- `-output-file` - Specify the output file path to use for export; defaults to `output.md`
- `-append` - Enable append mode on the output file; it truncates by default

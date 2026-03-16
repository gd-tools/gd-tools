# AGENTS.md
AI development context for the gd-tools project.

This file defines conventions and architectural rules that AI assistants
(ChatGPT, Copilot, Claude, etc.) should follow when generating code.

---

# Project Overview

Project: gd-tools

Purpose:
A Go-based infrastructure management system that installs and manages
server applications (WordPress, Nextcloud, Mail, RustDesk, etc.)
directly on the host system without containers.

Architecture separates:

- **Dev tool:** `gdt`
- **Production agent:** `gd-tools`

Both are built from the same Go module.

---

# Language and Communication

Code:
- Go

Comments:
- Must always be written in **English**

Communication with the project author:
- **German**

- CLI on dev: gdt (urfave/cli/v2)

Naming:
- WordPress prefix: wp
- Nextcloud prefix: nc
- RustDesk prefix: rd

---

# General Coding Rules

Prefer:

- simple, deterministic code
- minimal dependencies
- small functions
- explicit structures

Avoid:

- empty lines after function definition
- dynamic structures like map[string]any


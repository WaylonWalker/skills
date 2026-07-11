---
name: create-specification
description: Write a self-contained specification for a software project or feature. Target self-hosted web apps, static sites, CLIs, or mixed systems.
inspiration: https://raw.githubusercontent.com/github/awesome-copilot/refs/heads/main/skills/create-specification/SKILL.md
---

# Create Specification

Your goal is to create a new specification file for `${input:SpecPurpose}`.

Write a specification that a human or another agent can use without needing extra context. The spec must state the problem, the scope, the requirements, the system boundaries, and the acceptance criteria in clear Markdown.

This skill is for software projects in one or more of these categories:

- Self-hosted web apps
- Static sites
- Command-line tools
- Hybrid systems that combine the above

Do not assume a specific stack unless the repository already makes it clear. In particular, do not assume .NET, GitHub Actions, a managed cloud platform, or any named test framework unless the project already uses them.

## Workflow

1. Inspect the repository, existing docs, and naming conventions before writing.
2. Determine which project type applies: self-hosted web app, static site, CLI, or hybrid.
3. Identify the main user or operator goal, the system boundary, and the main constraints.
4. Ask a short clarifying question only if a missing detail would make the spec misleading.
5. Write one self-contained spec file in `/spec/`.

## Writing Rules

- Use precise, concrete language.
- Distinguish requirements from constraints, recommendations, and non-goals.
- Define acronyms and project-specific terms.
- Prefer stable interfaces and behaviors over implementation trivia.
- Include examples where they remove ambiguity.
- Include edge cases and failure modes.
- Keep the document self-contained. Do not require the reader to open other files to understand core requirements.
- If the repo already has relevant terminology or architecture, reuse it consistently.

## Project-Type Coverage

Include the sections that matter for the project type you identify.

### Self-Hosted Web Apps

Cover:

- User roles and primary flows
- Server-side and client-side responsibilities
- Authentication and authorization model
- Data storage and persistence needs
- Configuration, secrets, and environment variables
- Deployment shape, runtime assumptions, and operational constraints
- Logging, monitoring, backup, and recovery requirements when relevant

### Static Sites

Cover:

- Content structure and page types
- Build pipeline and asset generation
- Navigation, routing, and URL rules
- SEO, metadata, sitemap, and feed requirements when relevant
- Accessibility, responsiveness, and browser support expectations
- Hosting assumptions and cache or CDN behavior when relevant

### Command-Line Tools

Cover:

- Command structure and subcommands
- Flags, arguments, defaults, and configuration files
- Input and output behavior, including stdin, stdout, and stderr
- Exit codes and error handling
- Non-interactive versus interactive behavior
- Local filesystem, network, and credential access expectations
- Packaging, installation, and upgrade expectations when relevant

### Hybrid Systems

Cover:

- Boundaries between the web app, static site, and CLI parts
- Shared data contracts and configuration rules
- How users and operators move between components
- Which component owns each responsibility

## File Location And Naming

Save the specification in `/spec/`.

Use this filename pattern:

`spec-[purpose]-[topic].md`

Where `[purpose]` is one of:

- `architecture`
- `design`
- `process`
- `tool`
- `data`
- `infrastructure`
- `schema`

Use lowercase letters, numbers, and hyphens only.

## Required Template

Use this template. Remove bracketed guidance from the final file. Omit optional bullets that do not apply, but keep the section headings.

```md
---
title: [Short, specific title]
version: [Optional version or date tag]
date_created: [YYYY-MM-DD]
last_updated: [Optional YYYY-MM-DD]
owner: [Optional team or person]
tags: [Optional list such as architecture, cli, web-app, static-site, self-hosted]
---

# Summary

[One short paragraph that states what this spec covers and why it exists.]

## 1. Purpose & Scope

- **Purpose**: [What this spec defines]
- **In Scope**: [What is included]
- **Out of Scope**: [What is excluded]
- **Audience**: [Who should use this spec]

## 2. Product Context

- **Project Type**: [Self-hosted web app, static site, CLI, or hybrid]
- **Primary Users**: [End users, operators, developers, admins]
- **Primary Use Cases**: [Short list]
- **Operating Environment**: [Browser, server, terminal, local machine, container, VPS, on-prem, etc.]

## 3. Definitions

[Define acronyms, roles, domain terms, and internal shorthand used in this spec.]

## 4. Requirements

List specific, testable requirements.

- **REQ-001**: [Functional requirement]
- **REQ-002**: [Functional requirement]
- **SEC-001**: [Security or privacy requirement]
- **OPS-001**: [Operational requirement]
- **UX-001**: [User experience or interface requirement]

## 5. Constraints & Non-Goals

- **CON-001**: [Constraint such as runtime, hosting, compatibility, or legal limit]
- **NOG-001**: [Explicit non-goal]

## 6. Architecture / Design Notes

[Describe the proposed system shape at the level needed to implement correctly. Focus on components, boundaries, data flow, and responsibility split.]

### Components

| Component | Responsibility | Notes |
| --- | --- | --- |
| [Name] | [What it owns] | [Important detail] |

### Data Flow

[Describe the main request, command, build, or publishing flows.]

## 7. Interfaces & Data Contracts

[Document APIs, command syntax, file formats, events, or integration contracts. Use only the parts that apply.]

### Example API Contract

```json
{
  "example": true
}
```

### Example CLI Contract

```text
tool-name command --flag value
```

## 8. Configuration & Operations

[Document runtime configuration, secrets handling, deployment assumptions, observability, backup, restore, and maintenance expectations as needed.]

- **CFG-001**: [Configuration rule]
- **DEP-001**: [Deployment or hosting rule]
- **OBS-001**: [Logging or monitoring rule]

## 9. Edge Cases & Failure Modes

- **EDGE-001**: [Important edge case]
- **FAIL-001**: [Failure condition and required behavior]

## 10. Acceptance Criteria

Use testable statements.

- **AC-001**: Given [context], when [action], then [expected result]
- **AC-002**: Given [context], when [failure case], then [expected result]

## 11. Validation Strategy

[State how compliance with this spec will be checked. Mention the relevant test levels without assuming a specific framework.]

- **Unit Validation**: [What should be validated in isolation]
- **Integration Validation**: [What should be validated across boundaries]
- **End-to-End Validation**: [What user or operator flow should be verified]
- **Manual Checks**: [Any necessary manual review]

## 12. Dependencies & External Integrations

[List external systems, protocols, services, and platform assumptions that materially affect the design. Focus on required capabilities, not package names, unless a package is itself a hard constraint.]

- **EXT-001**: [External dependency and why it matters]

## 13. Examples

[Add short examples that remove ambiguity. Include representative input, output, page behavior, command usage, or deployment shape.]

## 14. Related Documents

- [Related spec or doc]
- [External standard or protocol reference]
```

## Completion Standard

Before you finish, verify that the spec:

- Matches the project type and does not include irrelevant sections from another type
- States clear requirements and explicit non-goals
- Defines interfaces and operational expectations at the right level
- Includes edge cases, failure behavior, and acceptance criteria
- Can be understood on its own by a new contributor or another agent

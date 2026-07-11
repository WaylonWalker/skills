---
name: web-design-guidelines
description: Review the UI code of self-hosted web apps for usability, accessibility, responsiveness, and implementation quality. Use this whenever the user asks to review a web UI, audit accessibility, check UX, inspect frontend code, or evaluate the design of a static site, server-rendered app, dashboard, admin panel, internal tool, settings page, data table, form flow, or CRUD interface. Use it even when the user does not say "design review" explicitly but is clearly asking whether an app screen, admin surface, or operational workflow is clear, usable, responsive, and production-ready.
metadata:
  author: opencode
  version: "1.1.0"
  argument-hint: <file-or-pattern>
---

# Web Design Guidelines

Review files for self-hosted web app UI quality.

## How It Works

1. Read the specified files, or ask the user which files or directories to review.
2. Infer the interface type: static marketing site, content site, server-rendered app, dashboard, admin surface, or CRUD workflow.
3. Review against the rules below, with extra attention to task flows, forms, tables, navigation, and empty or error states.
4. Output findings in terse `file:line` format grouped by file.

## Review Focus

Prioritize issues that affect real self-hosted applications:

- Broken or awkward CRUD flows
- Missing labels, feedback, loading, and error handling
- Layouts that fail on mobile, narrow laptop widths, or long content
- Tables, filters, and forms that are hard to scan or use
- Navigation or state that makes bookmarking, back/forward, or deep linking unreliable
- Risky interaction patterns for admin and destructive actions

## Rules

### Accessibility

- Icon-only buttons need an accessible name such as `aria-label`.
- Form controls need a visible label or a reliable accessible name.
- Use native interactive elements first: `button` for actions, `a` or framework links for navigation.
- Do not attach primary click behavior to `div` or `span` unless there is a strong reason and the full keyboard semantics are implemented.
- Images need meaningful `alt`, or `alt=""` when decorative.
- Decorative icons should be hidden from assistive tech.
- Validation, save states, and async updates should be announced when needed.
- Headings should form a sensible outline, especially in admin pages and settings screens.

### Layout and Responsiveness

- The page should remain usable on mobile, tablet, and common laptop widths.
- Avoid horizontal scrolling caused by fixed widths, unbounded tables, or overflowing code and content.
- Use layout primitives that adapt naturally before adding JS measurement.
- Long labels, IDs, emails, URLs, and filenames should wrap, truncate, or scroll intentionally.
- Sticky headers, sidebars, and in-page anchors should not obscure target content.
- Empty states should preserve layout clarity instead of collapsing into broken spacing.

### Navigation and Information Architecture

- Navigation should make location and next steps obvious.
- Links should behave like links, including open-in-new-tab and copy-link behavior.
- Important app state such as filters, tabs, pagination, selected records, and search should be shareable through the URL when reasonable.
- Back and forward navigation should feel natural in multi-step flows.
- Distinguish global navigation from local page actions.

### Forms and CRUD Workflows

- Inputs should use appropriate types, names, autocomplete, and input modes.
- Labels, help text, validation, and required state should be clear before submit.
- Do not block paste unless the user explicitly needs that protection.
- Submit actions should stay available until the request actually starts.
- Pending saves should show progress and prevent duplicate submissions.
- Validation errors should appear near the relevant field and be easy to recover from.
- Create, edit, and delete flows should make consequences clear.
- Destructive actions need confirmation, undo, or another safety mechanism.
- Bulk actions should clearly show selection scope and impact.

### Tables, Lists, and Dense Data

- Tables should remain readable with many columns, long content, or empty cells.
- Column headers should be specific and easy to scan.
- Sorting, filtering, selection, and row actions should be discoverable.
- Long lists should support search, filters, grouping, pagination, or virtualization when needed.
- Avoid rendering large datasets in ways that make the page sluggish.

### Feedback and States

- Every async action should communicate idle, loading, success, and failure states.
- Empty states should explain what happened and what the user can do next.
- Error states should include a recovery path, not just a problem statement.
- Disabled controls should have a visible reason when the reason is not obvious.
- Hover, active, focus, and selected states should be visually distinct.

### Typography and Copy

- Prefer short, direct labels that describe the exact action.
- Button text should be specific: `Save API Key` is better than `Continue`.
- Avoid vague placeholders as the only instruction.
- Use consistent heading hierarchy and spacing.
- Numeric or operational data should be easy to compare visually.
- Use plain language that matches software workflows, especially in admin and CRUD interfaces.

### Performance and Robustness

- Avoid patterns that cause heavy work on every keystroke or render.
- Large lists, tables, and dashboards should use pagination, chunking, or virtualization when appropriate.
- Avoid layout thrashing and unnecessary DOM measurement.
- Prefer resilient markup that still works if scripts load slowly.
- Hydration-sensitive UI should not render obviously different server and client output without intent.

### Self-Hosted App Concerns

- Settings, dashboards, and admin pages should favor clarity over decorative motion.
- Audit logs, statuses, environment names, and identifiers should be legible and copyable.
- Authentication, permissions, and dangerous operations should be reflected clearly in the UI.
- Deployment-specific details such as hostnames, ports, service names, and environment values should not break layout.
- System messages should help operators act quickly during failures.

### Common Anti-Patterns

- Clickable non-interactive elements standing in for buttons or links
- Missing labels on forms and icon buttons
- Actions hidden behind ambiguous kebab menus without primary affordances
- Tables that become unusable on smaller screens with no fallback
- Destructive actions placed next to common safe actions without separation
- Layouts that assume short English-only content
- Spinners or skeletons with no completion, error, or timeout handling
- Excessively generic copy such as `Submit`, `Confirm`, or `Update` when a concrete action is known

## Output Format

Group findings by file. Use `file:line` format so the paths are clickable.

- State the issue and why it matters in a few words.
- Focus on concrete findings, not general praise.
- If a file looks good, mark it with `✓ pass`.
- Prefer higher-severity issues first.

## Usage

When a user provides a file or pattern argument:

1. Read the specified files.
2. Infer the UI surface and primary user task.
3. Review against the rules above, emphasizing the issues most relevant to that surface.
4. Output findings in the format below.

If no files specified, ask the user which files to review.

```text
## src/users/UserTable.tsx

src/users/UserTable.tsx:48 - row actions hidden in overflow menu; common edit action too hard to discover
src/users/UserTable.tsx:91 - email column can force overflow on narrow screens

## src/settings/ProfileForm.tsx

src/settings/ProfileForm.tsx:33 - input missing label; placeholder should not carry field meaning
src/settings/ProfileForm.tsx:84 - save button has no pending state; duplicate submits likely

## src/dashboard/Overview.tsx

✓ pass
```

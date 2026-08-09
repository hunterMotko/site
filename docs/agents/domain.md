# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root, or
- **`CONTEXT-MAP.md`** at the repo root if it exists — it points at one `CONTEXT.md` per context. Read each one relevant to the topic.
- **`docs/adr/`** — read ADRs that touch the area you're about to work in. In multi-context repos, also check `src/<context>/docs/adr/` for context-scoped decisions.

If any of these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The `/domain-modeling` skill (reached via `/grill-with-docs` and `/improve-codebase-architecture`) creates them lazily when terms or decisions actually get resolved.

## `local/` is off-limits

`local/` is gitignored and holds private planning and career-strategy notes.
It is **not** domain documentation:

- Don't read it to answer domain questions, and don't treat anything in it as a project decision.
- Don't write skill output into it.
- Never reference it from a tracked file — it is absent from every clone, so the link would dangle for anyone else and in CI.

Domain material that should survive a clone belongs in `CONTEXT.md` or `docs/adr/`.

## File structure

This is a single-context repo:

```
/
├── CONTEXT.md
├── docs/adr/
│   ├── 0001-embed-templates-and-assets.md
│   └── 0002-pure-go-sqlite-for-static-binary.md
├── cmd/
└── internal/
```

A multi-context layout (root `CONTEXT-MAP.md` plus per-context `CONTEXT.md` and `src/<context>/docs/adr/`) is not in use here and shouldn't be introduced without a reason — this repo is one Go module serving one site.

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `CONTEXT.md`. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, that's a signal — either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/domain-modeling`).

## Flag ADR conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0007 (event-sourced orders) — but worth reopening because…_

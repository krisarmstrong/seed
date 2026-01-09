<!--
UPDATE WHEN:
- Adding new entry points or key files
- Introducing new patterns
- Discovering non-obvious behavior

Helps quickly navigate the codebase when resuming work.
-->

# Code Landmarks

Quick reference to important parts of the codebase.

## Entry Points
| Location | Purpose |
|----------|---------|
| cmd/seed/ | Main Go application entry |
| ui/src/main.tsx | Frontend entry point |

## Core Business Logic
| Location | Purpose |
|----------|---------|
| internal/api/ | HTTP handlers |
| internal/network/ | Network testing logic |
| internal/database/ | SQLite + migrations |

## Configuration
| Location | Purpose |
|----------|---------|
| internal/config/ | Go configuration |
| configs/ | Configuration files |
| ui/vite.config.ts | Frontend build config |

## Key Patterns
| Pattern | Example Location | Notes |
|---------|------------------|-------|
| Module structure | ui/src/ | Roots/Canopy/Shell/Sap/Harvest |
| API handlers | internal/api/ | REST conventions |

## Testing
| Location | Purpose |
|----------|---------|
| *_test.go | Go unit tests |
| ui/src/**/*.test.ts | Frontend tests (Vitest) |
| ui/e2e/ | Playwright E2E tests |

## Gotchas & Non-Obvious Behavior
| Location | Issue | Notes |
|----------|-------|-------|
| - | - | Add as discovered |

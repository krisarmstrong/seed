# Mustard Seed Networks - Comprehensive Rename Plan

## Overview

This document outlines the complete plan to rename the project from **NetScope/LuminetIQ** to **Mustard Seed Networks**.

## Brand Hierarchy

| Level | Name | Description |
|-------|------|-------------|
| **Company** | Mustard Seed Networks | The company/organization name |
| **Platform** | Mustard Seed Networks Platform | The overall product suite |
| **Core Product** | The Seed | The main application (binary) |
| **Modules** | Roots, Canopy, Shell, Sap, Harvest | Feature modules |

## Module Definitions

| Module | Function | Metaphor |
|--------|----------|----------|
| **Roots** | Path Analysis (traceroute, deep connectivity) | Foundation beneath the surface |
| **Canopy** | Wi-Fi Planning & Surveys | The expansive coverage above |
| **Shell** | Security Posture & Hardening | Protective outer layer |
| **Sap** | Live Telemetry & Monitoring | Lifeblood flowing through |
| **Harvest** | Reports, compliance, exports | The fruits of your labor |

## Naming Convention

| Context | Format | Example |
|---------|--------|---------|
| Company/Brand | Mustard Seed Networks | "Welcome to Mustard Seed Networks" |
| Platform | Mustard Seed Networks Platform | "The Mustard Seed Networks Platform" |
| Core Product | The Seed | "The Seed v1.0.0" |
| Binary/CLI | `seed` | `./seed --version` |
| Go Module | `github.com/krisarmstrong/seed` | import paths |
| Package Name | `seed` | Go package name |
| Storage Keys | `seed-*` | `seed-token`, `seed-theme` |
| Service Account | `seed` | systemd user |
| Config Directory | `~/.config/seed/` | user config |
| System Config | `/etc/seed/` | system config |

## Current → New Mapping

| Current | New |
|---------|-----|
| NetScope | Mustard Seed Networks (company) / The Seed (product) |
| LuminetIQ | Mustard Seed Networks (company) / The Seed (product) |
| `luminetiq` (binary) | `seed` |
| `netscope` (code refs) | `seed` |
| `netscope-*` (storage) | `seed-*` |

---

## Phase 1: Pre-Rename Preparation

### 1.1 Secure External Assets
- [ ] Register domain: `mustardseednetworks.com`
- [ ] Register domain: `mustardseednetworks.io`
- [ ] Create GitHub organization: `mustardseednetworks`
- [ ] Secure Twitter/X handle: `@mustardseednet` or `@theseednet`
- [ ] Secure LinkedIn company page: `/company/mustardseednetworks`
- [ ] Consult trademark attorney

### 1.2 Rename GitHub Repository

- [ ] Rename `krisarmstrong/netscope` to `krisarmstrong/seed`
- [ ] Update repository description
- [ ] Configure repository settings (branch protection, etc.)
- [ ] Update GitHub Actions secrets if needed

> **Note:** Business GitHub organization (`mustardseednetworks`) will be created later.
> Initial development continues on personal account.

### 1.3 Backup Current State
- [ ] Tag current release: `pre-rename-backup`
- [ ] Export all issues and PRs
- [ ] Document current CI/CD configuration

---

## Phase 2: Go Module Rename

### 2.1 Update go.mod

**File:** `go.mod`

```diff
- module github.com/krisarmstrong/luminetiq
+ module github.com/krisarmstrong/seed
```

### 2.2 Update All Go Imports (30+ files, ~150 occurrences)

**Pattern to replace:**

```text
github.com/krisarmstrong/luminetiq/internal/*
→
github.com/krisarmstrong/seed/internal/*
```

**Files affected:**

- `cmd/luminetiq/main.go` → `cmd/seed/main.go`
- `internal/api/*.go` (15+ files)
- `internal/auth/*.go`
- `internal/cable/*.go`
- `internal/config/*.go`
- `internal/dhcp/*.go`
- `internal/discovery/*.go`
- `internal/dns/*.go`
- `internal/gateway/*.go`
- `internal/iperf/*.go`
- `internal/logging/*.go`
- `internal/network/*.go`
- `internal/publicip/*.go`
- `internal/snmp/*.go`
- `internal/speedtest/*.go`
- `internal/survey/*.go`
- `internal/system/*.go`
- `internal/validation/*.go`
- `internal/version/*.go`
- `internal/vlan/*.go`
- `internal/wifi/*.go`
- All `*_test.go` files

### 2.3 Update .golangci.yml

```diff
- - github.com/krisarmstrong/luminetiq
+ - github.com/krisarmstrong/seed
```

---

## Phase 3: Directory Rename

### 3.1 Rename Command Directory

```bash
git mv cmd/luminetiq cmd/seed
```

### 3.2 Rename Config File

```bash
git mv configs/netscope.yaml configs/seed.yaml
```

### 3.3 Rename Packaging Files

```bash
git mv packaging/luminetiq.spec packaging/seed.spec
git mv packaging/luminetiq.service packaging/seed.service
git mv deploy/systemd/luminetiq.service deploy/systemd/seed.service
git mv deploy/luminetiq-dev.service deploy/seed-dev.service
```

### 3.4 Rename Web Assets

```bash
git mv web/public/luminetiq.svg web/public/seed.svg
```

---

## Phase 4: Build System Updates

### 4.1 Makefile Changes

| Line | Current | New |
|------|---------|-----|
| 43 | `BINARY_NAME=luminetiq` | `BINARY_NAME=seed` |
| 61-63 | LDFLAGS with luminetiq path | Update module path to `krisarmstrong/seed` |
| All | `./cmd/luminetiq` | `./cmd/seed` |
| All | Deploy paths with luminetiq | Deploy paths with seed |

### 4.2 Dockerfile

```diff
- go build ... -o /luminetiq ./cmd/luminetiq
+ go build ... -o /seed ./cmd/seed
```

---

## Phase 5: Frontend Updates

### 5.1 Storage Keys (web/src/hooks/useAuth.ts)

```diff
- "netscope-token"
- "netscope-token-expiry"
- "netscope-username"
+ "seed-token"
+ "seed-token-expiry"
+ "seed-username"
```

### 5.2 Theme Storage (web/src/hooks/useTheme.ts)

```diff
- const STORAGE_KEY = "netscope-theme"
+ const STORAGE_KEY = "seed-theme"
```

### 5.3 HTML Updates (web/index.html)

```diff
- <link rel="icon" type="image/svg+xml" href="/luminetiq.svg" />
- <title>LuminetIQ</title>
- <meta name="description" content="LuminetIQ - Illuminate Your Network" />
+ <link rel="icon" type="image/svg+xml" href="/seed.svg" />
+ <title>The Seed | Mustard Seed Networks</title>
+ <meta name="description" content="The Seed - Network diagnostic platform by Mustard Seed Networks" />
```

### 5.4 SVG Label (web/public/seed.svg)

```diff
- aria-label="NetScope"
+ aria-label="The Seed"
```

### 5.5 Component Updates
- `web/src/App.tsx` - Update comments
- `web/src/main.tsx` - Update comments
- `web/src/components/settings/SettingsDrawer.tsx` - Update version display
- `web/src/lib/api.ts` - Update comments
- `web/src/hooks/useWebSocket.ts` - Update comments

---

## Phase 6: Backend Code Updates

### 6.1 Display Text (cmd/seed/main.go)

- Update ASCII art banner to "THE SEED" or "MUSTARD SEED NETWORKS"
- Update version display: `"The Seed %s (Mustard Seed Networks)\n"`
- Update setup wizard text
- Update credential display

### 6.2 ICMP/Network Payloads

**internal/discovery/icmp.go:**

```diff
- data := []byte("netscope")
+ data := []byte("seed")
```

**internal/discovery/traceroute.go:**

```diff
- data := []byte("NETSCOPE")
+ data := []byte("SEED")
```

### 6.3 Config Paths (internal/discovery/oui.go)

```diff
- filepath.Join(os.Getenv("HOME"), ".config", "netscope", "oui.txt")
+ filepath.Join(os.Getenv("HOME"), ".config", "seed", "oui.txt")
```

### 6.4 Export Filenames (internal/api/handlers_status.go)

```diff
- "netscope-export.json"
+ "seed-export.json"
```

### 6.5 Default Password (internal/auth/auth.go)

Consider updating default password from "netscope" to "seed"

---

## Phase 7: Service & Deployment Updates

### 7.1 Systemd Service Files

**deploy/systemd/seed.service:**

```diff
- Description=LuminetIQ Network Diagnostic Tool
- ExecStart=/usr/local/bin/luminetiq
- User=luminetiq
- Group=luminetiq
+ Description=The Seed - Mustard Seed Networks Platform
+ ExecStart=/usr/local/bin/seed
+ User=seed
+ Group=seed
```

### 7.2 Install/Uninstall Scripts

- `deploy/systemd/install.sh` - Update all references
- `deploy/systemd/uninstall.sh` - Update all references
- `scripts/deploy-systemd.sh` - Update all references

### 7.3 Packaging Files

**packaging/control:**

```diff
- Package: netscope
- Description: NetScope is...
+ Package: seed
+ Description: The Seed is the core product of Mustard Seed Networks...
```

**packaging/seed.spec:**

- Update Name, URL, description
- Update all file paths

---

## Phase 8: CI/CD Updates

### 8.1 GitHub Workflows

**.github/workflows/ci.yml:**

- Update comment header to "The Seed - Mustard Seed Networks"
- Update binary artifact names: `seed-linux-amd64`, `seed-darwin-amd64`, etc.

**.github/workflows/release-please.yml:**

- Update all binary references
- Update artifact names to `seed-*`

**.github/workflows/release.yml:**

- Update binary names
- Update download URLs

**.github/workflows/docker-publish.yml:**

- Update comment header
- Update image names if applicable

### 8.2 Issue Templates

- `.github/ISSUE_TEMPLATE/hardware-report.yml` - Update branding to "The Seed"
- `.github/ISSUE_TEMPLATE/bug_report.md` - Update version label to "The Seed Version"

---

## Phase 9: Documentation Updates

### 9.1 Root Documentation
| File | Changes Required |
|------|------------------|
| `README.md` | Title, description, installation commands, examples |
| `CONTRIBUTING.md` | All NetScope/LuminetIQ references |
| `SECURITY.md` | Product name references |
| `LICENSE` | "Licensed Work: Mustard Seed Networks" |
| `STYLE_GUIDE.md` | Title, import examples |
| `PROJECT_PLAN.md` | Project name references |
| `HARDWARE.md` | Product name references |

### 9.2 Docs Directory
- `docs/VERIFICATION.md` - All LuminetIQ references
- All wiki files in `docs/wiki/` (25+ files)

### 9.3 Web Documentation
- `web/docs/ui-style.md` - Title update
- `web/THEMING.md` - Check for references

### 9.4 Claude Configuration

Claude Code reads CLAUDE.md from multiple locations (all are merged):

1. `~/.claude/CLAUDE.md` - User-level global (not project-specific)
2. `PROJECT_ROOT/CLAUDE.md` - Project root
3. `PROJECT_ROOT/.claude/CLAUDE.md` - Project hidden directory

#### File: `.claude/CLAUDE.md`

This file requires extensive updates:

| Line | Current | New |
|------|---------|-----|
| 1 | `# Claude Operating Guidelines for LuminetIQ` | `# Claude Operating Guidelines for The Seed` |
| 35 | `cd ~/luminetiq && git pull` | `cd ~/seed && git pull` |
| 63 | `/home/krisarmstrong/luminetiq` | `/home/krisarmstrong/seed` |
| 64 | `./luminetiq` | `./seed` |
| 75-76 | `systemctl ... luminetiq` | `systemctl ... seed` |
| 76 | `journalctl -u luminetiq` | `journalctl -u seed` |
| 78-81 | `luminetiq credentials` | `seed credentials` |
| 104 | `luminetiq/` | `seed/` |
| 105 | `cmd/luminetiq/` | `cmd/seed/` |

**Full rewrite recommended** to reflect:

- Brand: "The Seed" (product) by "Mustard Seed Networks" (company)
- Binary: `seed`
- Service: `seed.service`
- Paths: `/home/krisarmstrong/seed`, `cmd/seed/`
- Module metaphors: Roots, Canopy, Shell, Sap, Harvest

#### File: `.claude/settings.local.json`

This file contains permission allowlist entries with old paths/names:

| Line(s) | Current | New |
|---------|---------|-----|
| 7 | `netscope` in SSH path | `seed` |
| 23 | `netscope` in deploy cmd | `seed` |
| 34, 56-57 | `/netscope` paths | `/seed` |
| 59-60 | `netscope` binary refs | `seed` |
| 70 | `./luminetiq:*` | `./seed:*` |
| 92-99, 124-125 | `/netscope` in git paths | `/seed` |

**After rename:** These permissions may need regenerating as paths change.
Consider clearing and re-allowing as you use the new paths.

---

## Phase 10: Test Updates

### 10.1 Go Tests
- `internal/auth/auth_test.go` - Update test passwords
- All test files with import path updates

### 10.2 Frontend Tests
- `web/src/App.test.tsx` - Update storage key references
- `web/src/hooks/useAuth.test.ts` - Update all netscope references

---

## Phase 11: GitHub Repository Rename

### 11.1 Rename Repository

```bash
# Rename repository via GitHub UI or CLI
gh repo rename seed --repo krisarmstrong/netscope

# Update local remote
git remote set-url origin git@github.com:krisarmstrong/seed.git
```

### 11.2 Update All References

- Update any hardcoded GitHub URLs in code/docs
- Update CI/CD badge URLs in README
- Update issue/PR templates with new repo URL

### 11.3 Future: Business Organization Migration

When ready to move to business account:

```bash
# Transfer repository to organization
# Done via GitHub UI: Settings > Transfer ownership

# Update remote after transfer
git remote set-url origin git@github.com:mustardseednetworks/seed.git
```

---

## Phase 12: Verification

### 12.1 Build Verification

```bash
make clean
make all
./seed version
```

### 12.2 Search for Remaining References

```bash
# Check for any remaining old names (should return no results)
grep -ri "netscope" --include="*.go" --include="*.ts" --include="*.tsx" \
  --include="*.md" --include="*.yaml" --include="*.yml" --include="*.json" .

grep -ri "luminetiq" --include="*.go" --include="*.ts" --include="*.tsx" \
  --include="*.md" --include="*.yaml" --include="*.yml" --include="*.json" .
```

### 12.3 Test All Features
- [ ] Build passes
- [ ] All tests pass
- [ ] Frontend builds
- [ ] Docker image builds
- [ ] Service starts correctly
- [ ] Web UI loads
- [ ] Authentication works
- [ ] All API endpoints respond

---

## Execution Order

1. **Phase 1** - Secure external assets (domains, GitHub org, social)
2. **Phase 2** - Go module rename (most critical)
3. **Phase 3** - Directory renames
4. **Phase 4** - Build system updates
5. **Phase 5** - Frontend updates
6. **Phase 6** - Backend code updates
7. **Phase 7** - Service & deployment updates
8. **Phase 8** - CI/CD updates
9. **Phase 9** - Documentation updates
10. **Phase 10** - Test updates
11. **Phase 11** - GitHub repository migration
12. **Phase 12** - Verification

---

## Rollback Plan

If issues arise:
1. Revert to `pre-rename-backup` tag
2. Keep old repository active
3. Document lessons learned

---

## Post-Rename Tasks

- [ ] Update any external links
- [ ] Notify users of name change
- [ ] Update package registries
- [ ] Update Docker Hub/container registries
- [ ] Monitor for broken links/references
- [ ] Close issues #501, #504, #505

---

## Detailed File Audit

### Files Containing "netscope" (39 unique files)

| Category | File | Reference Type |
|----------|------|----------------|
| **Auth** | `internal/auth/auth.go` | Default password "netscope", storage key comments |
| **Auth** | `internal/auth/auth_test.go` | Test password references |
| **Discovery** | `internal/discovery/icmp.go` | ICMP payload `[]byte("netscope")` |
| **Discovery** | `internal/discovery/traceroute.go` | Traceroute payload `[]byte("NETSCOPE")` |
| **Discovery** | `internal/discovery/oui.go` | Config path `~/.config/netscope/` |
| **API** | `internal/api/handlers_status.go` | Export filename `netscope-export.json` |
| **Config** | `internal/config/config.go` | Example domain references |
| **Config** | `configs/netscope.yaml` | **FILE RENAME REQUIRED** |
| **Frontend** | `web/src/hooks/useAuth.ts` | Storage keys `netscope-token`, etc. |
| **Frontend** | `web/src/hooks/useAuth.test.ts` | Test storage key references |
| **Frontend** | `web/src/hooks/useTheme.ts` | Storage key `netscope-theme` |
| **Frontend** | `web/src/types/settings.ts` | Storage key references |
| **Frontend** | `web/src/App.test.tsx` | Test storage keys |
| **Frontend** | `web/src/index.css` | CSS comment "NetScope Theme" |
| **Assets** | `web/public/luminetiq.svg` | aria-label="NetScope" |
| **Build** | `Dockerfile` | Build output path references |
| **Packaging** | `packaging/luminetiq.spec` | GitHub URL |
| **Packaging** | `packaging/control` | Package name "netscope" |
| **Scripts** | `scripts/build-iperf3.sh` | NetScope references |
| **Docs** | `README.md` | Product name references |
| **Docs** | `CONTRIBUTING.md` | NetScope references |
| **Docs** | `SECURITY.md` | NetScope references |
| **Docs** | `PROJECT_PLAN.md` | NetScope references |
| **Docs** | `LICENSE` | "Licensed Work" reference |
| **Docs** | `web/docs/ui-style.md` | NetScope references |
| **Issues** | `.github/ISSUE_TEMPLATE/bug_report.md` | "NetScope Version" label |
| **Types** | `web/src/types/index.ts` | NetScope type references |
| **Components** | `web/src/components/cards/SystemHealthCard.tsx` | NetScope references |
| **Components** | `web/src/components/settings/SettingsDrawer.tsx` | NetScope references |

### Files Containing "luminetiq" (99 unique files)

| Category | Count | Examples |
|----------|-------|----------|
| **Go Imports** | 48 files | All `internal/*/*.go` files with `github.com/krisarmstrong/luminetiq` |
| **Makefile** | 1 file | BINARY_NAME, LDFLAGS, all paths |
| **Dockerfile** | 1 file | Build output, entry point |
| **Web Config** | 8 files | vite.config.ts, eslint.config.mjs, package.json, etc. |
| **Service Files** | 5 files | systemd services, install/uninstall scripts |
| **Workflows** | 4 files | CI, release, docker-publish |
| **Documentation** | 25+ files | Wiki, templates, README, guides |
| **Frontend** | 15+ files | Components, hooks, lib, styles |

### Directories Requiring Rename

| Current Path | New Path |
|--------------|----------|
| `cmd/luminetiq/` | `cmd/seed/` |
| `configs/netscope.yaml` | `configs/seed.yaml` |
| `web/public/luminetiq.svg` | `web/public/seed.svg` |
| `packaging/luminetiq.spec` | `packaging/seed.spec` |
| `packaging/luminetiq.service` | `packaging/seed.service` |
| `deploy/systemd/luminetiq.service` | `deploy/systemd/seed.service` |
| `deploy/luminetiq-dev.service` | `deploy/seed-dev.service` |

### Special Cases

1. **CSS Theme Comment** (`web/src/index.css` line 48):

   ```css
   /* NetScope Theme - WiFi Vigilante Color Scheme */
   ```

   → Change to: `/* The Seed Theme - Mustard Seed Networks */`

2. **HTML Meta/Title** (`web/index.html`):
   - favicon: `/luminetiq.svg` → `/seed.svg`
   - apple-touch-icon: `/luminetiq.svg` → `/seed.svg`
   - description: "LuminetIQ - Illuminate Your Network" → "The Seed - Network diagnostics by Mustard Seed Networks"
   - title: "LuminetIQ" → "The Seed | Mustard Seed Networks"

3. **SVG aria-label** (`web/public/luminetiq.svg`):

   ```svg
   aria-label="NetScope"
   ```

   → Change to: `aria-label="The Seed"`

4. **Makefile grep pattern** (line 216):

   ```makefile
   ps aux | grep '[l]uminetiq'
   ```

   → Change to: `ps aux | grep '[s]eed'`

5. **Go module path in 48 files**:

   ```go
   github.com/krisarmstrong/luminetiq/internal/*
   ```

   → Change to: `github.com/krisarmstrong/seed/internal/*`

---

## Files Summary

### Total Files to Modify: ~100+

| Category | Count |
|----------|-------|
| Go source files | 35+ |
| Go test files | 15+ |
| TypeScript/React files | 20+ |
| Documentation (MD) | 30+ |
| Configuration files | 10+ |
| Workflow files | 5 |
| Shell scripts | 8 |
| Service files | 4 |
| HTML files | 2 |
| Other | 5+ |

---

## Related Issues

- #501 - Implement New Naming Convention
- #504 - Comprehensive Codebase & Asset Renaming
- #505 - Update Documentation & Readmes

---

*Plan created: December 15, 2025*
*Target completion: TBD*

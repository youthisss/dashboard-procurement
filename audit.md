# Dashboard Procurement Repo Audit

Last updated: 2026-04-30

## Scope

This audit covers the current `dashboard_procurement` working tree.

Excluded from review:

- `wipe_data.go`
- Local secrets such as `.env`
- Generated or ignored artifacts: `frontend/node_modules`, `frontend/.next`, `.gocache`, uploads, Excel files, and local build outputs

The repository is currently dirty with many expected implementation changes from prior batches. This audit treats the current working tree as the source of truth and does not assume a clean Git baseline.

## Current Verdict

Status: Security hardening pass completed; release verification is still not complete.

The application has a stronger baseline than the earlier audit:

- JWT secret validation exists at backend startup.
- Auth token storage uses an HttpOnly cookie, not browser localStorage.
- Frontend API calls use the shared credential-aware API client.
- Upload filenames are sanitized.
- Import preview row editing, Actual Contract disabled state, and edited-sheet import confirmation were fixed.
- Frontend dependency advisory check currently passes.

Production deployment should still wait until the remaining release gates below are verified outside the current local blockers.

## Latest Validation

Passed on 2026-04-30:

- Backend tests: `go test ./...`
- Backend vet: `go vet ./...`
- Frontend typecheck: `cmd /c npx tsc --noEmit --pretty false`
- Frontend lint: `cmd /c npm run lint`
- Frontend dependency audit: `cmd /c npm audit --audit-level=moderate`
  - Result: `found 0 vulnerabilities`
- `git diff --check -- audit.md`

Validation limits still present:

- `cmd /c npm run build` compiles successfully, then fails at Next.js TypeScript worker startup with Windows `spawn EPERM`.
- Backend Docker build is not verified because Docker Desktop/Linux engine is not running.
- Browser/manual flows are not verified in this audit pass.

Current deprecated scan result:

- No `console.log`, `debugger`, `React.FC`, or CommonJS `require(...)` problems were found in source.
- Compatibility cleanup lines still intentionally delete old `pagination.position` before rendering AntD tables.
- No live deprecated AntD `Space orientation` prop remains in source.

## Open Security Findings

No open security findings remain from this execution order after the 2026-04-30 implementation pass.

## Open Quality And Release Findings

### Q2: Production Build Is Blocked By Local Windows `spawn EPERM`

Severity: Medium

Status: Open

Evidence:

- `cmd /c npm run build` compiled successfully with Next.js 16.2.4.
- The build failed afterward during TypeScript worker startup with `Error: spawn EPERM`.

Required fix:

- Re-run production build in a clean terminal or CI environment.
- If it still fails outside this machine, investigate Next.js worker spawning, antivirus/process policy, or Node install permissions.

### Q3: Backend Docker Build Is Not Verified

Severity: Medium

Status: Open

Evidence:

- `backend/Dockerfile:2` now uses `golang:1.25-alpine`.
- Docker build could not be verified because Docker Desktop/Linux engine was not running.

Required fix:

- Start Docker Desktop.
- Run `docker build -t dashboard-procurement-backend-audit .` from `backend/`.

### Q4: Import Browser Flows Are Not Manually Verified

Severity: Low

Status: Open

Evidence:

- Automated tests and typecheck pass.
- Browser upload/edit/confirm workflows were not run during this audit refresh.

Required fix:

- Test import preview editing across page 2 or later with a multi-row Excel file.
- Test confirm import after the uploaded temp file is missing but edited `sheets` payload is present.

## Fixed Findings

### F1: Frontend Dependency Advisory Check

Status: Fixed on 2026-04-30

Evidence:

- `cmd /c npm audit --audit-level=moderate` completed after explicit approval to contact the npm registry.
- Result: `found 0 vulnerabilities`.

### F2: Backend Dockerfile Go Version Mismatch

Status: Fixed in code on 2026-04-30

Evidence:

- `backend/go.mod:3` declares `go 1.25.0`.
- `backend/Dockerfile:2` now uses `golang:1.25-alpine`.

### F3: Import Preview Row Editing Across Pagination

Status: Fixed on 2026-04-30

Evidence:

- `frontend/src/app/import/page.tsx` now creates preview-only stable row indexes.
- Editable cell updates and `rowKey` use the stable preview row index instead of AntD rendered row index.
- Frontend typecheck and lint pass.

### F4: Actual Contract Import Mode Before Backend Support

Status: Fixed on 2026-04-30

Evidence:

- The import page renders `ContractModeSwitch` without `actualEnabled`.
- Actual Contract remains visible but disabled until backend/database support exists.

### F5: Edited Import Confirm Required Temporary File

Status: Fixed on 2026-04-30

Evidence:

- `ConfirmImport` only requires the uploaded temp file when no parsed `sheets` payload is supplied.
- Edited preview confirmation can proceed through `ConfirmParsedSheets`.
- Cleanup ignores already-missing temp files after successful import.

### F6: Previous Hardening Batch

Status: Fixed before this refresh

Summary:

- Upload filename sanitization and tests.
- Admin gating now uses backend-provided `is_admin`, not hard-coded frontend username checks.
- Contract mutations and audit inserts are transaction-scoped.
- Import logs no longer expose row-level procurement data.
- Import-created vendor/mill codes use stable generated short hashes.
- AntD pagination placement warnings were addressed.

### F7: Login Brute-Force Protection

Status: Fixed on 2026-04-30

Evidence:

- `backend/internal/handlers/login_limiter.go` tracks failed login attempts by IP plus username.
- `backend/internal/handlers/auth_handler.go` blocks locked login keys with generic failure responses and resets attempts after successful login.
- `backend/internal/handlers/login_limiter_test.go` covers lockout and reset behavior.

### F8: Request Body Limits

Status: Fixed on 2026-04-30

Evidence:

- `backend/internal/middleware/security.go` adds `RequestBodyLimit`.
- `backend/internal/router/router.go` applies a global API body limit plus route-specific import upload and import-confirm limits.
- `backend/internal/handlers/contract_handler.go` replaces unbounded contract update reads with a bounded reader.
- `backend/internal/middleware/security_test.go` covers oversized body rejection.

### F9: Production Server Runtime Hardening

Status: Fixed on 2026-04-30

Evidence:

- `backend/cmd/server/main.go` now starts an explicit `http.Server` with `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`.
- `backend/internal/router/router.go` sets Gin mode from config and configures trusted proxies explicitly.
- `backend/internal/config/config.go` adds `GIN_MODE` and `TRUSTED_PROXIES`.

### F10: Strict Origin Protection For Cookie Mutations

Status: Fixed on 2026-04-30

Evidence:

- `backend/internal/middleware/security.go` adds strict `Origin`/`Referer` validation for unsafe methods.
- `backend/internal/router/router.go` applies it to authenticated routes.
- `backend/internal/middleware/security_test.go` covers allowed and blocked origins.

### F11: Admin-Only Mutation Routes

Status: Fixed on 2026-04-30

Evidence:

- `backend/internal/middleware/security.go` adds reusable admin-only middleware.
- `backend/internal/router/router.go` applies admin-only protection to create, update, delete, agreement, import, and user-management mutations while leaving read routes authenticated.

### F12: Shared Password Policy

Status: Fixed on 2026-04-30

Evidence:

- `backend/internal/services/password_policy.go` centralizes password validation.
- Admin bootstrap/sync, user creation, and password update now use the shared policy.
- `backend/internal/services/password_policy_test.go` covers weak/default rejection and a strong password case.

### F13: Credentialed Cookie Deployment Documentation

Status: Fixed on 2026-04-30

Evidence:

- `DEPLOYMENT.md` now documents exact `FRONTEND_ORIGINS` for production cookie auth.
- Wildcard Vercel origin guidance was removed from production instructions and preview guidance is explicit allowlist only.

### F14: Deprecated AntD `Space orientation`

Status: Fixed on 2026-04-30

Evidence:

- `frontend/src/app/login/page.tsx` now uses `direction="vertical"`.
- Frontend typecheck and lint pass.

## Completed Execution Order

Completed on 2026-04-30 in the requested order:

1. S1 fixed: login brute-force protection.
2. S6 fixed: request body size limits.
3. S5 fixed: server runtime hardening.
4. S2 fixed: strict origin protection for cookie-authenticated mutation routes.
5. S3 fixed: admin-only mutation authorization.
6. S4 fixed: shared password policy.
7. S7 fixed: deployment documentation for credentialed cookie auth.
8. Q1 fixed: deprecated AntD prop cleanup.

Remaining release verification:

1. Re-run `cmd /c npm run build` in a clean terminal or CI environment because this Windows session still fails after successful compile with `spawn EPERM`.
2. Start Docker Desktop/Linux engine and run `docker build -t dashboard-procurement-backend-audit .` from `backend/`.
3. Manually test import browser flows with real Excel files.

# Changelog

This tracks notable changes to the Vantage-specific subsystem (scanner engine, dashboard,
deployment) on top of upstream Gophish. It replaces the earlier internal audit/implementation
reports, which had gone stale and contradicted each other.

## Unreleased

### Known open items
- `scanner/` has no automated test coverage (`ToolRegistry`, `BuildArgs`, orchestration error
  paths are all untested today).
- The per-tool parallel goroutine in `RunTask` (`scanner/engine.go`) does not have its own
  `recover()`; a panic inside one parallel tool run can crash the whole process even though the
  outer scan goroutines are guarded.
- Single concurrent scan per target (lock-based) — no scan queue yet.
- SQLite is single-writer; fine for typical team usage, not for high-concurrency multi-tenant
  deployments.

### Resolved since the original April 2026 internal audit
- Panic recovery: scan goroutines (`RunDiscovery`, `RunTask` top-level, `ScannerLogHub`) now
  wrap their work in `defer recover()`.
- Dashboard now renders through Go's `html/template` (`vantage_dashboard.html` via `base`
  layout) with a CSRF token, instead of being served as static HTML.
- Interactsh and Cloudlist now have dedicated tool adapters (`scanner/tools.go`).
- Network interface selection has a dedicated `UserNetworkConfig` model and API
  (`models/vantage.go`, `controllers/api/network.go`).

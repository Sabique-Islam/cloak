
---

# Security

This project creates controlled, ethical experiments with harmless artifacts. If you discover a security issue, please report it responsibly.

---

## Scope

- Cloak generates fake signals for OSINT research - all artifacts are harmless and use lab-controlled accounts.
- Security concerns include: API key exposure, database access, unintended artifact propagation, or web UI vulnerabilities.
- This is a research project; we prioritize learning and responsible disclosure.

---

## Reporting a Vulnerability

- Click on the Security tab.
- Click "Report a vulnerability".
- Include : Issue summary, Impact, PoC, Affected Components.

---

## Response process

- We will acknowledge reports and investigate as time permits.
- Fixes will be coordinated with reporters.
- Workarounds or patches will be provided depending on severity and project timeline.

---

## What we protect

- API keys (GitHub, Cloudflare, OSINT providers) stored in config/env files.
- SQLite database containing artifacts, detections, and audit logs.
- Web UI endpoints and user inputs.
- Lab domain and repository integrity.

---


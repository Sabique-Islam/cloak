
---

# Timeline

| Week |           Focus | Core Tasks                                                                                                          | Deliverables                                                                              |
| ---- | --------------: | ------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| 1    |    Prep & Infra | Team intro, project planning, stack discussion, define scope, setup GitHub & Cloudflare API keys, DB schema | Project charter, team roles, config samples, initial DB schema                            |
| 2    |    Artifact Gen | GitHub commit spoofing & repo creation, Cloudflare DNS injection & teardown, unit tests                             | `github-generator`, `dns-injector`, sample artifacts in SQLite DB, artifact creation docs |
| 3    | OSINT Ingestion | Build adapters for GitHub Search, VirusTotal, GreyNoise; parallel queries, log ingestion                            | OSINT adapters, `detections` table, logging & audit, prototype web UI                     |
| 4    |        Analysis | Compute metrics (visibility, latency, persistence), generate charts, refine experiments                             | Analysis scripts, charts & visualizations, updated UI reflecting results                  |
| 5    |      Validation | Re-run experiments, robustness testing, finalize dataset, execute teardown, confirm ethics                          | Final redacted dataset, final report, audit log confirmation  |

---

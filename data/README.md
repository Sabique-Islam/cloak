# Data Directory

## Structure

- `results.db` : SQLite db (created at runtime)
- `exports/` : Exported metrics and results (CSV/JSON)
- `logs/` : Raw provider responses and system logs

## Note

`results.db` is created automatically when application starts.
Export files are generated using `scripts/export.go`.

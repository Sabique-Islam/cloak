# IdeationDocs Test Suite

Comprehensive tests for the IdeationDocs directory changes.

## Setup
```bash
pip install -r requirements.txt
```

## Run Tests
```bash
# All tests
pytest

# With coverage
pytest --cov=. --cov-report=html

# Specific test file
pytest tests/ideation_docs/test_documentation.py
```

## Coverage

- **Documentation Tests**: Validation of PoC.md, SETUP.md, and notes.txt.
- **Database Tests**: Schema structure, data integrity, and relationships.
- **Integration Tests**: Cross-file consistency and binary file validation.
- **Configuration Tests**: .gitignore changes

Tests validate Markdown syntax, SQL schemas, database structure, data consistency,
and cross-references between documentation and database.
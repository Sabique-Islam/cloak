# Quick Start Guide - IdeationDocs Test Suite

## 🚀 Get Started in 3 Steps

### 1. Install Dependencies
```bash
cd /home/jailuser/git
pip install -r tests/requirements.txt
```

### 2. Run Tests
```bash
pytest tests/
```

### 3. View Results
All tests should pass ✓

## 📊 What's Being Tested?

### Documentation Files
- ✓ PoC.md structure and SQL schemas
- ✓ SETUP.md step-by-step guide
- ✓ notes.txt content and consistency

### Database
- ✓ cloak.db schema structure
- ✓ Data integrity and relationships
- ✓ Test data validation

### Binary Files
- ✓ PDF file format
- ✓ PNG

### Integration
- ✓ Cross-file consistency
- ✓ Schema documentation accuracy
- ✓ Reference resolution

## 🎯 Common Commands

```bash
# Run all tests
pytest tests/

# Run with verbose output
pytest tests/ -v

# Run specific test file
pytest tests/ideation_docs/test_documentation.py

# Run specific test class
pytest tests/ideation_docs/test_documentation.py::TestPoCMarkdown

# Run with coverage report
pytest --cov=IdeationDocs --cov-report=term-missing tests/

# Run and stop at first failure
pytest tests/ -x

# Run tests matching pattern
pytest tests/ -k "database"
```

## 📝 Test Organization
```text
```
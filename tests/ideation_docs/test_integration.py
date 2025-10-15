"""Integration tests across IdeationDocs files."""
import re
import sqlite3
from pathlib import Path

DOCS_PATH = Path(__file__).parent.parent.parent / "IdeationDocs"

class TestCrossFileConsistency:
    """Test consistency across multiple files."""
    
    def test_all_expected_files_present(self):
        expected = ["PoC.md", "SETUP.md", "notes.txt", "cloak.db",
                   "DNScheck.png", "Cloak-Adversarial_OSINT_Evasion.pdf"]
        for filename in expected:
            assert (DOCS_PATH / filename).exists(), f"Missing {filename}"
    
    def test_poc_schema_matches_database(self):
        poc = (DOCS_PATH / "PoC.md").read_text()
        
        conn = sqlite3.connect(DOCS_PATH / "cloak.db")
        cursor = conn.cursor()
        cursor.execute("PRAGMA table_info(artifacts)")
        db_cols = [row[1] for row in cursor.fetchall()]
        conn.close()
        
        for col in db_cols:
            assert col in poc, f"Column {col} should be documented"
    
    def test_notes_value_in_poc_and_db(self):
        notes = (DOCS_PATH / "notes.txt").read_text().strip()
        
        # Check PoC
        poc = (DOCS_PATH / "PoC.md").read_text()
        assert notes in poc
        
        # Check database
        conn = sqlite3.connect(DOCS_PATH / "cloak.db")
        cursor = conn.cursor()
        cursor.execute("SELECT value FROM artifacts WHERE value=?", (notes,))
        result = cursor.fetchone()
        conn.close()
        assert result is not None
    
    def test_image_referenced_in_poc_exists(self):
        poc = (DOCS_PATH / "PoC.md").read_text()
        images = re.findall(r'!\[.*?\]\((.*?)\)', poc)
        for img in images:
            assert (DOCS_PATH / img).exists(), f"Image {img} should exist"

class TestGitignoreChange:
    """Test .gitignore modification."""
    
    def test_gitignore_exists(self):
        gitignore = Path(__file__).parent.parent.parent / ".gitignore"
        assert gitignore.exists()
    
    def test_gitignore_ends_with_newline(self):
        gitignore = Path(__file__).parent.parent.parent / ".gitignore"
        content = gitignore.read_text()
        if len(content.split('\n')) > 1:
            assert content.endswith('\n'), "Multi-line .gitignore should end with newline"
    
    def test_gitignore_has_env_patterns(self):
        gitignore = Path(__file__).parent.parent.parent / ".gitignore"
        content = gitignore.read_text()
        assert '.env' in content
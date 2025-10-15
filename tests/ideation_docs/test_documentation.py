"""Tests for IdeationDocs markdown and text files."""
import re
from pathlib import Path

DOCS_PATH = Path(__file__).parent.parent.parent / "IdeationDocs"

class TestPoCMarkdown:
    """Test PoC.md documentation."""
    
    def test_poc_file_exists(self):
        assert (DOCS_PATH / "PoC.md").exists()
    
    def test_poc_has_title(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        assert content.startswith("# PoC")
    
    def test_poc_references_image(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        assert "DNScheck.png" in content
        assert (DOCS_PATH / "DNScheck.png").exists()
    
    def test_poc_contains_sql_schemas(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        assert "CREATE TABLE artifacts" in content
        assert "CREATE TABLE detections" in content
    
    def test_poc_artifacts_schema_complete(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        required_cols = ['id', 'type', 'name', 'subject', 'value', 'published_at', 'metadata']
        for col in required_cols:
            assert col in content
    
    def test_poc_has_valid_insert_statements(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        assert "INSERT INTO artifacts" in content
        assert "INSERT INTO detections" in content
    
    def test_poc_uuid_format(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        uuid_pattern = r'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}'
        assert re.search(uuid_pattern, content)
    
    def test_poc_timestamps_iso8601(self):
        content = (DOCS_PATH / "PoC.md").read_text()
        iso_pattern = r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z'
        assert re.search(iso_pattern, content)

class TestSetupMarkdown:
    """Test SETUP.md documentation."""
    
    def test_setup_exists(self):
        assert (DOCS_PATH / "SETUP.md").exists()
    
    def test_setup_has_steps(self):
        content = (DOCS_PATH / "SETUP.md").read_text()
        list_items = re.findall(r'^-\s+', content, re.MULTILINE)
        assert len(list_items) >= 3
    
    def test_setup_has_urls(self):
        content = (DOCS_PATH / "SETUP.md").read_text()
        assert "youtube.com" in content
        assert "dnschecker.org" in content
    
    def test_setup_mentions_key_services(self):
        content = (DOCS_PATH / "SETUP.md").read_text().lower()
        assert "github" in content
        assert "namecheap" in content
        assert "cloudflare" in content or "cloudfare" in content

class TestNotesFile:
    """Test notes.txt file."""
    
    def test_notes_exists(self):
        assert (DOCS_PATH / "notes.txt").exists()
    
    def test_notes_has_content(self):
        content = (DOCS_PATH / "notes.txt").read_text()
        assert content.strip()
    
    def test_notes_contains_test_value(self):
        content = (DOCS_PATH / "notes.txt").read_text()
        assert "cloak-test-001" in content
    
    def test_notes_referenced_in_poc(self):
        notes = (DOCS_PATH / "notes.txt").read_text().strip()
        poc = (DOCS_PATH / "PoC.md").read_text()
        assert notes in poc

class TestBinaryFiles:
    """Test binary files."""
    
    def test_pdf_exists_and_valid(self):
        pdf_path = DOCS_PATH / "Cloak-Adversarial_OSINT_Evasion.pdf"
        assert pdf_path.exists()
        with open(pdf_path, 'rb') as f:
            assert f.read(4) == b'%PDF'
    
    def test_png_exists_and_valid(self):
        png_path = DOCS_PATH / "DNScheck.png"
        assert png_path.exists()
        with open(png_path, 'rb') as f:
            assert f.read(8) == b'\x89PNG\r\n\x1a\n'
    
    def test_database_exists(self):
        assert (DOCS_PATH / "cloak.db").exists()
        assert (DOCS_PATH / "cloak.db").stat().st_size > 0
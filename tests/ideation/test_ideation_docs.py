"""
Comprehensive tests for IdeationDocs directory files.
Tests markdown documentation, database schema, configuration files, and data integrity.
"""

import pytest
import sqlite3
import json
import os
import re
from pathlib import Path
from typing import Dict, List, Tuple


# Base path for IdeationDocs
IDEATION_DOCS_PATH = Path(__file__).parent.parent.parent / "IdeationDocs"


class TestPoC:
    """Test suite for PoC.md markdown documentation."""
    
    @pytest.fixture
    def poc_content(self) -> str:
        """Load PoC.md content."""
        poc_path = IDEATION_DOCS_PATH / "PoC.md"
        with open(poc_path, 'r', encoding='utf-8') as f:
            return f.read()
    
    def test_poc_file_exists(self):
        """Test that PoC.md file exists."""
        poc_path = IDEATION_DOCS_PATH / "PoC.md"
        assert poc_path.exists(), "PoC.md file should exist"
        assert poc_path.is_file(), "PoC.md should be a file"
    
    def test_poc_not_empty(self, poc_content: str):
        """Test that PoC.md is not empty."""
        assert poc_content.strip(), "PoC.md should not be empty"
        assert len(poc_content) > 10, "PoC.md should have substantial content"
    
    def test_poc_has_title(self, poc_content: str):
        """Test that PoC.md has a proper title."""
        lines = poc_content.split('\n')
        assert lines[0].strip() == '# PoC', "First line should be the title '# PoC'"
    
    def test_poc_contains_image_reference(self, poc_content: str):
        """Test that PoC.md references the DNScheck.png image."""
        assert 'DNScheck.png' in poc_content, "Should reference DNScheck.png image"
        assert '![DNS Check](DNScheck.png)' in poc_content, "Should have proper markdown image syntax"
    
    def test_poc_image_file_exists(self):
        """Test that referenced image file actually exists."""
        image_path = IDEATION_DOCS_PATH / "DNScheck.png"
        assert image_path.exists(), "DNScheck.png should exist"
        assert image_path.is_file(), "DNScheck.png should be a file"
        assert image_path.stat().st_size > 0, "DNScheck.png should not be empty"
    
    def test_poc_contains_sql_code_blocks(self, poc_content: str):
        """Test that PoC.md contains SQL code blocks."""
        code_blocks = re.findall(r'```[^\n]*\n(.*?)```', poc_content, re.DOTALL)
        assert len(code_blocks) >= 3, "Should contain at least 3 code blocks"
    
    def test_poc_artifacts_schema_syntax(self, poc_content: str):
        """Test that artifacts table schema is valid SQL."""
        # Extract the CREATE TABLE statement for artifacts
        match = re.search(r'CREATE TABLE artifacts \((.*?)\);', poc_content, re.DOTALL)
        assert match, "Should contain artifacts table CREATE statement"
        
        schema = match.group(1)
        expected_columns = ['id', 'type', 'name', 'subject', 'value', 'published_at', 'metadata']
        
        for col in expected_columns:
            assert col in schema, f"Schema should define column '{col}'"
    
    def test_poc_detections_schema_syntax(self, poc_content: str):
        """Test that detections table schema is valid SQL."""
        match = re.search(r'CREATE TABLE detections \((.*?)\);', poc_content, re.DOTALL)
        assert match, "Should contain detections table CREATE statement"
        
        schema = match.group(1)
        expected_columns = ['id', 'artifact_id', 'provider', 'detected_at', 'raw_response']
        
        for col in expected_columns:
            assert col in schema, f"Schema should define column '{col}'"
    
    def test_poc_insert_statement_valid(self, poc_content: str):
        """Test that INSERT statements are syntactically valid."""
        insert_statements = re.findall(r'INSERT INTO \w+\([^)]+\)\s*VALUES\s*\([^)]+\);', 
                                       poc_content, re.DOTALL)
        assert len(insert_statements) >= 2, "Should contain at least 2 INSERT statements"
        
        # Check artifacts insert
        artifacts_insert = re.search(r'INSERT INTO artifacts\(([^)]+)\)', poc_content)
        assert artifacts_insert, "Should have artifacts INSERT statement"
        columns = [col.strip() for col in artifacts_insert.group(1).split(',')]
        assert 'id' in columns, "artifacts INSERT should include id column"
        assert 'type' in columns, "artifacts INSERT should include type column"
    
    def test_poc_uuid_format(self, poc_content: str):
        """Test that UUIDs in INSERT statements are properly formatted."""
        uuid_pattern = r'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}'
        uuids = re.findall(uuid_pattern, poc_content)
        assert len(uuids) >= 1, "Should contain at least one UUID"
        
        # Check the specific UUID used
        assert '11111111-1111-1111-1111-111111111111' in poc_content, \
            "Should contain the test UUID"
    
    def test_poc_metadata_json_valid(self, poc_content: str):
        """Test that metadata field contains valid JSON."""
        # Extract the metadata JSON string
        match = re.search(r'"metadata":"\{([^}]+)\}"', poc_content)
        if match:
            # Reconstruct the JSON (escaped in the SQL)
            json_str = '{' + match.group(1) + '}'
            # Basic JSON structure check
            assert 'created_via' in json_str or 'record_type' in json_str, \
                "Metadata should contain expected fields"
    
    def test_poc_timestamp_format(self, poc_content: str):
        """Test that timestamps are in ISO 8601 format."""
        timestamp_pattern = r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z'
        timestamps = re.findall(timestamp_pattern, poc_content)
        assert len(timestamps) >= 2, "Should contain at least 2 ISO 8601 timestamps"
    
    def test_poc_references_cloak_test_value(self, poc_content: str):
        """Test that PoC.md references the test value from notes.txt."""
        assert 'cloak-test-001' in poc_content, \
            "Should reference the cloak-test-001 value"


class TestSetupMd:
    """Test suite for SETUP.md documentation."""
    
    @pytest.fixture
    def setup_content(self) -> str:
        """Load SETUP.md content."""
        setup_path = IDEATION_DOCS_PATH / "SETUP.md"
        with open(setup_path, 'r', encoding='utf-8') as f:
            return f.read()
    
    def test_setup_file_exists(self):
        """Test that SETUP.md file exists."""
        setup_path = IDEATION_DOCS_PATH / "SETUP.md"
        assert setup_path.exists(), "SETUP.md file should exist"
        assert setup_path.is_file(), "SETUP.md should be a file"
    
    def test_setup_not_empty(self, setup_content: str):
        """Test that SETUP.md is not empty."""
        assert setup_content.strip(), "SETUP.md should not be empty"
    
    def test_setup_contains_steps(self, setup_content: str):
        """Test that SETUP.md contains setup steps."""
        lines = [line for line in setup_content.split('\n') if line.strip()]
        assert len(lines) >= 3, "Should contain at least 3 setup steps"
    
    def test_setup_has_list_format(self, setup_content: str):
        """Test that SETUP.md uses list format for steps."""
        list_items = re.findall(r'^-\s+', setup_content, re.MULTILINE)
        assert len(list_items) >= 3, "Should have at least 3 list items"
    
    def test_setup_contains_urls(self, setup_content: str):
        """Test that SETUP.md contains reference URLs."""
        url_pattern = r'https?://[^\s\)]+'
        urls = re.findall(url_pattern, setup_content)
        assert len(urls) >= 2, "Should contain at least 2 URLs"
    
    def test_setup_youtube_link_valid(self, setup_content: str):
        """Test that YouTube link is properly formatted."""
        assert 'youtube.com' in setup_content, "Should contain YouTube link"
        youtube_pattern = r'https://www\.youtube\.com/watch\?v=[A-Za-z0-9_-]+'
        assert re.search(youtube_pattern, setup_content), \
            "YouTube link should be properly formatted"
    
    def test_setup_dnschecker_link_valid(self, setup_content: str):
        """Test that dnschecker.org link is present."""
        assert 'dnschecker.org' in setup_content, \
            "Should contain dnschecker.org reference"
    
    def test_setup_mentions_github_student_pack(self, setup_content: str):
        """Test that setup mentions GitHub Student Developer Pack."""
        assert 'GitHub Pro Student Developer Pack' in setup_content or \
               'GitHub' in setup_content, \
               "Should mention GitHub Student Developer Pack"
    
    def test_setup_mentions_namecheap(self, setup_content: str):
        """Test that setup mentions domain provider."""
        assert 'namecheap' in setup_content.lower(), \
            "Should mention namecheap domain provider"
    
    def test_setup_mentions_cloudflare(self, setup_content: str):
        """Test that setup mentions Cloudflare."""
        assert 'cloudflare' in setup_content.lower() or 'Cloudfare' in setup_content, \
            "Should mention Cloudflare setup"
    
    def test_setup_mentions_dns_record(self, setup_content: str):
        """Test that setup mentions DNS record creation."""
        assert 'DNS record' in setup_content or 'dns' in setup_content.lower(), \
            "Should mention DNS record creation"


class TestNotesFile:
    """Test suite for notes.txt file."""
    
    @pytest.fixture
    def notes_content(self) -> str:
        """Load notes.txt content."""
        notes_path = IDEATION_DOCS_PATH / "notes.txt"
        with open(notes_path, 'r', encoding='utf-8') as f:
            return f.read()
    
    def test_notes_file_exists(self):
        """Test that notes.txt file exists."""
        notes_path = IDEATION_DOCS_PATH / "notes.txt"
        assert notes_path.exists(), "notes.txt file should exist"
        assert notes_path.is_file(), "notes.txt should be a file"
    
    def test_notes_not_empty(self, notes_content: str):
        """Test that notes.txt is not empty."""
        assert notes_content.strip(), "notes.txt should not be empty"
    
    def test_notes_contains_test_value(self, notes_content: str):
        """Test that notes.txt contains the expected test value."""
        assert 'cloak-test-001' in notes_content, \
            "notes.txt should contain 'cloak-test-001'"
    
    def test_notes_format_is_simple(self, notes_content: str):
        """Test that notes.txt is a simple single-line or minimal format."""
        lines = notes_content.strip().split('\n')
        assert len(lines) <= 5, "notes.txt should be concise"
    
    def test_notes_consistent_with_poc(self, notes_content: str):
        """Test that notes.txt value is consistent with PoC.md."""
        poc_path = IDEATION_DOCS_PATH / "PoC.md"
        with open(poc_path, 'r', encoding='utf-8') as f:
            poc_content = f.read()
        
        # The value in notes.txt should appear in PoC.md
        notes_value = notes_content.strip()
        assert notes_value in poc_content, \
            f"Value '{notes_value}' from notes.txt should be referenced in PoC.md"


class TestCloakDatabase:
    """Test suite for cloak.db SQLite database."""
    
    @pytest.fixture
    def db_connection(self) -> sqlite3.Connection:
        """Create a connection to the cloak.db database."""
        db_path = IDEATION_DOCS_PATH / "cloak.db"
        conn = sqlite3.connect(db_path)
        yield conn
        conn.close()
    
    def test_database_file_exists(self):
        """Test that cloak.db file exists."""
        db_path = IDEATION_DOCS_PATH / "cloak.db"
        assert db_path.exists(), "cloak.db file should exist"
        assert db_path.is_file(), "cloak.db should be a file"
        assert db_path.stat().st_size > 0, "cloak.db should not be empty"
    
    def test_database_is_valid_sqlite(self, db_connection: sqlite3.Connection):
        """Test that the database is a valid SQLite database."""
        cursor = db_connection.cursor()
        # This will raise an exception if the database is corrupted
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
        tables = cursor.fetchall()
        assert len(tables) > 0, "Database should contain at least one table"
    
    def test_artifacts_table_exists(self, db_connection: sqlite3.Connection):
        """Test that artifacts table exists."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='artifacts'")
        result = cursor.fetchone()
        assert result is not None, "artifacts table should exist"
    
    def test_detections_table_exists(self, db_connection: sqlite3.Connection):
        """Test that detections table exists."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='detections'")
        result = cursor.fetchone()
        assert result is not None, "detections table should exist"
    
    def test_artifacts_schema_structure(self, db_connection: sqlite3.Connection):
        """Test that artifacts table has correct schema."""
        cursor = db_connection.cursor()
        cursor.execute("PRAGMA table_info(artifacts)")
        columns = {row[1]: row[2] for row in cursor.fetchall()}
        
        expected_columns = {
            'id': 'TEXT',
            'type': 'TEXT',
            'name': 'TEXT',
            'subject': 'TEXT',
            'value': 'TEXT',
            'published_at': 'TEXT',
            'metadata': 'TEXT'
        }
        
        for col_name, col_type in expected_columns.items():
            assert col_name in columns, f"Column '{col_name}' should exist in artifacts table"
            assert columns[col_name] == col_type, \
                f"Column '{col_name}' should have type '{col_type}', got '{columns[col_name]}'"
    
    def test_artifacts_id_is_primary_key(self, db_connection: sqlite3.Connection):
        """Test that id is the primary key of artifacts table."""
        cursor = db_connection.cursor()
        cursor.execute("PRAGMA table_info(artifacts)")
        for row in cursor.fetchall():
            if row[1] == 'id':  # column name
                assert row[5] == 1, "id column should be primary key"
                break
    
    def test_detections_schema_structure(self, db_connection: sqlite3.Connection):
        """Test that detections table has correct schema."""
        cursor = db_connection.cursor()
        cursor.execute("PRAGMA table_info(detections)")
        columns = {row[1]: row[2] for row in cursor.fetchall()}
        
        expected_columns = {
            'id': 'TEXT',
            'artifact_id': 'TEXT',
            'provider': 'TEXT',
            'detected_at': 'TEXT',
            'raw_response': 'TEXT'
        }
        
        for col_name, col_type in expected_columns.items():
            assert col_name in columns, f"Column '{col_name}' should exist in detections table"
            assert columns[col_name] == col_type, \
                f"Column '{col_name}' should have type '{col_type}', got '{columns[col_name]}'"
    
    def test_artifacts_contains_test_data(self, db_connection: sqlite3.Connection):
        """Test that artifacts table contains the test data from PoC.md."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT * FROM artifacts WHERE id='11111111-1111-1111-1111-111111111111'")
        row = cursor.fetchone()
        
        assert row is not None, "Should contain test artifact with UUID 11111111-1111-1111-1111-111111111111"
    
    def test_artifacts_test_data_integrity(self, db_connection: sqlite3.Connection):
        """Test that test artifact data is correctly formatted."""
        cursor = db_connection.cursor()
        cursor.execute("""
            SELECT id, type, name, subject, value, published_at, metadata 
            FROM artifacts 
            WHERE id='11111111-1111-1111-1111-111111111111'
        """)
        row = cursor.fetchone()
        
        if row:
            _artifact_id, type_, name, subject, value, published_at, metadata = row
            
            assert type_ == 'dns', f"Type should be 'dns', got '{type_}'"
            assert name == 'cloak.nopejs.me', f"Name should be 'cloak.nopejs.me', got '{name}'"
            assert subject == 'cloak.nopejs.me', f"Subject should be 'cloak.nopejs.me', got '{subject}'"
            assert value == 'cloak-test-001', f"Value should be 'cloak-test-001', got '{value}'"
            
            # Test timestamp format
            timestamp_pattern = r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z'
            assert re.match(timestamp_pattern, published_at), \
                f"published_at should be ISO 8601 format, got '{published_at}'"
            
            # Test metadata is valid JSON
            if metadata:
                metadata_dict = json.loads(metadata)
                assert 'created_via' in metadata_dict or 'record_type' in metadata_dict, \
                    "Metadata should contain expected fields"
    
    def test_detections_contains_test_data(self, db_connection: sqlite3.Connection):
        """Test that detections table contains test data."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT * FROM detections WHERE id='det-0001'")
        row = cursor.fetchone()
        
        assert row is not None, "Should contain test detection with id 'det-0001'"
    
    def test_detections_test_data_integrity(self, db_connection: sqlite3.Connection):
        """Test that test detection data is correctly formatted."""
        cursor = db_connection.cursor()
        cursor.execute("""
            SELECT id, artifact_id, provider, detected_at, raw_response 
            FROM detections 
            WHERE id='det-0001'
        """)
        row = cursor.fetchone()
        
        if row:
            _det_id, artifact_id, provider, detected_at, raw_response = row
            
            assert artifact_id == '11111111-1111-1111-1111-111111111111', \
                "artifact_id should match test artifact UUID"
            assert provider == 'github_search', f"Provider should be 'github_search', got '{provider}'"
            
            # Test timestamp format
            timestamp_pattern = r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z'
            assert re.match(timestamp_pattern, detected_at), \
                f"detected_at should be ISO 8601 format, got '{detected_at}'"
            
            # Test raw_response is valid JSON
            if raw_response:
                response_dict = json.loads(raw_response)
                assert 'url' in response_dict, "raw_response should contain 'url' field"
    
    def test_foreign_key_relationship(self, db_connection: sqlite3.Connection):
        """Test that detection's artifact_id references an existing artifact."""
        cursor = db_connection.cursor()
        
        # Get all artifact_ids from detections
        cursor.execute("SELECT DISTINCT artifact_id FROM detections")
        detection_artifact_ids = [row[0] for row in cursor.fetchall()]
        
        # Verify each references an existing artifact
        for artifact_id in detection_artifact_ids:
            cursor.execute("SELECT id FROM artifacts WHERE id=?", (artifact_id,))
            result = cursor.fetchone()
            assert result is not None, \
                f"artifact_id '{artifact_id}' in detections should reference an existing artifact"
    
    def test_database_no_null_primary_keys(self, db_connection: sqlite3.Connection):
        """Test that primary keys are not null in any table."""
        cursor = db_connection.cursor()
        
        # Check artifacts
        cursor.execute("SELECT COUNT(*) FROM artifacts WHERE id IS NULL")
        assert cursor.fetchone()[0] == 0, "No artifact should have NULL id"
        
        # Check detections
        cursor.execute("SELECT COUNT(*) FROM detections WHERE id IS NULL")
        assert cursor.fetchone()[0] == 0, "No detection should have NULL id"
    
    def test_artifacts_value_consistent_with_notes(self, db_connection: sqlite3.Connection):
        """Test that artifact value matches notes.txt content."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT value FROM artifacts WHERE id='11111111-1111-1111-1111-111111111111'")
        row = cursor.fetchone()
        
        if row:
            artifact_value = row[0]
            
            # Load notes.txt
            notes_path = IDEATION_DOCS_PATH / "notes.txt"
            with open(notes_path, 'r', encoding='utf-8') as f:
                notes_content = f.read().strip()
            
            assert artifact_value == notes_content, \
                f"Artifact value '{artifact_value}' should match notes.txt content '{notes_content}'"


class TestGitignore:
    """Test suite for .gitignore file."""
    
    @pytest.fixture
    def gitignore_content(self) -> str:
        """Load .gitignore content."""
        gitignore_path = Path(__file__).parent.parent.parent / ".gitignore"
        with open(gitignore_path, 'r', encoding='utf-8') as f:
            return f.read()
    
    def test_gitignore_exists(self):
        """Test that .gitignore file exists."""
        gitignore_path = Path(__file__).parent.parent.parent / ".gitignore"
        assert gitignore_path.exists(), ".gitignore file should exist"
    
    def test_gitignore_not_empty(self, gitignore_content: str):
        """Test that .gitignore is not empty."""
        assert gitignore_content.strip(), ".gitignore should not be empty"
    
    def test_gitignore_no_trailing_ideation_docs(self, gitignore_content: str):
        """Test that .gitignore does not exclude IdeationDocs (as per the change)."""
        lines = gitignore_content.split('\n')
        # The last line should not be 'IdeationDocs' without newline
        assert lines[-1] != 'IdeationDocs', \
            "IdeationDocs should not be excluded at the end without newline"
    
    def test_gitignore_standard_patterns(self, gitignore_content: str):
        """Test that .gitignore contains standard patterns."""
        # Should contain common ignore patterns
        
        # Check for at least some common patterns (based on the diff context)
        assert '.env' in gitignore_content, "Should ignore .env files"
    
    def test_gitignore_proper_line_endings(self, gitignore_content: str):
        """Test that .gitignore has proper line endings."""
        # Should not end with a line without newline after other content
        if gitignore_content:
            # If file has content, last character should be newline or content should be single line
            lines = gitignore_content.split('\n')
            if len(lines) > 1:
                # Multi-line file should end with newline
                assert gitignore_content.endswith('\n'), \
                    ".gitignore should end with a newline character"


class TestIntegration:
    """Integration tests across multiple files."""
    
    def test_poc_schema_matches_database_schema(self):
        """Test that PoC.md schema matches actual database schema."""
        # Load PoC.md
        poc_path = IDEATION_DOCS_PATH / "PoC.md"
        with open(poc_path, 'r', encoding='utf-8') as f:
            poc_content = f.read()
        
        # Connect to database
        db_path = IDEATION_DOCS_PATH / "cloak.db"
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        try:
            # Check artifacts table
            cursor.execute("PRAGMA table_info(artifacts)")
            db_columns = [row[1] for row in cursor.fetchall()]
            
            # Extract schema from PoC.md
            poc_match = re.search(r'CREATE TABLE artifacts \((.*?)\);', poc_content, re.DOTALL)
            assert poc_match, "PoC.md should contain artifacts schema"
            
            # Verify column names match
            for col in db_columns:
                assert col in poc_match.group(1), \
                    f"Database column '{col}' should be documented in PoC.md schema"
        finally:
            conn.close()
    
    def test_all_documentation_files_present(self):
        """Test that all expected documentation files are present."""
        expected_files = [
            "PoC.md",
            "SETUP.md",
            "notes.txt",
            "cloak.db",
            "DNScheck.png",
            "Cloak-Adversarial_OSINT_Evasion.pdf"
        ]
        
        for filename in expected_files:
            file_path = IDEATION_DOCS_PATH / filename
            assert file_path.exists(), f"Expected file '{filename}' should exist in IdeationDocs"
    
    def test_cross_reference_consistency(self):
        """Test that cross-references between files are consistent."""
        # notes.txt value should appear in both PoC.md and database
        notes_path = IDEATION_DOCS_PATH / "notes.txt"
        with open(notes_path, 'r', encoding='utf-8') as f:
            notes_value = f.read().strip()
        
        # Check PoC.md
        poc_path = IDEATION_DOCS_PATH / "PoC.md"
        with open(poc_path, 'r', encoding='utf-8') as f:
            poc_content = f.read()
        assert notes_value in poc_content, \
            f"Value '{notes_value}' from notes.txt should be in PoC.md"
        
        # Check database
        db_path = IDEATION_DOCS_PATH / "cloak.db"
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        try:
            cursor.execute("SELECT value FROM artifacts WHERE value=?", (notes_value,))
            result = cursor.fetchone()
            assert result is not None, \
                f"Value '{notes_value}' from notes.txt should be in database"
        finally:
            conn.close()
    
    def test_pdf_file_is_valid(self):
        """Test that the PDF file exists and has reasonable size."""
        pdf_path = IDEATION_DOCS_PATH / "Cloak-Adversarial_OSINT_Evasion.pdf"
        assert pdf_path.exists(), "PDF file should exist"
        assert pdf_path.is_file(), "PDF should be a file"
        
        # PDF should be reasonably sized (not empty, not too small)
        file_size = pdf_path.stat().st_size
        assert file_size > 1000, "PDF should be larger than 1KB"
        assert file_size < 100_000_000, "PDF should be smaller than 100MB"
        
        # Check PDF magic number
        with open(pdf_path, 'rb') as f:
            header = f.read(4)
            assert header == b'%PDF', "File should start with PDF magic number"
    
    def test_image_file_is_valid(self):
        """Test that the PNG file exists and has valid format."""
        png_path = IDEATION_DOCS_PATH / "DNScheck.png"
        assert png_path.exists(), "PNG file should exist"
        assert png_path.is_file(), "PNG should be a file"
        
        # PNG should be reasonably sized
        file_size = png_path.stat().st_size
        assert file_size > 100, "PNG should be larger than 100 bytes"
        
        # Check PNG magic number
        with open(png_path, 'rb') as f:
            header = f.read(8)
            assert header == b'\x89PNG\r\n\x1a\n', "File should start with PNG magic number"
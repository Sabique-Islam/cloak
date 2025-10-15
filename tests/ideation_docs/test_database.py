"""Tests for cloak.db SQLite database."""
import sqlite3
import json
import re
from pathlib import Path

DOCS_PATH = Path(__file__).parent.parent.parent / "IdeationDocs"
DB_PATH = DOCS_PATH / "cloak.db"

class TestDatabaseStructure:
    """Test database schema and structure."""
    
    def test_database_is_valid_sqlite(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
        tables = cursor.fetchall()
        conn.close()
        assert len(tables) > 0
    
    def test_artifacts_table_exists(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='artifacts'")
        result = cursor.fetchone()
        conn.close()
        assert result is not None
    
    def test_detections_table_exists(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='detections'")
        result = cursor.fetchone()
        conn.close()
        assert result is not None
    
    def test_artifacts_schema(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("PRAGMA table_info(artifacts)")
        columns = {row[1]: row[2] for row in cursor.fetchall()}
        conn.close()
        
        expected = {
            'id': 'TEXT', 'type': 'TEXT', 'name': 'TEXT',
            'subject': 'TEXT', 'value': 'TEXT',
            'published_at': 'TEXT', 'metadata': 'TEXT'
        }
        for col, typ in expected.items():
            assert col in columns
            assert columns[col] == typ
    
    def test_detections_schema(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("PRAGMA table_info(detections)")
        columns = {row[1]: row[2] for row in cursor.fetchall()}
        conn.close()
        
        expected = {
            'id': 'TEXT', 'artifact_id': 'TEXT',
            'provider': 'TEXT', 'detected_at': 'TEXT',
            'raw_response': 'TEXT'
        }
        for col, typ in expected.items():
            assert col in columns
            assert columns[col] == typ

class TestDatabaseData:
    """Test database data integrity."""
    
    def test_artifacts_has_test_data(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM artifacts WHERE id='11111111-1111-1111-1111-111111111111'")
        row = cursor.fetchone()
        conn.close()
        assert row is not None
    
    def test_artifact_test_values(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("""SELECT type, name, value FROM artifacts 
                         WHERE id='11111111-1111-1111-1111-111111111111'""")
        row = cursor.fetchone()
        conn.close()
        
        if row:
            assert row[0] == 'dns'
            assert row[1] == 'cloak.nopejs.me'
            assert row[2] == 'cloak-test-001'
    
    def test_detections_has_test_data(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM detections WHERE id='det-0001'")
        row = cursor.fetchone()
        conn.close()
        assert row is not None
    
    def test_foreign_key_integrity(self):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT DISTINCT artifact_id FROM detections")
        artifact_ids = [r[0] for r in cursor.fetchall()]
        
        for aid in artifact_ids:
            cursor.execute("SELECT id FROM artifacts WHERE id=?", (aid,))
            assert cursor.fetchone() is not None
        conn.close()
    
    def test_artifact_value_matches_notes(self):
        notes_content = (DOCS_PATH / "notes.txt").read_text().strip()
        
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute("SELECT value FROM artifacts WHERE id='11111111-1111-1111-1111-111111111111'")
        row = cursor.fetchone()
        conn.close()
        
        if row:
            assert row[0] == notes_content
"""Private greenfield configuration and bounded disposable reset regression tests."""
import fcntl
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest

SCRIPTS=Path(__file__).resolve().parent
class LocalWorkspaceTests(unittest.TestCase):
    def test_owner_only_configuration_is_idempotent(self):
        with tempfile.TemporaryDirectory() as root:
            scripts=Path(root)/'scripts';scripts.mkdir()
            shutil.copy(SCRIPTS/'local-env.py',scripts)
            command=['python3',str(scripts/'local-env.py')]
            subprocess.run(command,check=True,capture_output=True)
            config=Path(root)/'.sentrame/local.env';first=config.read_bytes()
            self.assertIn(b'SENTRAME_PASSPHRASE',first)
            self.assertNotIn(b'SENTRAME_BROWSER_TOKEN',first)
            self.assertNotIn(b'SENTRAME_SERVER_TOKEN',first)
            self.assertEqual(config.stat().st_mode&0o777,0o600)
            subprocess.run(command,check=True,capture_output=True)
            self.assertEqual(first,config.read_bytes())

    def test_reset_respects_registry_lease_and_preserves_singleton_files(self):
        with tempfile.TemporaryDirectory() as root:
            scripts=Path(root)/'scripts';scripts.mkdir()
            shutil.copy(SCRIPTS/'reset-local.py',scripts)
            folder=Path(root)/'.sentrame';folder.mkdir()
            old=folder/'local.db';old.write_text('preserved')
            current=folder/'workspaces';current.mkdir()
            (current/'apps.json').write_text('{"version":1,"revision":1,"apps":[]}')
            with (current/'.apps.lock').open('a+') as lease:
                fcntl.flock(lease,fcntl.LOCK_EX|fcntl.LOCK_NB)
                result=subprocess.run(['python3',str(scripts/'reset-local.py')],capture_output=True)
                self.assertNotEqual(result.returncode,0)
                self.assertIn(b'Stop the local service',result.stderr)
                self.assertTrue((current/'apps.json').exists())
            self.assertEqual(old.read_text(),'preserved')

if __name__=='__main__':unittest.main()

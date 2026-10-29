"""Readiness-probe regression tests, with no Docker daemon or network required."""
import importlib.util
import pathlib
import unittest
import urllib.error
from unittest.mock import Mock

spec = importlib.util.spec_from_file_location("container_smoke", pathlib.Path(__file__).with_name("container-smoke.py"))
smoke = importlib.util.module_from_spec(spec)
spec.loader.exec_module(smoke)

class ProbeTests(unittest.TestCase):
    def test_startup_connection_reset_is_retried(self):
        response = Mock(status=200)
        response.__enter__ = Mock(return_value=response)
        response.__exit__ = Mock(return_value=False)
        request = Mock(side_effect=[ConnectionResetError("starting"), response])
        pause = Mock()
        smoke.wait_for_ready(request, pause=pause)
        self.assertEqual(request.call_count, 2)
        pause.assert_called_once_with(1)

    def test_fixture_is_inside_receiver_window_when_host_clock_is_ahead(self):
        server_now = 1_800_000_000_000_000_000
        from unittest.mock import patch
        with patch.object(smoke.time, "time_ns", return_value=server_now + 60_000_000_000):
            payload = smoke.smoke_payload("test", server_now)
        record = payload["resourceLogs"][0]["scopeLogs"][0]["logRecords"][0]
        timestamp = int(record["timeUnixNano"])
        self.assertLess(timestamp, server_now)
        self.assertGreater(timestamp, server_now - 30 * 86400 * 1_000_000_000)
        self.assertNotIn("attributes", record, "standard receipt needs no journey attributes")

    def test_bad_endpoint_is_not_hidden_by_retries(self):
        error = urllib.error.HTTPError("http://local/ready", 404, "missing", {}, None)
        request = Mock(side_effect=error)
        try:
            with self.assertRaises(urllib.error.HTTPError):
                smoke.wait_for_ready(request, pause=Mock())
        finally:
            error.close()

    def test_unavailable_startup_has_a_bounded_deadline(self):
        request = Mock(side_effect=ConnectionResetError("starting"))
        with self.assertRaisesRegex(RuntimeError, "readiness"):
            smoke.wait_for_ready(request, pause=Mock(), clock=Mock(side_effect=[0, 1, 91]))
        self.assertEqual(request.call_count, 1)

if __name__ == "__main__":
    unittest.main()

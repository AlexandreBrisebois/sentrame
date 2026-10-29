#!/usr/bin/env python3
"""Resolve official version tags once; all CI jobs use the same builder digests."""
import json
import os
import re
import subprocess

images = {"go_image": "golang:1.27.3-alpine3.24", "node_image": "node:26.13.0-alpine3.24"}
results = {}
for key, tag in images.items():
    raw = subprocess.check_output(["docker", "buildx", "imagetools", "inspect", tag, "--format", "{{json .Manifest}}"], text=True)
    digest = json.loads(raw).get("digest", "")
    if not re.fullmatch(r"sha256:[a-f0-9]{64}", digest):
        raise SystemExit("Registry did not return an immutable builder digest")
    results[key] = tag + "@" + digest
if os.environ.get("GITHUB_OUTPUT"):
    with open(os.environ["GITHUB_OUTPUT"], "a") as output:
        for key, value in results.items():
            output.write(key + "=" + value + "\n")
if os.environ.get("GITHUB_STEP_SUMMARY"):
    with open(os.environ["GITHUB_STEP_SUMMARY"], "a") as summary:
        summary.write("Builder manifests used for every architecture:\n")
        for value in results.values():
            summary.write("- `" + value + "`\n")
print(json.dumps(results, indent=2))

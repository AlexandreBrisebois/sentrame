#!/usr/bin/env python3
"""Reset only the explicitly disposable default local workspace; preserve old DBs."""
import fcntl, pathlib, socket, shutil
root=pathlib.Path(__file__).resolve().parent.parent
folder=root/'.sentrame/workspaces'
folder.mkdir(mode=0o700,parents=True,exist_ok=True)
with (folder/'.apps.lock').open('a+') as lease:
    try: fcntl.flock(lease,fcntl.LOCK_EX|fcntl.LOCK_NB)
    except BlockingIOError: raise SystemExit('Stop the local service before resetting its disposable workspace.')
    with socket.socket() as sock:
        sock.settimeout(.5)
        if sock.connect_ex(('127.0.0.1',8096))==0:
            raise SystemExit('Port8096 is occupied. Stop the local service; no files were removed.')
    for name in ['apps.json','.apps-initialized']:
        p=folder/name
        if p.is_symlink(): raise SystemExit('Refusing symlink in disposable workspace.')
    apps=folder/'apps'
    if apps.is_symlink(): raise SystemExit('Refusing symlink in disposable workspace.')
    if apps.exists(): shutil.rmtree(apps)
    for name in ['apps.json','.apps-initialized']:
        (folder/name).unlink(missing_ok=True)
print('Disposable local apps removed. Register apps again. Existing single-app files and owner configuration were preserved.')

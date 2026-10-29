#!/usr/bin/env python3
"""Reset only disposable default local synthetic files while holding their process lease."""
import fcntl, pathlib, socket
root=pathlib.Path(__file__).resolve().parent.parent
folder=root/'.sentrame';folder.mkdir(mode=0o700,exist_ok=True)
with (folder/'local.db.lock').open('a+') as lease:
    try: fcntl.flock(lease,fcntl.LOCK_EX|fcntl.LOCK_NB)
    except BlockingIOError: raise SystemExit('Stop the local service before resetting its synthetic database.')
    with socket.socket() as sock:
        sock.settimeout(.5)
        if sock.connect_ex(('127.0.0.1',8096))==0:
            raise SystemExit('Port8096 is occupied. Stop the local service before resetting; no files were removed.')
    for name in ['local.db','local.db-wal','local.db-shm','local.env']:
        p=folder/name
        if p.exists(): p.unlink()
print('Local synthetic state reset. Credentials will rotate; old queued requests are no longer authorized.')

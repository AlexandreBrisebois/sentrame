#!/usr/bin/env python3
"""Create local synthetic secrets only when absent. Never print their values."""
import pathlib, secrets, os
root=pathlib.Path(__file__).resolve().parent.parent
folder=root/'.sentrame'; folder.mkdir(mode=0o700,exist_ok=True)
p=folder/'local.env'
if not p.exists():
    vals={'SENTRAME_PASSPHRASE':secrets.token_urlsafe(24),'SENTRAME_BROWSER_TOKEN':secrets.token_hex(32),'SENTRAME_SERVER_TOKEN':secrets.token_hex(32),'SENTRAME_ORIGIN':'http://127.0.0.1:8096'}
    fd=os.open(p,os.O_WRONLY|os.O_CREAT|os.O_EXCL,0o600)
    with os.fdopen(fd,'w') as f:
        for k,v in vals.items(): f.write(f"export {k}='{v}'\n")
print('Local synthetic configuration: .sentrame/local.env (secrets omitted)')

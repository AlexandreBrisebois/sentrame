#!/usr/bin/env python3
"""Isolated loopback synthetic load through the actual receiver/store/read/SSE path."""
import concurrent.futures, http.cookiejar, json, os, pathlib, secrets, socket, statistics, subprocess, tempfile, threading, time, urllib.request
root=pathlib.Path(__file__).resolve().parent.parent

def main():
    with tempfile.TemporaryDirectory(prefix='sentrame-qualification-') as tmp:
        with socket.socket() as sock:
            sock.bind(('127.0.0.1',0)); port=sock.getsockname()[1]
        origin=f'http://127.0.0.1:{port}'; password=secrets.token_urlsafe(24); producer=secrets.token_hex(32)
        env=dict(os.environ,SENTRAME_ADDR=f'127.0.0.1:{port}',SENTRAME_ORIGIN=origin,SENTRAME_DB=f'{tmp}/load.db',SENTRAME_PASSPHRASE=password)
        proc=subprocess.Popen([str(root/'bin/sentrame')],cwd=root,env=env,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
        errors=[]; durations=[]; query_times=[]; samples=[]; stop=threading.Event()
        try:
            for _ in range(100):
                try:
                    urllib.request.urlopen(origin+'/health',timeout=.2).close(); break
                except OSError: time.sleep(.05)
            else: raise RuntimeError('isolated service failed readiness')
            jar=http.cookiejar.CookieJar(); owner=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))
            req=urllib.request.Request(origin+'/api/login',data=json.dumps({'passphrase':password}).encode(),headers={'Origin':origin,'Content-Type':'application/json'})
            with owner.open(req,timeout=5) as r: csrf=json.load(r)['csrf']
            headers={'Origin':origin,'Content-Type':'application/json','X-CSRF-Token':csrf}
            req=urllib.request.Request(origin+'/api/apps',data=json.dumps({'name':'Load qualification'}).encode(),headers=headers)
            with owner.open(req,timeout=5) as r: app_id=json.load(r)['id']
            with owner.open(origin+'/api/apps/'+app_id+'/connection',timeout=5) as r: producer=json.load(r)['browser']
            def view():
                while not stop.is_set():
                    start=time.perf_counter()
                    try:
                        with owner.open(origin+'/api/apps/'+app_id+'/snapshot',timeout=5) as r: json.load(r)
                        query_times.append(time.perf_counter()-start)
                    except Exception as e: errors.append(type(e).__name__)
                    stop.wait(.1)
            def stream():
                try:
                    with owner.open(origin+'/api/apps/'+app_id+'/events',timeout=6) as r:
                        while not stop.is_set():
                            if not r.readline(): break
                except Exception as e:
                    if not stop.is_set(): errors.append(type(e).__name__)
            def memory():
                while not stop.is_set():
                    try: samples.append(int(subprocess.check_output(['ps','-o','rss=','-p',str(proc.pid)],text=True).strip())*1024)
                    except Exception: pass
                    stop.wait(.2)
            threads=[threading.Thread(target=f,daemon=True) for f in (view,stream,memory)]
            for t in threads: t.start()
            def send(index):
                body={'resourceLogs':[{'scopeLogs':[{'logRecords':[{'eventName':'qualification.synthetic','timeUnixNano':str(time.time_ns())}]}]}]}
                req=urllib.request.Request(origin+'/v1/logs',data=json.dumps(body).encode(),headers={'Authorization':'Bearer '+producer,'Content-Type':'application/json'})
                start=time.perf_counter()
                try:
                    with urllib.request.urlopen(req,timeout=5) as r:
                        if r.status!=200: errors.append(f'HTTP {r.status}')
                    durations.append(time.perf_counter()-start)
                except Exception as e: errors.append(type(e).__name__)
            with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
                for rate,seconds in [(10,10),(100,10)]:
                    futures=[]; start=time.perf_counter()
                    for i in range(rate*seconds):
                        wait=start+i/rate-time.perf_counter()
                        if wait>0: time.sleep(wait)
                        futures.append(pool.submit(send,i))
                    for f in futures: f.result()
            with owner.open(origin+'/api/apps/'+app_id+'/snapshot',timeout=5) as r: snap=json.load(r)
            if len(snap['evidence'])!=1100: errors.append('acknowledged receipt count mismatch')
            sizes=sum(p.stat().st_size for p in pathlib.Path(tmp).rglob('telemetry.db*'))
            percentile=lambda values: sorted(values)[int(.95*(len(values)-1))] if values else None
            result={'scope':'isolated Mac loopback, JSON standard logs without optional journey fields; one query viewer + one SSE viewer + 1-second cleanup','requests':1100,'successfulRequests':len(durations),'phases':'10/s for10s, then100/s for10s;4 bounded senders','ingestP95Seconds':percentile(durations),'snapshotP95Seconds':percentile(query_times),'snapshotSamples':len(query_times),'peakObservedRSSBytes':max(samples) if samples else None,'databaseAndWALBytes':sizes,'errors':errors,'limits':'No 100,000-record retention soak, physical power-loss, Home Server, Tunnel, browser/.NET exporter, or external telemetry qualification.'}
            result['passed']=not errors and result['ingestP95Seconds']<1 and result['snapshotP95Seconds']<2 and result['peakObservedRSSBytes'] is not None and result['peakObservedRSSBytes']<512*1024*1024
            (root/'.sentrame'/'qualification.json').write_text(json.dumps(result,indent=2)+'\n')
            print(json.dumps(result,indent=2))
            if not result['passed']: raise SystemExit(1)
        finally:
            stop.set();proc.terminate()
            try: proc.wait(timeout=5)
            except subprocess.TimeoutExpired: proc.kill();proc.wait()
if __name__=='__main__': main()

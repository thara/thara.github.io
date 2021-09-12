---
title: Gunicorn起動中にUnixドメインソケットが消える
tags: [python]
---

会社のデベロッパーブログでも書いたけど再現手順は記載してなかったので、改めて。

Gunicorn 19.6.0 をthread workerで稼働していると、max_requestsに達してworker threadが再起動したときにUnixドメインソケットが削除される現象にあった。

以下、再現手順。   


## 事前準備

- Python 3.6
- [h2o](https://github.com/h2o/h2o)  : Unixドメインソケット プロキシ用
- [wrk](https://github.com/wg/wrk)  : ベンチマークツール

```
$ pip install gunicorn==19.6.0 bottle==0.12.13
```

以下のPythonモジュールを作成しておく。

web.py
(WSGIアプリケーション)

```python
# -*- coding:utf-8 -*-
from bottle import route, run, default_app

@route('/')
def index():
  return 'Hello World'

app = default_app()
```


config.py
(gunicorn設定)

```python
bind = 'unix:/tmp/gunicorn.sock'
backlog = 2048

workers = 5
threads = 5
worker_class = 'sync'
worker_connections = 1000
timeout = 30
keepalive = 2

max_requests = 512

spew = False

daemon = False
pidfile = None
umask = 0
user = None
group = None
tmp_upload_dir = None

#   Logging
errorlog = '-'
loglevel = 'info'
accesslog = '-'
access_log_format = '%(h)s %(l)s %(u)s %(t)s "%(r)s" %(s)s %(b)s "%(f)s" "%(a)s"'

# Process naming
proc_name = None

# Server hooks
def post_fork(server, worker):
    server.log.info("Worker spawned (pid: %s)", worker.pid)

def pre_fork(server, worker):
    pass

def pre_exec(server):
    server.log.info("Forked child, re-executing.")

def when_ready(server):
    server.log.info("Server is ready. Spawning workers")

def worker_int(worker):
    worker.log.info("worker received INT or QUIT signal")

    ## get traceback info
    import threading, sys, traceback
    id2name = dict([(th.ident, th.name) for th in threading.enumerate()])
    code = []
    for threadId, stack in sys._current_frames().items():
        code.append("\n# Thread: %s(%d)" % (id2name.get(threadId,""),
            threadId))
        for filename, lineno, name, line in traceback.extract_stack(stack):
            code.append('File: "%s", line %d, in %s' % (filename,
                lineno, name))
            if line:
                code.append("  %s" % (line.strip()))
    worker.log.debug("\n".join(code))

def worker_abort(worker):
    worker.log.info("worker received SIGABRT signal")
```


h2o.conf
(h2o設定)

```
listen:
  port: 18080
hosts:
  default:
    paths:
      "/":
        proxy.reverse.url: http://[unix:/tmp/gunicorn.sock]/
        proxy.preserve-host: ON
    access-log: /dev/stdout
```


## サーバー起動


gunicorn 起動

```
$ gunicorn web:app -c config.py
```


h2o 起動

```
$ h2o -c h2o.conf
```


## 再現

```
$ wrk -t10 -c200 -d10s http://127.0.0.1:18080
```

以下のようなエラーログが出力される。   
（場合によっては、スレッド数や接続数等のオプションを変更する必要があるかもしれない）

```
[2017-02-18 01:26:57 +0900] [56927] [INFO] Autorestarting worker after current request.
 - - [18/Feb/2017:01:26:57 +0900] "GET / HTTP/1.1" 200 11 "-" "-"
 - - [18/Feb/2017:01:26:57 +0900] "GET / HTTP/1.1" 200 11 "-" "-"
 - - [18/Feb/2017:01:26:57 +0900] "GET / HTTP/1.1" 200 11 "-" "-"
[2017-02-18 01:26:57 +0900] [56927] [ERROR] Exception in worker process
Traceback (most recent call last):
  File "/Users/thara/.pyenv/versions/3.6.0/envs/gunicorn_bug/lib/python3.6/site-packages/gunicorn/arbiter.py", line 557, in spawn_worker
    worker.init_process()
  File "/Users/thara/.pyenv/versions/3.6.0/envs/gunicorn_bug/lib/python3.6/site-packages/gunicorn/workers/gthread.py", line 109, in init_process
    super(ThreadWorker, self).init_process()
  File "/Users/thara/.pyenv/versions/3.6.0/envs/gunicorn_bug/lib/python3.6/site-packages/gunicorn/workers/base.py", line 132, in init_process
    self.run()
  File "/Users/thara/.pyenv/versions/3.6.0/envs/gunicorn_bug/lib/python3.6/site-packages/gunicorn/workers/gthread.py", line 240, in run
    s.close()
  File "/Users/thara/.pyenv/versions/3.6.0/envs/gunicorn_bug/lib/python3.6/site-packages/gunicorn/sock.py", line 123, in close
    os.unlink(self.cfg_addr)
FileNotFoundError: [Errno 2] No such file or directory: '/tmp/gunicorn.sock'
```

この状態になると、Unixドメインソケットのファイルが紛失するため、gunicornに対してリバースプロキシ（この例だとH2O）から接続できなくなり、
HTTPリクエストしたクライアントには502が返される。

該当issueは [これ](https://github.com/benoitc/gunicorn/issues/1298)。   
該当箇所は19.6.0での修正箇所であるため、19.6.0より前のバージョンではこの現象は発生しないと思われる。
また、このissueはすでに対応されcloseされているが、まだリリースには至っていない。


## 回避方法

現状では、Unixドメインソケットを使わず、portのlistenのみで対応するのが確実。

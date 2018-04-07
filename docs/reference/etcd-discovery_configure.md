## etcd-discovery configure

Configure certs for etcd-discovery

### Synopsis

Configure certs for etcd-discovery

```
etcd-discovery configure [flags]
```

### Options

```
      --addr string       Address of server ip (default "127.0.0.1")
      --cert-dir string   Path to directory where pki files are stored. (default "etcd.local.config/certificates")
  -h, --help              help for configure
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --enable-analytics                 Send usage events to Google Analytics (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [etcd-discovery](etcd-discovery.md)	 - etcd discovery server


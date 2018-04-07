## etcd-discovery run

Launch a etcd discovery server

### Synopsis

Launch a etcd discovery server

```
etcd-discovery run [flags]
```

### Options

```
      --audit-log-format string                        Format of saved audits. "legacy" indicates 1-line text format for each event. "json" indicates structured json format. Requires the 'AdvancedAuditing' feature gate. Known formats are legacy,json. (default "json")
      --audit-log-maxage int                           The maximum number of days to retain old audit log files based on the timestamp encoded in their filename.
      --audit-log-maxbackup int                        The maximum number of old audit log files to retain.
      --audit-log-maxsize int                          The maximum size in megabytes of the audit log file before it gets rotated.
      --audit-log-path string                          If set, all requests coming to the apiserver will be logged to this file.  '-' means standard out.
      --audit-policy-file string                       Path to the file that defines the audit policy configuration. Requires the 'AdvancedAuditing' feature gate. With AdvancedAuditing, a profile is required to enable auditing.
      --audit-webhook-batch-buffer-size int            The size of the buffer to store events before batching and sending to the webhook. Only used in batch mode. (default 10000)
      --audit-webhook-batch-initial-backoff duration   The amount of time to wait before retrying the first failed requests. Only used in batch mode. (default 10s)
      --audit-webhook-batch-max-size int               The maximum size of a batch sent to the webhook. Only used in batch mode. (default 400)
      --audit-webhook-batch-max-wait duration          The amount of time to wait before force sending the batch that hadn't reached the max size. Only used in batch mode. (default 30s)
      --audit-webhook-batch-throttle-burst int         Maximum number of requests sent at the same moment if ThrottleQPS was not utilized before. Only used in batch mode. (default 15)
      --audit-webhook-batch-throttle-qps float32       Maximum average number of requests per second. Only used in batch mode. (default 10)
      --audit-webhook-config-file string               Path to a kubeconfig formatted file that defines the audit webhook configuration. Requires the 'AdvancedAuditing' feature gate.
      --audit-webhook-mode string                      Strategy for sending audit events. Blocking indicates sending events should block server responses. Batch causes the webhook to buffer and send events asynchronously. Known modes are batch,blocking. (default "batch")
      --bind-address ip                                The IP address on which to listen for the --secure-port port. The associated interface(s) must be reachable by the rest of the cluster, and by CLI/web clients. If blank, all interfaces will be used (0.0.0.0). (default 0.0.0.0)
      --cert-dir string                                The directory where the TLS certs are located. If --peer-cert-file and --peer-private-key-file are provided, this flag will be ignored. (default "etcd.local.config/certificates")
      --cert-file string                               File containing the default x509 Certificate used for SSL/TLS connections to etcd. When this option is set, advertise-client-urls can use the HTTPS schema. If HTTPS serving is enabled, and --cert-file and --private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir.
      --contention-profiling                           Enable lock contention profiling, if profiling is enabled
      --enable-swagger-ui                              Enables swagger ui on the apiserver at /swagger-ui
      --etcd-backup-store string                       Backup store location
      --etcd-cluster-name string                       Name of cluster
      --etcd-cluster-size int                          Size of cluster size
      --etcd-data-dir string                           Directory for storing etcd data
  -h, --help                                           help for run
      --initial-cluster stringToString                 Initial cluster configuration (default [])
      --initial-cluster-state ClusterState             Initial cluster state (default New)
      --peer-cert-file string                          File containing the default x509 Certificate used for SSL/TLS connections between peers. This will be used both for listening on the peer address as well as sending requests to other peers. If HTTPS serving is enabled, and --peer-cert-file and --peer-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir.
      --peer-private-key-file string                   File containing the default x509 private key matching --peer-cert-file.
      --peer-trusted-ca-file string                    File containing the certificate authority will used for secure access from peer etcd servers. This must be a valid PEM-encoded CA bundle.
      --private-key-file string                        File containing the default x509 private key matching --cert-file.
      --profiling                                      Enable profiling via web interface host:port/debug/pprof/ (default true)
      --secure-port int                                The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all. (default 2381)
      --trusted-ca-file string                         File containing the certificate authority will used for secure client-to-server communication. This must be a valid PEM-encoded CA bundle.
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


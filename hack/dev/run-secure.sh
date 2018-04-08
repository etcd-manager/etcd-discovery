/opt/etcd-v3.2.18-linux-amd64/etcd --name infra1 \
  --initial-advertise-peer-urls https://127.0.0.1:2380 \
  --listen-peer-urls https://127.0.0.1:2380 \
  --listen-client-urls https://127.0.0.1:2379 \
  --advertise-client-urls https://127.0.0.1:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra1=https://127.0.0.1:2380 \
  --initial-cluster-state new \
  --data-dir=/tmp/infra1 \
  --cert-file=etcd.local.config/certificates/db-server.crt \
  --key-file=etcd.local.config/certificates/db-server.key \
  --trusted-ca-file=etcd.local.config/certificates/db-ca.crt \
  --client-cert-auth=true \
  --peer-cert-file=etcd.local.config/certificates/peer-localhost.crt \
  --peer-key-file=etcd.local.config/certificates/peer-localhost.key \
  --peer-trusted-ca-file=etcd.local.config/certificates/peer-ca.crt \
  --peer-client-cert-auth=true


/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member list
e92d66acd89ecf29: name=infra1 peerURLs=https://127.0.0.1:2380 clientURLs=https://127.0.0.1:2379 isLeader=true

# ---------------------------------------------------------------------
/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member add infra2 https://127.0.0.2:2380
Added member named infra2 with ID 30439ed47a326d0 to cluster

ETCD_NAME="infra2"
ETCD_INITIAL_CLUSTER="infra2=https://127.0.0.2:2380,infra1=https://127.0.0.1:2380"
ETCD_INITIAL_CLUSTER_STATE="existing"

/opt/etcd-v3.2.18-linux-amd64/etcd --name infra2 \
  --initial-advertise-peer-urls https://127.0.0.2:2380 \
  --listen-peer-urls https://127.0.0.2:2380 \
  --listen-client-urls https://127.0.0.2:2379 \
  --advertise-client-urls https://127.0.0.2:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra1=https://127.0.0.1:2380,infra2=https://127.0.0.2:2380 \
  --initial-cluster-state existing \
  --data-dir=/tmp/infra2 \
  --cert-file=etcd.local.config/certificates/db-server.crt \
  --key-file=etcd.local.config/certificates/db-server.key \
  --trusted-ca-file=etcd.local.config/certificates/db-ca.crt \
  --client-cert-auth=true \
  --peer-cert-file=etcd.local.config/certificates/peer-localhost.crt \
  --peer-key-file=etcd.local.config/certificates/peer-localhost.key \
  --peer-trusted-ca-file=etcd.local.config/certificates/peer-ca.crt \
  --peer-client-cert-auth=true

/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member list
e75025783d78419: name=infra2 peerURLs=https://127.0.0.2:2380 clientURLs=https://127.0.0.2:2379 isLeader=false
e92d66acd89ecf29: name=infra1 peerURLs=https://127.0.0.1:2380 clientURLs=https://127.0.0.1:2379 isLeader=true

/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member add infra2 https://127.0.0.2:2380
membership: peerURL exists

# ---------------------------------------------------------------------
/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member add infra3 https://127.0.0.3:2380
Added member named infra3 with ID af695637f1c2a94c to cluster

ETCD_NAME="infra3"
ETCD_INITIAL_CLUSTER="infra2=https://127.0.0.2:2380,infra3=https://127.0.0.3:2380,infra1=https://127.0.0.1:2380"
ETCD_INITIAL_CLUSTER_STATE="existing"

/opt/etcd-v3.2.18-linux-amd64/etcd --name infra3 \
  --initial-advertise-peer-urls https://127.0.0.3:2380 \
  --listen-peer-urls https://127.0.0.3:2380 \
  --listen-client-urls https://127.0.0.3:2379 \
  --advertise-client-urls https://127.0.0.3:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra1=https://127.0.0.1:2380,infra2=https://127.0.0.2:2380,infra3=https://127.0.0.3:2380 \
  --initial-cluster-state existing \
  --data-dir=/tmp/infra3 \
  --cert-file=etcd.local.config/certificates/db-server.crt \
  --key-file=etcd.local.config/certificates/db-server.key \
  --trusted-ca-file=etcd.local.config/certificates/db-ca.crt \
  --client-cert-auth=true \
  --peer-cert-file=etcd.local.config/certificates/peer-localhost.crt \
  --peer-key-file=etcd.local.config/certificates/peer-localhost.key \
  --peer-trusted-ca-file=etcd.local.config/certificates/peer-ca.crt \
  --peer-client-cert-auth=true

/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member add infra3 https://127.0.0.3:2380
client: etcd cluster is unavailable or misconfigured; error #0: client: etcd member https://127.0.0.3:2379 has no leader
; error #1: client: etcd member https://127.0.0.2:2379 has no leader
; error #2: client: etcd member https://127.0.0.1:2379 has no leader

/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --endpoints=https://127.0.0.1:2379 \
  member add infra3 https://127.0.0.3:2380
membership: peerURL exists

/opt/etcd-v3.2.18-linux-amd64/etcdctl \
  --ca-file=etcd.local.config/certificates/db-ca.crt \
  --cert-file=etcd.local.config/certificates/db-discovery-client.crt \
  --key-file=etcd.local.config/certificates/db-discovery-client.key \
  --endpoints=https://127.0.0.1:2379 \
  member list
e75025783d78419: name=infra2 peerURLs=https://127.0.0.2:2380 clientURLs=https://127.0.0.2:2379 isLeader=false
af695637f1c2a94c: name=infra3 peerURLs=https://127.0.0.3:2380 clientURLs=https://127.0.0.3:2379 isLeader=false
e92d66acd89ecf29: name=infra1 peerURLs=https://127.0.0.1:2380 clientURLs=https://127.0.0.1:2379 isLeader=true

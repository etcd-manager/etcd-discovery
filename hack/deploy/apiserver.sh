#!/bin/bash
set -eou pipefail

echo "checking kubeconfig context"
kubectl config current-context || { echo "Set a context (kubectl use-context <context>) out of the following:"; echo; kubectl config get-contexts; exit 1; }
echo ""

# ref: https://stackoverflow.com/a/27776822/244009
case "$(uname -s)" in
    Darwin)
        curl -fsSL -o onessl https://github.com/pharmer/onessl/releases/download/0.1.0/onessl-darwin-amd64
        chmod +x onessl
        export ONESSL=./onessl
        ;;

    Linux)
        curl -fsSL -o onessl https://github.com/pharmer/onessl/releases/download/0.1.0/onessl-linux-amd64
        chmod +x onessl
        export ONESSL=./onessl
        ;;

    CYGWIN*|MINGW32*|MSYS*)
        curl -fsSL -o onessl.exe https://github.com/pharmer/onessl/releases/download/0.1.0/onessl-windows-amd64.exe
        chmod +x onessl.exe
        export ONESSL=./onessl.exe
        ;;
    *)
        echo 'other OS'
        ;;
esac

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup {
    rm -rf $ONESSL ca.crt ca.key server.crt server.key
}
trap cleanup EXIT

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export ETCD_NAMESPACE=kube-system
export ETCD_SERVICE_ACCOUNT=default
export ETCD_ENABLE_RBAC=false
export ETCD_RUN_ON_MASTER=0
export ETCD_DOCKER_REGISTRY=pharmer
export ETCD_IMAGE_PULL_SECRET=
export ETCD_UNINSTALL=0

show_help() {
    echo "etcd-discovery.sh - install etcd discovery server"
    echo " "
    echo "etcd-discovery.sh [options]"
    echo " "
    echo "options:"
    echo "-h, --help                         show brief help"
    echo "-n, --namespace=NAMESPACE          specify namespace (default: kube-system)"
    echo "    --rbac                         create RBAC roles and bindings"
    echo "    --docker-registry              docker registry used to pull etcd discovery server images (default: pharmer)"
    echo "    --image-pull-secret            name of secret used to pull etcd discovery server images"
    echo "    --run-on-master                run etcd discovery server on master"
    echo "    --uninstall                    uninstall etcd discovery server"
}

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        -n)
            shift
            if test $# -gt 0; then
                export ETCD_NAMESPACE=$1
            else
                echo "no namespace specified"
                exit 1
            fi
            shift
            ;;
        --namespace*)
            export ETCD_NAMESPACE=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --rbac)
            export ETCD_SERVICE_ACCOUNT=etcd-discovery
            export ETCD_ENABLE_RBAC=true
            shift
            ;;
        --docker-registry*)
            export ETCD_DOCKER_REGISTRY=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --image-pull-secret*)
            secret=`echo $1 | sed -e 's/^[^=]*=//g'`
            export ETCD_IMAGE_PULL_SECRET="name: '$secret'"
            shift
            ;;
        --run-on-master)
            export ETCD_RUN_ON_MASTER=1
            shift
            ;;
        --uninstall)
            export ETCD_UNINSTALL=1
            shift
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done

if [ "$ETCD_UNINSTALL" -eq 1 ]; then
    kubectl delete deployment -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete service -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete secret -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete validatingwebhookconfiguration -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete mutatingwebhookconfiguration -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete apiservice -l app=pharmer --namespace $ETCD_NAMESPACE
    # Delete RBAC objects, if --rbac flag was used.
    kubectl delete serviceaccount -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete clusterrolebindings -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete clusterrole -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete rolebindings -l app=pharmer --namespace $ETCD_NAMESPACE
    kubectl delete role -l app=pharmer --namespace $ETCD_NAMESPACE

    exit 0
fi

env | sort | grep ETCD*
echo ""

# create necessary TLS certificates:
# - a local CA key and cert
# - a apiserver key and cert signed by the local CA
$ONESSL create ca-cert
$ONESSL create server-cert server --domains=etcd-discovery.$ETCD_NAMESPACE.svc
export SERVICE_SERVING_CERT_CA=$(cat ca.crt | $ONESSL base64)
export TLS_SERVING_CERT=$(cat server.crt | $ONESSL base64)
export TLS_SERVING_KEY=$(cat server.key | $ONESSL base64)
export KUBE_CA=$($ONESSL get kube-ca | $ONESSL base64)
rm -rf $ONESSL ca.crt ca.key server.crt server.key

curl -fsSL https://raw.githubusercontent.com/etcd-manager/etcd-discovery/master/hack/deploy/operator.yaml | envsubst | kubectl apply -f -

if [ "$ETCD_ENABLE_RBAC" = true ]; then
    kubectl create serviceaccount $ETCD_SERVICE_ACCOUNT --namespace $ETCD_NAMESPACE
    kubectl label serviceaccount $ETCD_SERVICE_ACCOUNT app=pharmer --namespace $ETCD_NAMESPACE
    curl -fsSL https://raw.githubusercontent.com/etcd-manager/etcd-discovery/master/hack/deploy/rbac-list.yaml | envsubst | kubectl auth reconcile -f -
    curl -fsSL https://raw.githubusercontent.com/etcd-manager/etcd-discovery/master/hack/deploy/user-roles.yaml | envsubst | kubectl auth reconcile -f -
fi

if [ "$ETCD_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy etcd-discovery -n $ETCD_NAMESPACE \
      --patch="$(curl -fsSL https://raw.githubusercontent.com/etcd-manager/etcd-discovery/master/hack/deploy/run-on-master.yaml)"
fi

#!/bin/bash -x

set -euo pipefail

CLUSTERNAME="test-argocd-plugin"

# kind delete cluster --name ${CLUSTERNAME} --quiet 
# kind create cluster --name ${CLUSTERNAME} 
minikube delete -p ${CLUSTERNAME}
minikube start -p ${CLUSTERNAME}

kubectl create --context ${CLUSTERNAME} namespace argocd
kubectl apply --context ${CLUSTERNAME} -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl patch --context ${CLUSTERNAME} svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'

GOOS=linux GOARCH=amd64 go build -o tomlrocks main.go

kubectl wait --timeout=10m -n argocd --for=condition=available deployment/argocd-server
container_id=`kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o jsonpath='{.items..metadata.name}'`
repo_id=`kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-repo-server -o jsonpath='{.items..metadata.name}'`
kubectl cp ./tomlrocks argocd/${repo_id}:tomlrocks
kubectl -n argocd patch cm/argocd-cm -p "$(cat argocd-plugin.yaml)"

echo ${container_id} | pbcopy
while [ 1 ]
do
  service_url=`minikube service -p ${CLUSTERNAME} -n argocd argocd-server --url` || echo "failed to retrieve url"
  echo ${service_url} | grep -q "http://" || continue
  url=`echo ${service_url} |head -1| awk '{print $1}'`
  break
done


argocd login --insecure --username admin --password ${container_id} ${service_url##*/}
open ${service_url}

argocd app create testapp --config-management-plugin tomlrocks --repo https://github.com/sledigabel/test-argocd-plugin.git --dest-namespace default --dest-server  https://kubernetes.default.svc --path path --sync-option Prune=True

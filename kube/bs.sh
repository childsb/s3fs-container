#!/bin/sh
alias kubectl.sh=/Users/bc/dev/go-code/src/k8s.io/kubernetes/cluster/kubectl.sh

export KUBERNETES_PROVIDER=local

  kubectl.sh config set-cluster local --server=https://localhost:6443 --certificate-authority=/var/run/kubernetes/apiserver.crt
  kubectl.sh config set-credentials myself --username=admin --password=admin
  kubectl.sh config set-context local --cluster=local --user=myself
  kubectl.sh config use-context local

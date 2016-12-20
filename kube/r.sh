#!/bin/sh

alias kubectl.sh=/opt/go/src/k8s.io/kubernetes/cluster/kubectl.sh

 kubectl.sh create -f pv.yaml
 kubectl.sh create -f pvc.yaml
 kubectl.sh create -f pod.yaml

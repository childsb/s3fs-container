#!/bin/sh

alias kubectl.sh=/opt/go/src/k8s.io/kubernetes/cluster/kubectl.sh

 kubectl.sh delete -f pv.yaml
 kubectl.sh delete -f pvc.yaml
 kubectl.sh delete -f pod.yaml

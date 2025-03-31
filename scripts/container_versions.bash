#!/bin/bash

namespaces=(
  calico-apiserver
  calico-system
  cert-manager
  concourse
  external-secrets-operator
  gatekeeper-system
  ingress-controllers
  kube-system
  kuberos
  logging
  monitoring
  overprovision
  tigera-operator
  trivy-system
  velero
)

if [ -z "$1" ] || [[ ! " ${namespaces[@]} " =~ " $1 " ]]; then
  echo "Usage: $0 <namespace>"
  echo "Available namespaces: ${namespaces[@]}"
  exit 1
fi


remove_duplicates() {
  declare -a uniq_images=()
  for image in ${images[@]}; do
    if [[ ! " ${uniq_images[@]} " =~ " $image " ]]; then
      uniq_images+=($image)
    fi
  done
  images=("${uniq_images[@]}")
  
  for image in ${images[@]}; do
    echo $image
  done
}


declare -a containers=`kubectl get pods -n $1 -o json | jq -r '.items[] | .spec.containers[].name' | sort | uniq`
declare -a images=()

for container in ${containers[@]}; do
  images+=(`kubectl get pods -n $1 -o json | jq -r --arg NAME $container '.items[] | select(.spec.containers[].name == $NAME) | .spec.containers[] | select(.name == $NAME) | .image' | sort | uniq | while read line; do echo "$1-$container=$line"; done`)
done

remove_duplicates
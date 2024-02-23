---
name: Cloud Platform Kubernetes Upgrade story
about: Create new Cloud Platform Kubernetes Upgrade story
---

## Issue 1:
### Update vpc-cni from <current-version> to the version needed for current k8s version
Check for the latest supported addon version for the current kubernetes version and update the clusters

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)


## Issue 2:
### Update kube-proxy from <current-version> to the version needed for current k8s version
Check for the latest supported addon version for the current kubernetes version and update the clusters

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)


## Issue 3:
### Update core-dns from <current-version> to the version needed for current k8s version
Check for the latest supported addon version for the current kubernetes version and update the clusters

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)

## Issue 4:
### Test EKS <version> on test cluster

Follow instructions from Upgrade runbook: https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html#upgrade-eks-cluster

## Issue 5:
### Test EKS <version> on live-like cluster

Follow instructions from creating a live like cluster runbook: https://runbooks.cloud-platform.service.justice.gov.uk/creating-a-live-like.html#creating-a-live-like-cluster

## Issue 6:
### EKS: Upgrade Production clusters to Kubernetes <version>
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

## Issue 7:
### Update vpc-cni from <version> to the version needed for the upgraded k8s version
Check for the latest supported addon version for the upgraded kubernetes version and update the clusters

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)


## Issue 8:
### Update kube-proxy from <version> to the version needed for the upgraded k8s version
Check for the latest supported addon version for the upgraded kubernetes version and update the clusters

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)


## Issue 9:
### Update core-dns from <version> to the version needed for the upgraded k8s version
Check for the latest supported addon version for the upgraded kubernetes version and update the clusters

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)

## Issue 10:
### Review cluster components for upgrading
https://runbooks.cloud-platform.service.justice.gov.uk/container-images.html#container-images-used-by-cluster-components

Review the compatibility matrix for all cluster components and verify if the image is compatible with the upgraded kuebrnetes version

## Issue 11:
### Review and upgrade kube-state-metrics for the upgraded cluster version <version>
https://runbooks.cloud-platform.service.justice.gov.uk/container-images.html#container-images-used-by-cluster-components

https://github.com/kubernetes/kube-state-metrics?tab=readme-ov-file#compatibility-matrix

Review the compatibility matrix of kube-state-metrics and verify if the image is compatible with the upgraded kuebrnetes version. If the version is mismatch, check if the whole kube-prometheus-chart needs updating to get the default version of kube-state-metrics from the chart-> values. Otherwise, pin the kube-state-metrics image to match the kubernetes-version

## Issue 12:
### Upgrade cluster autoscalar for k8s version <version>
The Cloud Platform Cluster is in k8s version <version>. Hence upgrade the cluster-autoscalar to match the k8s version.

https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler
https://github.com/ministryofjustice/cloud-platform-terraform-cluster-autoscaler

## Issue 13:
### Upgrade cluster descheudler for k8s version <version>
The Cloud Platform Cluster is in k8s version <version>. Hence upgrade the descheduler to match the k8s version.

https://github.com/kubernetes-sigs/descheduler?tab=readme-ov-file#%EF%B8%8F--documentation-versions-by-release
https://github.com/ministryofjustice/cloud-platform-terraform-descheduler

## Issue 14:
### Post k8s version <version> Cleanup
Following upgrade from EKS <old-version > <version>, there will be a number of cleanup activities that need addressing:

Update user guide / runbook references
Update tools-image for kubectl version
Update cloud-platform-cli for kubectl version 
Update concourse pipelines to use updated tools-image and cli
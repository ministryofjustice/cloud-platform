---
name: Cloud Platform Kubernetes Upgrade story
about: Create new Cloud Platform Kubernetes Upgrade story
---

## Issue 1:
### Terraform AWS EKS Modules support for <version>
Check to see if there are any issues with using <current EKS module> with Kubernetes <version>

## Issue 2:
### Check for deprecated resources for <vesion>
Scan to see if we have any in the production clusters live/live-2/manager. If any present, take actions to address that. 

Refer to the deprecations note on runbook: https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html#compatibility-check

## Issue 3:
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


## Issue 4:
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


## Issue 5:
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

## Issue 6:
### Test EKS <version> on test cluster

Follow instructions from Upgrade runbook: https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html#upgrade-eks-cluster

## Issue 7:
### Test EKS <version> on live-like cluster

Follow instructions from creating a live like cluster runbook: https://runbooks.cloud-platform.service.justice.gov.uk/creating-a-live-like.html#creating-a-live-like-cluster

## Issue 8:
### EKS: Upgrade live-2 to Kubernetes <version>
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

## Issue 9:
### EKS: Upgrade manager to Kubernetes <version>
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

## Issue 10:
### EKS: Upgrade live to Kubernetes <version>
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

## Issue 11:
### EKS: Upgrade live to Kubernetes <version>
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

## Issue 12:
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


## Issue 13:
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


## Issue 14:
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

## Issue 15:
### Review cluster components for upgrading
https://runbooks.cloud-platform.service.justice.gov.uk/container-images.html#container-images-used-by-cluster-components

Review the compatibility matrix for all cluster components and verify if the image is compatible with the upgraded kuebrnetes version

## Issue 16:
### Review and upgrade kube-state-metrics for the upgraded cluster version <version>
https://runbooks.cloud-platform.service.justice.gov.uk/container-images.html#container-images-used-by-cluster-components

https://github.com/kubernetes/kube-state-metrics?tab=readme-ov-file#compatibility-matrix

Review the compatibility matrix of kube-state-metrics and verify if the image is compatible with the upgraded kuebrnetes version. If the version is mismatch, check if the whole kube-prometheus-chart needs updating to get the default version of kube-state-metrics from the chart-> values. Otherwise, pin the kube-state-metrics image to match the kubernetes-version

## Issue 17:
### Upgrade cluster autoscalar for k8s version <version>
The Cloud Platform Cluster is in k8s version <version>. Hence upgrade the cluster-autoscalar to match the k8s version.

https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler
https://github.com/ministryofjustice/cloud-platform-terraform-cluster-autoscaler

## Issue 18:
### Upgrade cluster descheudler for k8s version <version>
The Cloud Platform Cluster is in k8s version <version>. Hence upgrade the descheduler to match the k8s version.

https://github.com/kubernetes-sigs/descheduler?tab=readme-ov-file#%EF%B8%8F--documentation-versions-by-release
https://github.com/ministryofjustice/cloud-platform-terraform-descheduler

## Issue 19:
### Post k8s version <version> Cleanup
Following upgrade from EKS <old-version > <version>, there will be a number of cleanup activities that need addressing:

Update user guide / runbook references
Update tools-image for kubectl version
Update cloud-platform-cli for kubectl version 
Update concourse pipelines to use updated tools-image and cli
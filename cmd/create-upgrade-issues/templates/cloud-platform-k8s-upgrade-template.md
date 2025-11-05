## Issue 1:
### Planning upgrade to EKS <upgrade-version>

Go through the EKS and Kubernetes release notes for version <upgrade-version> and create a plan to upgrade our clusters

Things to consider:

Review EKS / Kubernetes changelogs & release notes
EKS Module supported at target upgrade version?
Are there any API deprecations & removals? (Check EKS insights)
Are there new components being added?
What changes are being introduced to current components?
Are there changes to core infra of the CP required? i.e. Are all our current components compatible with <upgrade-version>?
Are there changes users need to make?
Do we need to expand any of our smoke/integration testing?
Create additional tickets needed for any findings specific to this upgrade

Cluster upgrade Runbook:
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

## Issue 2:
### Update vpc-cni to latest supported release for current EKS version
Check for the latest supported addon release for the current EKS version and update the clusters
[Runbook link](https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-addons.html#listing-available-eks-upgrades)

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)


## Issue 3:
### Update kube-proxy to latest supported release for current EKS version
Check for the latest supported addon release for the current EKS version and update the clusters
[Runbook link](https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-addons.html#listing-available-eks-upgrades)

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)


## Issue 4:
### Update core-dns to latest supported release for current EKS version
Check for the latest supported addon release for the current EKS version and update the clusters
[Runbook link](https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-addons.html#listing-available-eks-upgrades)

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)

## Issue 5:
### Test EKS <upgrade-version> on test cluster

Follow instructions from Upgrade runbook: https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html#upgrade-eks-cluster

## Issue 6:
### Test EKS <upgrade-version> on live-like cluster

Follow instructions from creating a live like cluster runbook: https://runbooks.cloud-platform.service.justice.gov.uk/creating-a-live-like.html#creating-a-live-like-cluster

## Issue 7:
### EKS: Upgrade Production clusters to Kubernetes <upgrade-version>
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

## Issue 8:
### Update vpc-cni to latest supported release for upgraded EKS version
Check for the latest supported addon version for the upgraded EKS version and update the clusters
[Runbook link](https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-addons.html#listing-available-eks-upgrades)

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)


## Issue 9:
### Update kube-proxy to latest supported release for upgraded EKS version
Check for the latest supported addon version for the upgraded EKS version and update the clusters
[Runbook link](https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-addons.html#listing-available-eks-upgrades)

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)


## Issue 10:
### Update core-dns to latest supported release for upgraded EKS version
Check for the latest supported addon version for the upgraded EKS version and update the clusters
[Runbook link](https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-addons.html#listing-available-eks-upgrades)

Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)

## Issue 11:
### Review cluster components for upgrading
https://runbooks.cloud-platform.service.justice.gov.uk/container-images.html#container-images-used-by-cluster-components

Review the compatibility matrix for all cluster components and verify if the image is compatible with the upgraded Kubernetes version

## Issue 12:
### Review and upgrade kube-state-metrics for the upgraded EKS version <upgrade-version>
https://runbooks.cloud-platform.service.justice.gov.uk/container-images.html#container-images-used-by-cluster-components

https://github.com/kubernetes/kube-state-metrics?tab=readme-ov-file#compatibility-matrix

Review the compatibility matrix of kube-state-metrics and verify if the image is compatible with the upgraded Kubernetes version. If the version is a mismatch, check if the whole kube-prometheus-chart needs updating to get the default version of kube-state-metrics from the chart-> values. Otherwise, pin the kube-state-metrics image to match the Kubernetes version.

## Issue 13:
### Upgrade cluster autoscaler for EKS version <upgrade-version>
The Cloud Platform EKS Cluster is in Kubernetes version <upgrade-version>. Hence upgrade the cluster-autoscaler to match the Kubernetes version.

https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler
https://github.com/ministryofjustice/cloud-platform-terraform-cluster-autoscaler

## Issue 14:
### Upgrade cluster descheduler for EKS version <upgrade-version>
The Cloud Platform EKS Cluster is in Kubernetes version <upgrade-version>. Hence upgrade the descheduler to match the Kubernetes version.

https://github.com/kubernetes-sigs/descheduler?tab=readme-ov-file#%EF%B8%8F--documentation-versions-by-release
https://github.com/ministryofjustice/cloud-platform-terraform-descheduler

## Issue 15:
### Post EKS version <upgrade-version> Cleanup
Following upgrade to EKS <upgrade-version>, there will be a number of cleanup activities that need addressing:

Update user guide / runbook references to EKS version
Update tools-image for kubectl version
Update cloud-platform-cli for kubectl version
Update concourse pipelines to use updated tools-image and cli

## Issue 16:
### Add deprecated apis from <upgrade-version> to Gatekeeper
Following upgrade to EKS <upgrade-version>, create a new [gatekeeper constraint](https://github.com/ministryofjustice/cloud-platform-terraform-gatekeeper/tree/main/resources/constraints) that stops people from using any apis that are now deprecated as a part of the <upgrade-version> upgrade.

## Issue 17:
### Update Upgrade runbook and Cluster Upgrade Issue Template

Update the runbook for
- any changes needed in the steps to perform the upgrade
- any lessons learnt that could be useful for next upgrade
- any changes to the [upgrade issue template](https://github.com/ministryofjustice/cloud-platform/blob/main/.github/ISSUE_TEMPLATE/cloud-platform-k8s-upgrade-template.md)

Cluster upgrade Runbook:
https://runbooks.cloud-platform.service.justice.gov.uk/upgrade-eks-cluster.html

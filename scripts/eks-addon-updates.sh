#!/bin/bash

set -ex

CLUSTER=live

# get cluster versions
CLUSTER_VERSION=$(aws eks describe-cluster --name "$CLUSTER" | jq -r '.cluster.version')

# get addons
CLUSTER_ADDONS=($(aws eks list-addons --cluster-name "$CLUSTER" | jq -r '.addons[] | .'))

for CLUSTER_ADDON in "${CLUSTER_ADDONS[@]}";
do
  # get addon version for cluster
  CLUSTER_ADDON_VERSION=$(aws eks describe-addon --cluster-name "$CLUSTER" --addon-name "$CLUSTER_ADDON" | jq -r '.addon.addonVersion')

  # get latest supported addon version for the cluster/k8s version
  LATEST_SUPPORTED_ADDON_VERSION_FOR_KUBERNETES_VERSION=$(aws eks describe-addon-versions --addon-name "$CLUSTER_ADDON" --kubernetes-version "$CLUSTER_VERSION" | jq -r '.addons[0].addonVersions[0].addonVersion')

  TITLE="EKS addon Production Clusters: Update $CLUSTER_ADDON from $CLUSTER_ADDON_VERSION to the latest version"

  if [[ "$LATEST_SUPPORTED_ADDON_VERSION_FOR_KUBERNETES_VERSION" != "$CLUSTER_ADDON_VERSION" ]]; then # check if newer version is supported
    BODY=$(cat << END
## Background

There is a new version of the EKS add-on $CLUSTER_ADDON. $CLUSTER_ADDON needs updating on all of our clusters. When this issue was created, the latest supported add-on version for Kubernetes $CLUSTER_VERSION was $LATEST_SUPPORTED_ADDON_VERSION_FOR_KUBERNETES_VERSION.

## Production Clusters Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See the [Amazon EKS add-ons](https://docs.aws.amazon.com/eks/latest/userguide/eks-add-ons.html) documentation for more information about addons, or find the latest versions for these EKS add-ons directly:

- [coredns](https://docs.aws.amazon.com/eks/latest/userguide/managing-coredns.html)
- [kube-proxy](https://docs.aws.amazon.com/eks/latest/userguide/managing-kube-proxy.html)
- [vpc-cni](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)
- [runtime-monitoring](https://docs.aws.amazon.com/guardduty/latest/ug/how-does-runtime-monitoring-work.html)
END
)

    GITHUB_ISSUES=$(gh issue list --repo ministryofjustice/cloud-platform --state all --search "in:title \"$TITLE\"" --limit 50 --json title | jq -r "[ .[] | select(.title == \"$TITLE\") ] | length")

    # if no issues yet, create one
    if (( $(echo "0 == $GITHUB_ISSUES" | bc -l) )); then
      echo "No issue found for $TITLE, creating one..."
      gh issue create --title "$TITLE" --body "$BODY" --repo ministryofjustice/cloud-platform
    else
      echo "Issue already exists for $TITLE, skipping..."
    fi
  else
    echo "Up to date, skipping issue creation for $TITLE"
  fi

done

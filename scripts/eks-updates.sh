#!/bin/bash

set -ex

# get supported versions
VERSIONS=($(aws eks describe-addon-versions | jq -r ".addons[] | .addonVersions[] | .compatibilities[] | .clusterVersion" | sort | uniq))

CLUSTER=live

# get cluster versions
CLUSTER_VERSION=$(aws eks describe-cluster --name "$CLUSTER" | jq -r '.cluster.version')

for VERSION in "${VERSIONS[@]}";
do
  if [[ "$CLUSTER_VERSION" != "$VERSION" ]]; then
    if (( $(echo "$CLUSTER_VERSION < $VERSION" | bc -l) )); then # check if newer version is supported
      TITLE="EKS: Upgrade $CLUSTER to Kubernetes v$VERSION";
      BODY=$(cat << END
## Background

EKS supports Kubernetes $VERSION. The $CLUSTER cluster needs upgrading to Kubernetes $VERSION.

## Environments Checklist:

- [ ] live-2
- [ ] manager
- [ ] live

See [Amazon EKS Kubernetes versions](https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html) for more details.
END
)

      # get github issues and check if one already exists
      GITHUB_ISSUES=$(gh issue list --repo ministryofjustice/cloud-platform --state all --search "in:title \"$TITLE\"" --limit 50 --json title | jq -r "[ .[] | select(.title == \"$TITLE\") ] | length")

      # if no issues yet, create one
      if (( $(echo "0 == $GITHUB_ISSUES" | bc -l) )); then
        echo "No issue found for $TITLE, creating one..."
        gh issue create --title "$TITLE" --body "$BODY" --label EPIC --repo ministryofjustice/cloud-platform
      else
        echo "Issue already exists for $TITLE, skipping..."
      fi
    fi
  fi
done

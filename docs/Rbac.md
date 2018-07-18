# How to install RBAC on a new cluster.

## Purpose
This document is to provide an instructions on applying the ClusterRole to the new cluster, this allows the webops team to have cluster wide permissions. 

**Steps:**
2. `$ kubectl apply -f kubernetes-investigations/cluster-config/rbac/webops-cluster-admin.yml`
3. `clusterrolebinding.rbac.authorization.k8s.io "webops-cluster-admin" created`










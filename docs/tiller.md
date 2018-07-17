# How to install Tiller on a new cluster

Note: Helm should be installed locally before running the below.

```
$ brew install kubernetes-helm
```

Steps:
1. Clone the `ministryofjustice/kubernetes-investigations` repo 
https://github.com/ministryofjustice/kubernetes-investigations
2. `$ kubectl apply -f ./cluster-components/helm/rbac-config.yml`
3. `$ helm init --tiller-namespace kube-system --service-account tiller`


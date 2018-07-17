# HOW TO INSTALL TILLER ON A NEW CLUSTER

Note: Helm should be installed locally before running the below.

```
$ brew install kubernetes-helm
```

Steps:
1. Clone the `kubernetes-investigations` repo 
2. `$ kubectl apply -f ./cluster-components/helm/rbac-config.yml`
3. `$ helm init --tiller-namespace kube-system --service-account tiller`


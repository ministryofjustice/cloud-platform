# One cluster for dev/staging/prod

Date: 01/06/2018

## Status

Accepted

## Context

One Cluster or Many?, which would be the best approach out of running a single large cluster or running multiple cluster per dev/staging/prod. One of the many strengths of Kubernetes is just how much flexibility you have when deploying and operating the orchestrator for your containerised workloads. You get to choose the number of nodes, pods, containers per cluster and a host of other parameters to fit your needs. But, as always, flexibility comes with responsibility. Question is do you have the right tools and manpower to effectively manage multiple clusters so that they don't deviate in long run.

Managing single cluster keep operational overhead and costs minimum and easy to maintain, can a single cluster able to satisfy the needs of multiple environments, users or groups of users with out compromising security. This can be achieved using Kubernetes namespace abstraction which enable multiple workloads to operate within the same Kubernetes cluster. Kubernetes namespaces help different projects, teams, or customers to share a Kubernetes cluster. A Kubernetes namespace provides the scope for Pods, Services, and Deployments in the cluster. Users interacting with one namespace do not see the content in another namespace.

Single cluster with different Namespace for dev/staging/prod provides a unique scope for named resources (to avoid basic naming collisions) , delegated management authority to trusted users and ability to limit community resource consumption. Environment(dev/staging/prod) isolation can be achieved by managing each environment in a separate namespace, we can collect metrics at namespace level and can set resources quota at namespace level. That way you can take action at environment level for dev/staging/prod.


## Decision

After consideration of pros and cons of each approach, we went with one cluster and multiple namespaces.

Some important reasons behind this move were:

A single cluster (with namespaces and RBAC) is easier to setup and manage. 

A single k8s cluster does support high load. 

Reduced operational overhead and reduced costs for managing a single cluster.

It works well for us as it lets us manage all the services more easily with one ingress controller, centralised logs collection via fluentd, centralised monitoring with Prometheus, cluster policies, etc. so for now it looks like a good approach.


## Consequences

Since we only want to be running a single cluster, when we create a new cluster and want to migrate to new cluster, we will need to shut down old cluster as soon as it's no longer needed. Also services migrate from old cluster to new cluster sooner will avoid the complexities of running two parallel clusters.

The selection of the number of Kubernetes clusters may be a relatively static choice, only revisited occasionally. By contrast, the number of nodes in a cluster and the number of pods in a service may change frequently according to load and growth.

Monitor the amount of resources each namespace consume in order to limit the impact to other namespaces using the cluster.

There are multiple cluster components that are shared across all tenants within a cluster, regardless of namespace. These shared components include the master components such as the API server, controller manager, scheduler, and DNS, as well as worker components such as the Kubelet and Kube Proxy. Sharing these non-namespace-aware components between tenants necessarily exposes tenant resources to all other tenants in the cluster or on the same worker node.
### Service Catalog

This is the kubernetes side of things.

#### Reference

- https://svc-cat.io/
- https://svc-cat.io/docs/install/

#### Installation

To install, switch to the context of the cluster.

```
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm install svc-cat/catalog --name catalog --namespace catalog
```

Fetch the `svcat` binary (not required but useful).

```
curl -sLO https://download.svcat.sh/cli/latest/darwin/amd64/svcat
chmod +x svcat
./svcat get bindings
```

### AWS Service Broker

The AWS implementation of a service broker.

#### Reference
- https://github.com/awslabs/aws-servicebroker
- https://github.com/awslabs/aws-servicebroker/blob/master/docs/getting-started-k8s.md

#### Installation

Setup IAM roles and user (name the workspace after the cluster you'll be working on)

```
terraform init
terraform workspace new cloud-platform-test-1
terraform apply
```

Fetch the AWS service broker release.

```
curl -kLO https://s3.amazonaws.com/awsservicebroker/assets/aws-service-broker-install.tar.gz
mkdir -p asb
tar -xvf aws-service-broker-install.tar.gz -C asb
rm aws-service-broker-install.tar.gz
```

Edit `asb/k8s-variable.yaml` and populate the values from the terraform outputs (refer to the getting started document).

Don't forget to change the region and set the VPC id for the cluster.

```
aws --output json ec2 describe-vpcs --filters=Name=tag:Name,Values="$(terraform workspace show)"  | jq -r '.Vpcs[0].VpcId'
```

Install the broker.

```
cd asb
bash ./install_aws_service_broker.sh
```

Request a resouce.

```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: demo
---
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceInstance
metadata:
  name: demo-sqs
  namespace: demo
spec:
  clusterServiceClassExternalName: dh-sqs
  clusterServicePlanExternalName: standard
---
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceBinding
metadata:
  name: demo-sqs
  namespace: demo
spec:
  instanceRef:
    name: demo-sqs
EOF
```

### Cleanup

Remove `ServiceBindings` and `ServiceInstances` (and wait for deprovisioning to finish).

```
cd asb
bash uninstall.sh
helm delete --purge catalog
kubectl delete ns demo catalog
```

Cleanup IAM resources.

```
terraform destroy
```

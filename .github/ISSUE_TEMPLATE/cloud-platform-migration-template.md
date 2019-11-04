---
name: Cloud Platform Migration
about: Track migrations to cloud platform

---

## Information
Service name : <!--- Full name or acronym -->

Production service URL : <!--- *.service.gov.uk -->

Migration Date : <!--- 25/12/2019 -->  

Contact Person : <!-- John Smith -->



## Environments
Below are the expected steps to migrate to the Cloud Platfom

### **Development environment**
#### Infrastructure Resources

- [ ] The Kubernetes Namespaces have been created
- [ ] The AWS resource (RDS/Elasticache/ECR/S3/etc.) have been created
- [ ] The Route53 zone has been created (if applicable)

#### Kubernetes Deployments

- [ ] The Ingresses have been created (if applicable)
- [ ] The Certificates have been created (if applicable)
- [ ] The applications are deploying successfully

#### Data Migration

- [ ] A Data validation strategy has been defined
- [ ] The RDS instances have been migrated successfully
- [ ] The S3 buckets have been migrated (if applicable)

---
### **Staging environment**
#### Infrastructure Resources

- [ ] The Kubernetes Namespaces have been created
- [ ] The AWS resource (RDS/Elasticache/ECR/S3/etc.) have been created
- [ ] The Route53 zone has been created (if applicable)

#### Kubernetes Deployments

- [ ] The Ingresses have been created (if applicable)
- [ ] The Certificates have been created (if applicable)
- [ ] The applications are deploying successfully

#### Data Migration

- [ ] A Data validation strategy has been defined
- [ ] The RDS instances have been migrated successfully
- [ ] The S3 buckets have been migrated (if applicable)  



---
### **Production environment**
#### Infrastructure Resources

- [ ] The Kubernetes Namespaces have been created
- [ ] The AWS resource (RDS/Elasticache/ECR/S3/etc.) have been created
- [ ] The Route53 zone has been created (if applicable)

#### Kubernetes Deployments

- [ ] The Ingresses have been created (if applicable)
- [ ] The Certificates have been created (if applicable)
- [ ] The applications are deploying successfully

#### Data Migration

- [ ] A Data validation strategy has been defined
- [ ] The RDS instances have been migrated successfully
- [ ] The S3 buckets have been migrated (if applicable)

#### DNS Switch

- [ ] The DNS switch has happened, the URL is now routing to the new service.

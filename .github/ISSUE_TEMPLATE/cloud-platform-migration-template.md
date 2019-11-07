---
name: Cloud Platform Migration
about: Track migrations to cloud platform

---

## Information
Service name : <!--- Full name and abbreviation (if any) -->

Production service URL : <!--- *.service.gov.uk -->

Planned Migration Date : <!--- Leave blank if no date has been decided-->  

Contact Person : <!-- email address -->



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
- [ ] The Route53 zone has been created

#### Kubernetes Deployments

- [ ] The Ingresses have been created
- [ ] The Certificates have been created
- [ ] The applications are deploying successfully

#### Data Migration

- [ ] A Data validation strategy has been defined
- [ ] The RDS instances have been migrated successfully
- [ ] The S3 buckets have been migrated (if applicable)

#### DNS Switch

- [ ] DNS authority has been delegated to the cloud platform Route53 zone (if applicable)
- [ ] The DNS switch has happened, the URL is now routing to the new service.

# CIDR allocation register

| CIDR          | Mask | Allocated to                         | Notes                               |
| :------------ | :--- | :----------------------------------- | :---------------------------------- |
| 10.195.0.0/16 | /16  | Non-live + platform infrastructure   | Routable                            |
| 10.41.0.0/16  | /16  | Live                                 | Routable                            |
| 100.64.0.0/10 | /10  | Pod IPs                              | RFC6598, not externally routable    |
| 172.20.0.0/16 | /16  | Existing Cloud Platform Cluster      | Routable                            |

## Platform and Shared Infrastructure (10.195.0.0/16, lower range)

| VPC                            | CIDR          | Mask | Addresses | Purpose                                           |
| :----------------------------- | :------------ | :--- | --------: | :------------------------------------------------ |
| Hub cluster (Cloud Platform) | 10.195.0.0/20 | /20 | 4,096 | Argo CD, KRO, OpenSearch — ~250 pods, few nodes |
| Hub preproduction | 10.195.16.0/20 | /20 | 4,096 | Preproduction Hub |
| Dev clusters (shared) | 10.195.32.0/20 | /20 | 4,096 | Up to 10 CP engineer dev clusters sharing one VPC |

### BU Non-Live Clusters (10.195.0.0/16, upper range)

| VPC           | CIDR            | Mask | Addresses | Assigned BU |
| :------------ | :-------------- | :--- | --------: | :---------- |
| BU-1 non-live | 10.195.48.0/20 | /20 | 4,096 | OCTO |
| BU-2 non-live | 10.195.64.0/20 | /20 | 4,096 | HMPPS |
| BU-3 non-live | 10.195.80.0/20 | /20 | 4,096 | LAA |
| BU-4 non-live | 10.195.96.0/20 | /20 | 4,096 | HQ |
| BU-5 non-live | 10.195.112.0/20 | /20 | 4,096 | HMCTS |
| BU-6 non-live | 10.195.128.0/20 | /20 | 4,096 | CICA |
| BU-7 non-live | 10.195.144.0/20 | /20 | 4,096 | OPG |
| BU-8 non-live | 10.195.160.0/20 | /20 | 4,096 | Unassigned |
| BU-9 non-live | 10.195.176.0/20 | /20 | 4,096 | Unassigned |
| Free | 10.195.192.0/18 | /18 | 16,384 | Unassigned |

### BU Live Clusters (10.41.0.0/16)

| VPC      | CIDR          | Size | Addresses | Assigned BU |
| :------- | :------------ | :--- | --------: | :---------- |
| BU-1 live | 10.41.0.0/20 | /20 | 4,096 | OCTO |
| BU-2 live | 10.41.16.0/20 | /20 | 4,096 | HMPPS |
| BU-3 live | 10.41.32.0/20 | /20 | 4,096 | LAA |
| BU-4 live | 10.41.48.0/20 | /20 | 4,096 | HQ |
| BU-5 live | 10.41.64.0/20 | /20 | 4,096 | HMCTS |
| BU-6 live | 10.41.80.0/20 | /20 | 4,096 | CICA |
| BU-7 live | 10.41.96.0/20 | /20 | 4,096 | OPG |
| BU-8 live | 10.41.112.0/20 | /20 | 4,096 | Unassigned |
| BU-9 live | 10.41.128.0/20 | /20 | 4,096 | Unassigned |
| BU-10 live | 10.41.144.0/20 | /20 | 4,096 | Unassigned |
| Free | 10.41.160.0/19 | /19 | 8,192 | Unassigned |
| Free | 10.41.192.0/18 | /18 | 16,384 | Unassigned |

## Secondary CIDR Allocation (Pod IPs - 100.64.0.0/10)

Each EKS cluster receives a /16 from the RFC6598 range for pod IPs.

| Cluster                  | Secondary CIDR | Pod IPs  | Notes                                              |
| :-------------------- | :------------- | -------: | :------------------------------------------------- |
| Hub cluster | 100.64.0.0/16 | 65,536 | ~250 pods; included for configuration consistency. |
| Hub preproduction | 100.65.0.0/16 | 65,536 | ~250 pods; included for configuration consistency. |
| Dev clusters (shared) | 100.66.0.0/16 | 65,536 | Shared across up to 10 dev clusters |
| Free | 100.67.0.0/16 | 65,536 | Kept free for future platform use |
| BU-1 non-live | 100.68.0.0/16 | 65,536 | OCTO |
| BU-2 non-live | 100.69.0.0/16 | 65,536 | HMPPS |
| BU-3 non-live | 100.70.0.0/16 | 65,536 | LAA |
| BU-4 non-live | 100.71.0.0/16 | 65,536 | HQ |
| BU-5 non-live | 100.72.0.0/16 | 65,536 | HMCTS |
| BU-6 non-live | 100.73.0.0/16 | 65,536 | CICA |
| BU-7 non-live | 100.74.0.0/16 | 65,536 | OPG |
| BU-8 non-live | 100.75.0.0/16 | 65,536 | Unassigned |
| BU-9 non-live | 100.76.0.0/16 | 65,536 | Unassigned |
| Free | 100.77.0.0/16 | 65,536 | Unassigned |
| Free | 100.78.0.0/16 | 65,536 | Unassigned |
| Free | 100.79.0.0/16 | 65,536 | Unassigned |
| BU-1 live | 100.80.0.0/16 | 65,536 | OCTO |
| BU-2 live | 100.81.0.0/16 | 65,536 | HMPPS |
| BU-3 live | 100.82.0.0/16 | 65,536 | LAA |
| BU-4 live | 100.83.0.0/16 | 65,536 | HQ |
| BU-5 live | 100.84.0.0/16 | 65,536 | HMCTS |
| BU-6 live | 100.85.0.0/16 | 65,536 | CICA |
| BU-7 live | 100.86.0.0/16 | 65,536 | OPG |
| BU-8 live | 100.87.0.0/16 | 65,536 | Unassigned |
| BU-9 live | 100.88.0.0/16 | 65,536 | Unassigned |
| BU-10 live | 100.89.0.0/16 | 65,536 | Unassigned |
| Free | 100.90.0.0/15 | 131,072 | Unassigned |
| Free | 100.92.0.0/14 | 262,144 | Unassigned |
| Free | 100.96.0.0/11 | 2,097,152 | Unassigned |

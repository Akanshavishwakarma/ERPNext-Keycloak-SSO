\# Cloud Deployment



\## Overview



This project was deployed as a Proof of Concept on an Azure Ubuntu Virtual Machine.



The deployment included:



\- Azure infrastructure

\- ERPNext / Frappe

\- MariaDB

\- Redis

\- Nginx

\- Supervisor

\- Docker

\- MySQL

\- Keycloak

\- OpenID Connect (OIDC)



\## Azure Infrastructure



\### Resource Group



`rg-av`



\### Region



Central India



\### Virtual Network



`vnet-av`



Address space:



`10.0.0.0/16`



\### Subnet



`subnet-av`



Address space:



`10.0.0.0/24`



\### Virtual Machine



Name:



`VM-av`



Operating System:



Ubuntu Server 24.04 LTS



VM Size:



Standard D2ls\_v6



Resources:



\- 2 vCPU

\- 4 GiB RAM



The Azure infrastructure was created manually through the Azure Portal.



\## ERPNext Deployment



ERPNext was deployed on the Ubuntu VM using Frappe Bench.



\### ERPNext Components



\- Frappe 15.121.3

\- ERPNext 15.121.6

\- MariaDB 10.11.14

\- Redis

\- Nginx

\- Supervisor

\- Frappe Bench



\### ERPNext Site



`erp.av.local`



The ERPNext application, database, Redis services, Nginx and Supervisor-managed services were configured on the VM.



\## Keycloak Deployment



Keycloak was deployed using Docker.



\### Keycloak Stack



\- Keycloak 26.7.4

\- MySQL 8.0

\- Docker

\- Docker Compose



Keycloak and MySQL were connected through a dedicated Docker network.



A persistent Docker volume was used for MySQL data.



Keycloak was exposed on port:



`8080`



\## Keycloak Configuration



A Keycloak realm named:



`erpnext`



was configured for the SSO Proof of Concept.



An OIDC client named:



`erpnext`



was created for ERPNext integration.



The client used the standard Authorization Code flow.



\## ERPNext + Keycloak Integration



ERPNext was integrated with Keycloak using OpenID Connect.



The authentication flow was:



```text

User

&#x20; |

&#x20; v

ERPNext

&#x20; |

&#x20; v

Keycloak

&#x20; |

&#x20; | OIDC Authentication

&#x20; |

&#x20; v

ERPNext

&#x20; |

&#x20; v

Dashboard


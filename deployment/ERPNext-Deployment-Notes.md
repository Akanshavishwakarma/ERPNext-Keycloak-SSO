\# ERPNext Deployment



\## Overview



ERPNext was deployed on an Azure Ubuntu VM as part of the ERPNext + Keycloak SSO Proof of Concept.



\## Azure Environment



\- Resource Group: `rg-av`

\- Region: Central India

\- Virtual Network: `vnet-av`

\- VNet Address Space: `10.0.0.0/16`

\- Subnet: `subnet-av`

\- Subnet Address Space: `10.0.0.0/24`

\- VM: `VM-av`

\- OS: Ubuntu Server 24.04 LTS

\- VM Size: Standard D2ls\_v6

\- CPU: 2 vCPU

\- Memory: 4 GiB



\## ERPNext Stack



\- Frappe: 15.121.3

\- ERPNext: 15.121.6

\- MariaDB: 10.11.14

\- Redis

\- Nginx

\- Supervisor

\- Frappe Bench



\## ERPNext Site



Site name:



`erp.av.local`



The ERPNext installation was managed through Frappe Bench.



\## Services



The deployment used:



\- MariaDB for database storage

\- Redis for caching and background jobs

\- Nginx as the web server/reverse proxy

\- Supervisor for managing ERPNext services

\- Frappe Bench for ERPNext application management



\## Deployment Validation



The ERPNext services were verified using Supervisor.



The following components were running successfully:



\- Redis cache

\- Redis queue

\- ERPNext web service

\- Socket.IO

\- Background workers



Nginx was also configured to serve the ERPNext application.



\## Issue Resolved



An ERPNext asset loading issue resulted in HTTP 404 responses.



The issue was resolved by correcting directory traversal permission on the application user's home directory and reloading Nginx.



\## Keycloak Integration



After ERPNext deployment, Keycloak was integrated with ERPNext using OpenID Connect (OIDC).



The authentication flow was:



ERPNext Login  

→ Keycloak  

→ OIDC Authentication  

→ ERPNext Dashboard



The end-to-end SSO flow was successfully tested.



\## Security Note



Passwords, private keys, client secrets and other credentials are intentionally excluded from this repository.


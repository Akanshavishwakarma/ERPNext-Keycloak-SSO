\# Bibo Desktop - Keycloak Integration



\## Overview



Bibo Desktop is integrated with the ERPNext + Keycloak SSO flow so that users do not need to maintain a separate Bibo login.



\## Authentication Flow



User

→ ERPNext

→ Keycloak

→ ERPNext Dashboard

→ Bibo Desktop Activation

→ Bibo Backend

→ Bibo Desktop Session



\## How It Works



1\. The user logs in to ERPNext through Keycloak.

2\. ERPNext authenticates the user using Keycloak.

3\. The authenticated user starts Bibo Desktop.

4\. Bibo uses the Keycloak-authenticated identity for desktop activation.

5\. A temporary one-time activation code is generated.

6\. The Bibo Desktop application sends the activation code to the Bibo backend.

7\. The backend validates the activation code.

8\. Bibo receives its access and refresh tokens.

9\. The desktop application stores the authenticated session.



\## Backend Activation Endpoints



\### Create Desktop Activation



```text

POST /auth/desktop/activation


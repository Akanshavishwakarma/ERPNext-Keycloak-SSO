\# ERPNext + Keycloak SSO Integration



\## Overview



Keycloak was integrated with ERPNext as the Identity Provider (IdP) using OpenID Connect (OIDC).



The purpose of the integration was to allow users to authenticate through Keycloak and access ERPNext using Single Sign-On (SSO).



\## Architecture



User

&#x20; |

&#x20; v

ERPNext Login

&#x20; |

&#x20; v

Keycloak

&#x20; |

&#x20; | OpenID Connect (OIDC)

&#x20; |

&#x20; v

ERPNext

&#x20; |

&#x20; v

ERPNext Dashboard



\## Keycloak Configuration



\### Realm



Realm: `erpnext`



\### Client



Client ID: `erpnext`



Protocol: `OpenID Connect`



Client Authentication: Enabled



Standard Authorization Code Flow: Enabled



Direct Access Grants: Disabled



Implicit Flow: Disabled



\## ERPNext Configuration



ERPNext Social Login was configured with Keycloak as the authentication provider.



The following endpoints were configured:



\- Authorization Endpoint

\- Token Endpoint

\- UserInfo Endpoint



User ID mapping:



`preferred\_username`



Authentication parameters:



\- `response\_type=code`

\- `scope=openid`



\## Callback URL



The ERPNext OAuth callback endpoint was configured as:



`http://40.80.93.70/api/method/frappe.integrations.oauth2\_logins/login\_via\_keycloak`



\## Keycloak Base URL



The Keycloak ERPNext realm was accessed through:



`http://40.80.93.70:8080/realms/erpnext`



\## User Mapping



A Keycloak test user was created and mapped to a corresponding ERPNext user.



The ERPNext user was configured as a System User for SSO testing.



\## End-to-End SSO Flow



1\. User opens ERPNext.

2\. User selects Login with Keycloak.

3\. ERPNext redirects the user to Keycloak.

4\. Keycloak authenticates the user.

5\. Keycloak returns the authorization code.

6\. ERPNext exchanges the code for tokens.

7\. ERPNext retrieves the authenticated user information.

8\. User is authenticated into ERPNext.

9\. ERPNext Dashboard opens.



\## Validation



The complete ERPNext + Keycloak SSO flow was successfully tested.



Successful flow:



`ERPNext Login → Keycloak Authentication → ERPNext Dashboard`



\## Security



The following are NOT stored in this repository:



\- Passwords

\- Keycloak client secrets

\- Database credentials

\- Private keys

\- Production credentials



The HTTP configuration was used for the Proof of Concept environment.



For production deployment, HTTPS, secure secret management and restricted infrastructure access should be implemented.


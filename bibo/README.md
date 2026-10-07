\# Bibo Employee Tracking



Bibo is an employee tracking application consisting of a Desktop application and backend services.



This folder contains the Bibo Desktop source code and the related monorepo configuration.



\## Project Structure



```text

bibo/

├── apps/

│   ├── backend/          # Bibo backend service

│   └── desktop/          # Bibo Desktop application

├── docs/                 # Project documentation

├── marketing/            # Marketing-related resources

├── scripts/              # Project scripts

├── package.json          # Root package configuration

├── pnpm-lock.yaml        # Locked dependencies

├── pnpm-workspace.yaml   # PNPM workspace configuration

├── .npmrc                # PNPM/npm configuration

└── README.md

```



\## Technology Stack



\### Desktop



\* Tauri 2

\* Rust

\* TypeScript / JavaScript

\* Node.js

\* PNPM



\### Backend



\* Go

\* HTTP APIs

\* Authentication and session management



\### Authentication



\* Keycloak

\* OpenID Connect (OIDC)

\* Bibo backend authentication APIs



\## Package Manager



This project uses \*\*PNPM\*\*.



Do not use `npm install` to generate a new `package-lock.json`.



Install dependencies with:



```bash

pnpm install

```



The repository contains `pnpm-lock.yaml`, which should be used to keep dependency versions consistent.



\## Requirements



Before running the project, install the required development tools.



Typical requirements include:



\* Node.js 20+

\* PNPM 10+

\* Go

\* Rust

\* Cargo

\* Tauri prerequisites for the target operating system



Check installed versions:



```bash

node --version

pnpm --version

go version

rustc --version

cargo --version

```



\## Install Dependencies



From the Bibo project root:



```bash

pnpm install

```



\## Build / Run



The exact development commands depend on the workspace scripts.



Check available scripts with:



```bash

pnpm run

```



For the Desktop application, use the Tauri development/build commands defined by the project configuration.



\## Keycloak Integration



The Bibo Desktop application has been integrated with the Keycloak authentication flow used by the ERPNext + Keycloak POC.



The objective is to allow an authenticated ERPNext/Keycloak user to activate the Bibo Desktop application without requiring a separate Bibo login.



\### Authentication Flow



```text

User

&#x20;|

&#x20;v

ERPNext

&#x20;|

&#x20;| Login with Keycloak

&#x20;v

Keycloak

&#x20;|

&#x20;| Authentication

&#x20;v

ERPNext Dashboard

&#x20;|

&#x20;| One-time activation

&#x20;v

Bibo Backend

&#x20;|

&#x20;| Session / Tokens

&#x20;v

Bibo Desktop

```



\## Desktop Activation Flow



After the user is authenticated through Keycloak:



1\. The authenticated user is identified by the Bibo backend.

2\. The backend verifies the Keycloak authentication token.

3\. The backend finds the corresponding Bibo employee.

4\. A cryptographically secure one-time activation code is generated.

5\. The activation code is valid for a limited period.

6\. Bibo Desktop sends the activation code to the Bibo backend.

7\. The backend validates the code.

8\. The backend returns a Bibo access token and refresh token.

9\. Bibo Desktop stores the authenticated session.



\## Backend Activation APIs



The Keycloak/Bibo integration adds desktop activation endpoints.



```text

POST /auth/desktop/activation

POST /auth/desktop/activate

```



The Desktop application uses the backend activation API:



```text

POST /v1/auth/desktop/activate

```



The activation request contains a one-time code.



A successful response contains the Bibo session tokens and their expiry information.



\## One-Time Activation Code



The activation code is:



\* Cryptographically generated.

\* Short-lived.

\* Intended for one-time use.

\* Removed after successful activation.

\* Not printed as the actual code in backend logs.



For the current POC, pending activation codes are stored in backend memory.



For a production multi-instance deployment, the activation state should be moved to a shared persistent store so that it works correctly across backend instances and survives backend restarts.



\## Desktop Local Server



The Bibo Desktop application runs a local server for communication between the Desktop application and related components.



The local server includes an activation endpoint:



```text

/activate

```



The Desktop server receives the activation code and communicates with the Bibo backend to exchange it for a valid Bibo session.



\## Authentication State



The Desktop application loads its persisted authentication state during startup.



The local server receives the authentication state and Bibo backend URL so that the Desktop application can participate in the Keycloak/Bibo activation flow.



\## Security



Do not commit any of the following to GitHub:



\* Passwords

\* Keycloak admin credentials

\* Keycloak client secrets

\* Access tokens

\* Refresh tokens

\* API keys

\* Azure credentials

\* SSH private keys

\* Production `.env` files

\* Employee personal data

\* Other production credentials



Use environment variables or a secure secret-management system for environment-specific credentials.



\## Development Notes



The Bibo project uses PNPM rather than npm.



The repository should retain:



```text

pnpm-lock.yaml

pnpm-workspace.yaml

package.json

```



Do not add an unrelated `package-lock.json`.



\## Testing



Before pushing changes, perform appropriate checks for the modified components.



Recommended checks include:



```bash

pnpm install

```



Then run the project-specific tests, linting and formatting commands defined in the repository.



For Rust/Tauri changes, also run the relevant Cargo/Tauri checks.



For Go backend changes, run:



```bash

go build ./...

```



and the applicable Go tests.



\## Git Workflow



Recommended workflow:



```text

Feature / Fix

&#x20;    |

&#x20;    v

Local Development

&#x20;    |

&#x20;    v

Tests

&#x20;    |

&#x20;    v

Git Commit

&#x20;    |

&#x20;    v

GitHub Branch

&#x20;    |

&#x20;    v

Pull Request

&#x20;    |

&#x20;    v

Code Review

&#x20;    |

&#x20;    v

CI/CD

```



The Keycloak integration was developed on the branch:



```text

keycloak-bibo-sso

```



\## Current Integration Status



\### Completed



\* Bibo Desktop source prepared for GitHub.

\* Keycloak authentication integration added.

\* Backend desktop activation functionality added.

\* One-time activation code flow added.

\* Desktop activation endpoint added.

\* Desktop authentication state connected to the local server.

\* Bibo session token handling added.

\* Changes committed to Git.



\### Future Improvements



\* Production HTTPS configuration.

\* Centralized secret management.

\* Persistent/shared activation-code storage.

\* Automated CI/CD.

\* Automated tests for the complete Keycloak → Bibo activation flow.

\* Production monitoring and logging.

\* Secure production deployment configuration.



\## Relation to ERPNext + Keycloak



This Bibo project is part of the broader ERPNext + Keycloak SSO POC.



The overall target architecture is:



```text

&#x20;                   Keycloak

&#x20;                      |

&#x20;            +---------+---------+

&#x20;            |                   |

&#x20;            v                   v

&#x20;         ERPNext             Bibo

&#x20;            |                   |

&#x20;            +---------+---------+

&#x20;                      |

&#x20;                Authenticated

&#x20;                   Employee

```



Keycloak acts as the centralized identity provider, while ERPNext and Bibo consume the authenticated identity through their respective integration flows.

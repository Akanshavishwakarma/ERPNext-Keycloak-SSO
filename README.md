# ERPNext + Keycloak SSO

## Overview

This project contains the deployment and integration documentation for an ERPNext and Keycloak Single Sign-On (SSO) Proof of Concept.

## Architecture

- Azure Ubuntu VM
- ERPNext v15
- MariaDB
- Redis
- Nginx
- Keycloak 26.7.4
- MySQL 8.0
- OpenID Connect (OIDC)

## SSO Flow

ERPNext -> Keycloak -> OIDC Authentication -> ERPNext Dashboard

## Project Status

The ERPNext + Keycloak end-to-end SSO proof of concept was successfully validated.

## Security

Secrets, passwords, private keys and environment-specific credentials must not be committed to this repository.

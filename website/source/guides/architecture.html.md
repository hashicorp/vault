---
layout: "guides"
page_title: "Vault Architecture - Guides"
sidebar_current: "guides-vault-architecture"
description: |-
  In production environments, the availability of the Vault server is critical
  since any downtime may affect the system that is leveraging Vault. Vault is
  designed to support a highly available deploy to ensure a machine or process
  failure is minimally disruptive. The understanding of the Vault architecture
  can help deciding on the Vault backends for your organization's requirements.
---

# Vault Architecture

Vault documentation explained the components.

This guide demonstrates regenerating a root token using a one-time-password (OTP).

## Steps to Regenerate Root Tokens

1. Make sure that the Vault server is unsealed
2. Generate a one-time-password (OTP) to share
3. Each unseal key holder runs `generate-root` with the OTP
4. Decode the generated root token

### Step 1: Make sure that the Vault server is unsealed

First, verify the status:

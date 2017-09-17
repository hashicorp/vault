---
layout: "intro"
page_title: "Vault vs. HSMs"
sidebar_current: "vs-other-hsm"
description: |-
  Comparison between Vault and HSM systems.
---

# Vault vs. HSMs

A [hardware security module
(HSM)](https://en.wikipedia.org/wiki/Hardware_security_module) is a hardware
device that is meant to secure various secrets using protections against access
and tampering at both the software and hardware layers.

The primary issue with HSMs is that they are expensive and not very cloud
friendly. An exception to the latter is Amazon's CloudHSM service, which is
friendly for AWS users but still costs more than $14k per year per instance,
and not as useful for heterogenous cloud architectures.

Once an HSM is up and running, configuring it is generally very tedious, and
the API to request secrets is also difficult to use. Example: CloudHSM requires
SSH and setting up various keypairs manually. It is difficult to automate. APIs
tend to require the use of specific C libraries (e.g. PKCS#11) or
vendor-specific libraries.

However, although configuring and running an HSM can be a challenge, they come
with a significant advantage in that they conform to government-mandated
compliance requirements (e.g. FIPS 140), which often require specific hardware
protections and security models in addition to software.

Vault doesn't replace an HSM. Instead, they can be complementary; a compliant
HSM can protect Vault's master key to help Vault comply with regulatory
requirements, and Vault can provide easy client APIs for tasks such as
encryption and decryption.

Vault can also do many things that HSMs cannot currently do, such as generating
_dynamic secrets_. Instead of storing AWS access keys directly within Vault,
Vault can generate access keys according to a specific policy on the fly. Vault
has the potential of doing this for any system through its mountable secret
backend system.

For many companies' security requirements, Vault alone is enough. For companies
that can afford an HSM or with specific regulatory requirements, it can be used
with Vault to get the best of both worlds.

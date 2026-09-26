/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Static descriptive content for the "About this engine" section on mount-form pages
 * for the five common secrets engines: KV, AWS, Azure, GCP, and Database.
 *
 * Each entry has three bullet arrays (preferredFor, keyFeatures, notSuitedFor). Each bullet
 * is a `BulletPart[]` of inline segments: plain strings, `LinkItem` (`{ text, href }`), or
 * `CodeItem` (`{ code }`) for bold emphasis.
 */

export interface LinkItem {
  text: string;
  href: string;
}

/** Inline bold span — rendered as `<strong>` — for emphasis within a bullet, e.g. `persist_app=true`. */
export interface CodeItem {
  code: string;
}

export type BulletPart = string | LinkItem | CodeItem;

export interface AboutEngineEntry {
  /** Engine type matching `EngineDisplayData.type` in all-engines-metadata */
  type: string;
  /** Bullet points for the "Preferred for" row */
  preferredFor: BulletPart[][];
  /** Bullet points for the "Key features" row */
  keyFeatures: BulletPart[][];
  /** Bullet points for the "Not suited for" row */
  notSuitedFor: BulletPart[][];
}

export const ABOUT_ENGINE_INFO: AboutEngineEntry[] = [
  {
    type: 'kv',
    preferredFor: [
      ['Bearer tokens to proprietary systems'],
      [
        'Public/private key pairs (see also ',
        { text: 'SSH secrets engine', href: 'https://developer.hashicorp.com/vault/docs/secrets/ssh' },
        ')',
      ],
      ["Systems that don't integrate with other Vault secrets engines"],
    ],
    keyFeatures: [
      ['Storing arbitrary data in key-value format (string → string)'],
      [
        {
          text: 'Custom metadata',
          href: 'https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2#metadata',
        },
        ' for context visible to auditors without exposing secret contents',
      ],
      [
        {
          text: 'Version management',
          href: 'https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2#versions',
        },
        ' to track history and roll back to any prior version',
      ],
    ],
    notSuitedFor: [
      ['Frequently updated or short-lived secrets'],
      ["Mandatory rotation cadence: Vault won't enforce or trigger rotation"],
      ['Validating whether credentials are still active or have gone stale'],
    ],
  },
  {
    type: 'aws',
    preferredFor: [
      ['Apps and pipelines that need short-lived, scoped AWS credentials on demand'],
      ['Federating across a wide AWS estate'],
      ['Legacy systems requiring IAM user keys (use static roles so Vault owns rotation)'],
    ],
    keyFeatures: [
      [
        'Four credential types: ',
        { text: 'iam_user', href: 'https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users.html' },
        ', ',
        {
          text: 'assumed_role',
          href: 'https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html',
        },
        ', ',
        {
          text: 'federation_token',
          href: 'https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html',
        },
        ', ',
        {
          text: 'session_token',
          href: 'https://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html',
        },
      ],
      [
        'Per-role ',
        {
          text: 'inline or managed IAM policies',
          href: 'https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_inline-vs-managed.html',
        },
        ', with optional ',
        {
          text: 'permissions boundaries',
          href: 'https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html',
        },
      ],
      ['Static roles: Vault owns a real IAM user and rotates its access key on schedule'],
    ],
    notSuitedFor: [
      [
        'iam_user keys are eventually consistent in IAM. Newly issued keys may fail briefly; prefer assumed_role for modern workloads',
      ],
      ["Root credentials aren't validated at config time, bad IAM perms fail silently"],
    ],
  },
  {
    type: 'azure',
    preferredFor: [
      ['Apps that need scoped, time-bound access to Azure resources via RBAC'],
      ['Teams managing a large Entra estate. Vault centralizes credential issuance across tenants'],
    ],
    keyFeatures: [
      ['By default, creates a new SP per lease with role bindings, deleted on revocation'],
      ['One Vault role → one or more Azure roles, optionally with AD group memberships'],
      [
        'WIF: when Vault is hosted on Azure, its managed identity can serve as root credential instead of static client secret',
      ],
    ],
    notSuitedFor: [
      [
        'Default (fully dynamic) mode hammers Entra on every request, use ',
        { code: 'persist_app=true' },
        ' to avoid throttling or DOS',
      ],
      ['Latency-sensitive flows — issuance can take tens of seconds with many role assignments'],
    ],
  },
  {
    type: 'gcp',
    preferredFor: [
      ['Short-term or one-off GCP access (batch jobs, introspection)'],
      ['Multi-cloud setups where users auth to Vault centrally and get GCP creds on demand'],
    ],
    keyFeatures: [
      ['Two types: access_token (primary) and service_account_key (for legacy, low-activity workloads)'],
      ['Keys are automatically deleted in GCP on lease expiry — no orphaned credentials'],
      ['Vault can directly impersonate reusable roles (access_token only).'],
    ],
    notSuitedFor: [
      ['High-volume key issuance since GCP caps keys per service account; use access_token instead'],
      [
        "Blank config falls back to Application Default Credentials, and if they're unintended, unexpected access can be granted",
      ],
    ],
  },
  {
    type: 'database',
    preferredFor: [
      ['Generating unique, short-lived credentials per request'],
      ['Existing DB users whose passwords Vault should own and rotate on a schedule'],
    ],
    keyFeatures: [
      ['Supports Postgres, MySQL, MongoDB, MSSQL, Redis, Snowflake, Oracle etc.'],
      ['Dynamic roles: unique credentials per request, auto-revoked on lease expiry'],
      ['Static roles: 1-to-1 mapping to a DB user, Vault rotates the password on schedule'],
    ],
    notSuitedFor: [
      ['Root DB credentials assigned to a static role: rotation breaks all users on connection'],
      ['Password rotates instantly on static role onboarding, apps must re-fetch before use'],
      ['Dynamic creds expire with the lease, long-running processes must renew or reconnect'],
    ],
  },
];

/**
 * Returns the `AboutEngineEntry` for the given engine type, or `null` if none exists.
 */
export function getAboutEngineInfo(type: string): AboutEngineEntry | null {
  return ABOUT_ENGINE_INFO.find((entry) => entry.type === type) ?? null;
}

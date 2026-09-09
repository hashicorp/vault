/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/** localStorage key for the set of checklist IDs the user has hidden. */
export const HIDDEN_CHECKLISTS_KEY = 'vault:checklist-hidden';

/** Internal key for known checklist detector implementations. */
export type DetectorKey = 'namespaces' | 'policy' | 'auth' | 'kv';

/** Navigation permission requirement for a checklist step. */
export type NavPermission = [string, string] | null;

/** A call-to-action link rendered inside an expanded checklist step. */
export interface StepCta {
  text: string;
  /** Helios icon name for the CTA button. */
  icon?: string;
  /** Icon placement within the CTA button. */
  iconPosition?: 'leading' | 'trailing';
  /** Internal Ember route to navigate to. */
  route?: string;
  /** Dynamic segments for the route, if required. */
  models?: string[];
  /** External URL (opens in new tab). */
  href?: string;
}

/** Configuration for a single checklist step. */
export interface ChecklistStepConfig {
  label: string;
  description: string;
  ctas: StepCta[];
  /** null means always visible; otherwise requires nav permission. */
  navPermission: NavPermission;
  /** Completion mode controls whether the UI renders manual completion. */
  completion: 'explicit' | 'inferred';
  /** Detector implementation to use for inferred completion. */
  detector?: DetectorKey;
}

/** Full display and behavior configuration for a checklist. */
export interface ChecklistConfig {
  id: string;
  title: string;
  hideButtonText: string;
  steps: Record<string, ChecklistStepConfig>;
  order: string[];
}

/** Step IDs for the cluster-startup checklist. */
export const CLUSTER_STARTUP_CHECKLIST_ID = 'cluster-startup';
export const CLUSTER_STARTUP_STEPS = ['tvp-cli', 'namespaces', 'policy', 'auth', 'kv'] as const;
export type ClusterStartupStepId = (typeof CLUSTER_STARTUP_STEPS)[number];

/** Reusable config for the onboarding checklist shown on the dashboard. */
export const CLUSTER_STARTUP_CHECKLIST: ChecklistConfig = {
  id: CLUSTER_STARTUP_CHECKLIST_ID,
  title: 'Set up your cluster',
  hideButtonText: 'Hide setup guide',
  order: [...CLUSTER_STARTUP_STEPS],
  steps: {
    'tvp-cli': {
      label: 'Install Terraform Vault Provider and CLI',
      description: 'Install Terraform Vault Provider, CLI or UI to get your Vault started.',
      navPermission: null,
      completion: 'explicit',
      ctas: [
        {
          text: 'Installation guide',
          icon: 'docs-link',
          iconPosition: 'leading',
          href: 'https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli',
        },
      ],
    },
    namespaces: {
      label: 'Set up namespaces',
      description:
        'Create secure, isolated environments where independent teams can manage their own resources.',
      navPermission: ['access', 'namespaces'],
      completion: 'inferred',
      detector: 'namespaces',
      ctas: [{ text: 'Namespaces', route: 'vault.cluster.access.namespaces' }],
    },
    policy: {
      label: 'Set up policies',
      description:
        'Define rules to explicitly grant or forbid access to specific paths and operations within your cluster.',
      navPermission: ['policies', 'acl'],
      completion: 'inferred',
      detector: 'policy',
      ctas: [{ text: 'ACL policies', route: 'vault.cluster.policies.index', models: ['acl'] }],
    },
    auth: {
      label: 'Set up authentication',
      description:
        'Set up authentication methods so users and apps can log in and verify their identity with Vault (e.g., via LDAP or Kubernetes) to get a token and access secrets.',
      navPermission: ['access', 'methods'],
      completion: 'inferred',
      detector: 'auth',
      ctas: [{ text: 'Authentication methods', route: 'vault.cluster.access.methods' }],
    },
    kv: {
      label: 'Create your first secrets engine',
      description:
        'Create your first configured secrets engine in the current cluster, from dynamic database credentials and more.',
      navPermission: null,
      completion: 'inferred',
      detector: 'kv',
      ctas: [
        { text: 'Secrets engines', route: 'vault.cluster.secrets' },
        { text: 'Secrets sync', route: 'vault.cluster.sync.secrets.overview' },
      ],
    },
  },
};

/** Registry for checklist configs keyed by checklist ID. */
export const CHECKLIST_CONFIGS: Record<string, ChecklistConfig> = {
  [CLUSTER_STARTUP_CHECKLIST_ID]: CLUSTER_STARTUP_CHECKLIST,
};

/** Looks up a checklist configuration by ID. */
export function getChecklistConfig(checklistId: string): ChecklistConfig | null {
  return CHECKLIST_CONFIGS[checklistId] ?? null;
}

/** Looks up a checklist step configuration by checklist ID and step ID. */
export function getChecklistStepConfig(checklistId: string, stepId: string): ChecklistStepConfig | null {
  return getChecklistConfig(checklistId)?.steps[stepId] ?? null;
}

/** Backward-compatible export used by the current checklist widget. */
export const CLUSTER_STARTUP_STEP_CONFIG = CLUSTER_STARTUP_CHECKLIST.steps;

/*
 * ACL policy names that ship with every Vault cluster and must be excluded when
 * inferring whether a user-created policy exists (policy step detection).
 */
export const BUILTIN_POLICIES = new Set(['default', 'default-ceiling', 'root', 'hcp-root']);

/*
 * The token auth backend path is present on every Vault cluster and must be
 * excluded when inferring whether a user-enabled auth method exists.
 * Some APIs may include leading 'auth/' or omit the trailing slash.
 */
export const BUILTIN_AUTH_PATHS = new Set(['token/', 'token', 'auth/token/', 'auth/token']);

/*
 * KV secret engine paths that are auto-created by the Vault dev server and
 * should not count as user-created mounts. In standard (non-dev) production
 * deployments, no KV mount exists by default.
 */
export const BUILTIN_KV_PATHS = new Set(['secret/']);

/*
 * Steps that require explicit user action to complete — they cannot be inferred
 * from the Vault API. These get a manual completion control in the checklist UI.
 */
export const USER_EXPLICIT_STEPS = new Set(
  CLUSTER_STARTUP_STEPS.filter((stepId) => {
    const step = CLUSTER_STARTUP_CHECKLIST.steps[stepId];
    return step ? step.completion === 'explicit' : false;
  })
);

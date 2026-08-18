/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PREFIX = 'vault_ui';

/*
  buildEventName is a helper to build conformant analytics event names.

  While event names are not strictly controlled in the data warehouse, consistent
    naming helps find things predictably.


*/
const buildEventName = (category: string, resource: string, action: string) =>
  `${PREFIX}_${category}_${resource}_${action}`;

// IBM Tracking Plan generic event names; the specifics (namespace, elementId,
// action, CTA, object, etc.) are carried as event properties.
const UI_INTERACTION = 'UI Interaction';
const CTA_CLICKED = 'CTA Clicked';
const CREATED_OBJECT = 'Created Object';

export const TOGGLE_WEB_REPL = buildEventName('core', 'web-repl', 'toggle');

// Side navigation — cluster level
export const NAV_DASHBOARD = UI_INTERACTION;
export const NAV_SECRETS = UI_INTERACTION;
export const NAV_RESILIENCE_RECOVERY = UI_INTERACTION;
export const NAV_ACCESS_CONTROL = UI_INTERACTION;
export const NAV_OPERATIONAL_TOOLS = UI_INTERACTION;
export const NAV_REPORTING = UI_INTERACTION;
export const NAV_CLIENT_COUNT = UI_INTERACTION;
export const NAV_BILLING_METRICS = UI_INTERACTION;
export const NAV_RAFT_STORAGE = UI_INTERACTION;

// Side navigation — access control sub-items
export const NAV_AUTH_METHODS = UI_INTERACTION;
export const NAV_ACL_POLICIES = UI_INTERACTION;
export const NAV_NAMESPACES = UI_INTERACTION;

// Namespace wizard — step 1 (security policy selection)
export const WIZARD_NAMESPACE_STEP1_POLICY_SELECT = UI_INTERACTION;

// Namespace wizard — step 2 (tips reveal + namespace field input)
export const WIZARD_NAMESPACE_STEP2_TIPS_COLLAPSE = UI_INTERACTION;
export const WIZARD_NAMESPACE_STEP2_FIELD_INPUT = UI_INTERACTION;

// Namespace wizard — step 3 (creation-method selection)
export const WIZARD_NAMESPACE_STEP3_SELECT_TERRAFORM = UI_INTERACTION;
export const WIZARD_NAMESPACE_STEP3_SELECT_CLI_API = UI_INTERACTION;
export const WIZARD_NAMESPACE_STEP3_SELECT_UI = UI_INTERACTION;

// Resource creation — Created Object events fire on both success and failure
// (distinguished by successFlag), with process: 'UI'.
export const SECRET_ENGINE_CREATED = CREATED_OBJECT;
export const AUTH_METHOD_CREATED = CREATED_OBJECT;
export const POLICY_CREATED = CREATED_OBJECT;
// Namespace wizard — namespace creation (fires on success and failure)
export const NAMESPACE_CREATED = CREATED_OBJECT;

// Policy creation cancelled
export const POLICY_CREATION_CANCELLED = UI_INTERACTION;

// Dashboard quick actions and secrets engines widget
export const DASHBOARD_QUICK_ACTION_CTA_CLICKED = CTA_CLICKED;
export const DASHBOARD_SECRETS_ENGINE_VIEW_ALL = UI_INTERACTION;
export const DASHBOARD_SECRETS_ENGINE_VIEW = UI_INTERACTION;

// Intro pages
export const INTRO_SECRETS_ENGINE_CTA_CLICKED = CTA_CLICKED;
export const INTRO_AUTH_METHODS_CTA_CLICKED = CTA_CLICKED;
export const INTRO_ACL_POLICIES_CTA_CLICKED = CTA_CLICKED;
// Namespace wizard — intro CTA (guided start / docs / dismiss)
export const INTRO_NAMESPACES_CTA_CLICKED = CTA_CLICKED;

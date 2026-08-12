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

// IBM Tracking Plan generic event name; the specifics (namespace, elementId,
// action, CTA) are carried as event properties.
const UI_INTERACTION = 'UI Interaction';

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

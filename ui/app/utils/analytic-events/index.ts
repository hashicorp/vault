/**
 * Copyright (c) HashiCorp, Inc.
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

export const TOGGLE_WEB_REPL = buildEventName('core', 'web-repl', 'toggle');

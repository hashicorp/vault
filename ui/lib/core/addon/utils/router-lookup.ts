/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { getOwner } from '@ember/owner';

import type RouterService from '@ember/routing/router-service';

// In components shared across engines, we have to look up the router dynamically
// and use getOwner because Ember Engine's use the alias 'app-router'.
// This way we avoid throwing an error by looking up a service that doesn't exist.
// https://guides.emberjs.com/release/services/#toc_accessing-services
export default function routerLookup(context: object) {
  const owner = getOwner(context);
  return (
    (owner?.lookup('service:router') as RouterService) ||
    (owner?.lookup('service:app-router') as RouterService)
  );
}

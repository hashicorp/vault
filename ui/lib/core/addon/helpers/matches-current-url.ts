/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Helper from '@ember/component/helper';
import routerLookup from 'core/utils/router-lookup';
import { sanitizeStart } from 'core/utils/sanitize-path';
import { assert } from '@ember/debug';

import type RouterService from '@ember/routing/router-service';

type RouteComputeOptions = {
  isExactMatch?: boolean;
};

/**
 * Uses recognize() to determine the route from the current URL and whether
 * it is an ancestor of the passed route (a substring) or an exact match.
 *
 * @see https://api.emberjs.com/ember/release/classes/routerservice/#recognize
 * @example
 * <LinkTo @current-when={{matches-current-url "vault.secrets.backend" isExactMatch=true}} />
 */

export function matchesCurrentUrl(
  router: RouterService,
  routeName: string,
  options: RouteComputeOptions = {}
): boolean {
  assert(
    `routeName is required, you passed ${routeName === '' ? 'an empty string' : routeName}`,
    !!routeName
  );

  if (!router) {
    return false;
  }

  // Remove leading slash because rootUrl has trailing and leading slash
  const sanitizedUrl = sanitizeStart(router.currentURL);
  // Must pass a url that begins with the application's rootURL.
  // Set dynamically but it is `/ui/`
  const url = `${router.rootURL}${sanitizedUrl}`;
  const currentRouteInfo = router.recognize(url);

  if (!currentRouteInfo) {
    return false;
  }

  if (options.isExactMatch) {
    return currentRouteInfo.name === routeName;
  }

  // Check if current route includes the passed route string (returns false for empty strings)
  return currentRouteInfo.name.includes(routeName);
}

export default class MatchesCurrentUrl extends Helper {
  get router(): RouterService {
    return routerLookup(this);
  }

  compute([routeName]: [string], options: RouteComputeOptions): boolean {
    return matchesCurrentUrl(this.router, routeName, options);
  }
}

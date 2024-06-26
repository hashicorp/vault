/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class OidcConfigureController extends Controller {
  @service router;

  @tracked header = null;

  constructor() {
    super(...arguments);
    this.router.on('routeDidChange', (transition) => this.setHeader(transition));
  }

  setHeader(transition) {
    // set correct header state based on child route
    // when no clients have been created, display create button as call to action
    // list views share the same header with tabs as resource links
    // the remaining routes are responsible for their own header
    const routeName = transition.to.name;
    if (routeName.includes('oidc.index')) {
      this.header = 'cta';
    } else {
      const isList = ['clients', 'assignments', 'keys', 'scopes', 'providers'].find((resource) => {
        return routeName.includes(`${resource}.index`);
      });
      this.header = isList ? 'list' : null;
    }
  }

  get isCta() {
    return this.header === 'cta';
  }
}

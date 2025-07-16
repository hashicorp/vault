/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';

export default Component.extend({
  router: service(),
  controlGroup: service(),
  api: service(),

  // public attrs
  model: null,
  controlGroupResponse: null,

  //internal state
  error: null,
  unwrapData: null,

  unwrap: task(function* (token) {
    this.set('error', null);
    try {
      const response = yield this.api.sys.unwrap({}, this.api.buildHeaders({ token }));
      this.set('unwrapData', response.auth || response.data);
      this.controlGroup.deleteControlGroupToken(this.model.id);
    } catch (e) {
      const { message } = yield this.api.parseError(e);
      this.error = `Token unwrap failed: ${message}`;
    }
  }).drop(),

  markAndNavigate: task(function* () {
    this.controlGroup.markTokenForUnwrap(this.model.id);
    const { url } = this.controlGroupResponse.uiParams;
    yield this.router.transitionTo(url);
  }).drop(),
});

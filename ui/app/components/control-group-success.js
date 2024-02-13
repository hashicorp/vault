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
  store: service(),

  // public attrs
  model: null,
  controlGroupResponse: null,

  //internal state
  error: null,
  unwrapData: null,

  unwrap: task(function* (token) {
    const adapter = this.store.adapterFor('tools');
    this.set('error', null);
    try {
      const response = yield adapter.toolAction('unwrap', null, { clientToken: token });
      this.set('unwrapData', response.auth || response.data);
      this.controlGroup.deleteControlGroupToken(this.model.id);
    } catch (e) {
      this.set('error', `Token unwrap failed: ${e.errors[0]}`);
    }
  }).drop(),

  markAndNavigate: task(function* () {
    this.controlGroup.markTokenForUnwrap(this.model.id);
    const { url } = this.controlGroupResponse.uiParams;
    yield this.router.transitionTo(url);
  }).drop(),
});

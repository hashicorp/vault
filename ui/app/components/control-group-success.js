/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

export default Component.extend({
  router: service(),
  controlGroup: service(),
  api: service(),
  store: service(),

  // public attrs
  model: null,
  controlGroupResponse: null,

  //internal state
  error: null,
  unwrapData: null,

  unwrap: task(function* (token) {
    this.set('error', null);
    try {
      const adapter = this.store.adapterFor('application');
      const response = yield adapter.ajax('/v1/sys/wrapping/unwrap', 'POST', { clientToken: token });
      this.set('unwrapData', response.auth || response.data);
      this.controlGroup.deleteControlGroupToken(this.model.id);
    } catch (e) {
      this.error = `Token unwrap failed: ${errorMessage(e)}`;
    }
  }).drop(),

  markAndNavigate: task(function* () {
    this.controlGroup.markTokenForUnwrap(this.model.id);
    const { url } = this.controlGroupResponse.uiParams;
    yield this.router.transitionTo(url);
  }).drop(),
});

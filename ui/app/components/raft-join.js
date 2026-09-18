/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import RaftJoinForm from 'vault/forms/storage/raft-join';

/**
 * @module RaftJoin
 * RaftJoin component presents the user with a choice to join an existing raft cluster when a new Vault
 * server is brought up
 *
 *
 * @example
 * ```js
 * <RaftJoin @onDismiss={{action (mut attr)}} />
 * ```
 * @param {function} onDismiss - This function will be called if the user decides not to join an existing
 * raft cluster
 *
 */

import Component from '@ember/component';

export default Component.extend({
  classNames: 'raft-join',
  api: service(),
  router: service(),
  onDismiss() {},
  preference: 'join',
  showJoinForm: false,
  form: null,
  modelValidations: null,
  errorMessage: null,

  save: task(
    waitFor(function* (event) {
      event.preventDefault();
      this.set('errorMessage', null);
      const { isValid, state, data } = this.form.toJSON();
      // always refresh modelValidations (not just when invalid) so a previously-shown error
      // clears once the field is corrected, instead of lingering stale on a valid resubmit
      this.set('modelValidations', state);
      if (!isValid) {
        return;
      }
      try {
        yield this.api.request.post('/sys/storage/raft/join', data);
      } catch (e) {
        const { message } = yield this.api.parseError(e);
        this.set('errorMessage', message);
        return;
      }
      this.router.transitionTo('vault.cluster.unseal');
    })
  ).drop(),

  actions: {
    advanceFirstScreen(event) {
      event.preventDefault();
      if (this.preference !== 'join') {
        this.onDismiss();
        return;
      }
      this.set('form', new RaftJoinForm());
      this.set('showJoinForm', true);
    },
    cancel() {
      this.set('showJoinForm', false);
      this.set('form', null);
      this.set('modelValidations', null);
      this.set('errorMessage', null);
    },
  },
});

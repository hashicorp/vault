/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';

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
  store: service(),
  onDismiss() {},
  preference: 'join',
  showJoinForm: false,
  actions: {
    advanceFirstScreen(event) {
      event.preventDefault();
      if (this.preference !== 'join') {
        this.onDismiss();
        return;
      }
      this.set('showJoinForm', true);
    },
    newModel() {
      return this.store.createRecord('raft-join');
    },
  },
});

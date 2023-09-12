/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

export default class SecretDeleteMenu extends Component {
  @service router;
  @service store;

  @action
  handleDelete() {
    this.args.model.destroyRecord().then(() => {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    });
  }
}

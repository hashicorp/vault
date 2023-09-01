/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { alias } from '@ember/object/computed';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';

export default class SecretDeleteMenu extends Component {
  @service router;

  @maybeQueryRecord(
    'capabilities',
    (context) => {
      if (!context.args.model || context.args.mode === 'create') {
        return;
      }
      const { backend, id } = context.args.model;
      const path = `${backend}/${id}`;
      return {
        id: path,
      };
    },
    'model',
    'model.id',
    'mode'
  )
  secretDataPath;
  @alias('secretDataPath.canDelete') canDeleteSecretData;

  @action
  handleDelete() {
    this.args.model.destroyRecord().then(() => {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    });
  }
}

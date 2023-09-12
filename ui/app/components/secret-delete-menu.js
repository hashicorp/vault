/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class SecretDeleteMenu extends Component {
  @service router;
  @service flashMessages;

  get canUndeleteVersion() {
    return this.args.modelForData.canUndeleteVersion;
  }
  get canDestroyVersion() {
    return this.args.modelForData.canDestroyVersion;
  }
  get canDestroyAllVersions() {
    return this.args.modelForData.canDestroyAllVersions;
  }
  get canDeleteSecretData() {
    return this.args.modelForData.canDeleteSecretData;
  }
  get canSoftDeleteSecretData() {
    return this.args.modelForData.canSoftDeleteSecretData;
  }

  get isLatestVersion() {
    // must have metadata access.
    const { model } = this.args;
    if (!model) return false;
    const latestVersion = model.currentVersion;
    const selectedVersion = model.selectedVersion.version;
    if (latestVersion !== selectedVersion) {
      return false;
    }
    return true;
  }

  @action
  handleDelete() {
    this.args.model.destroyRecord().then(() => {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    });
  }
}

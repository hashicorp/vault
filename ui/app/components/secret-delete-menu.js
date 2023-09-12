/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import Ember from 'ember';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

const getErrorMessage = (errors) => {
  const errorMessage =
    errors?.join('. ') || 'Something went wrong. Check the Vault logs for more information.';
  return errorMessage;
};
export default class SecretDeleteMenu extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @tracked showDeleteModal = false;

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
  handleDelete(deleteType) {
    // deleteType should be 'delete', 'destroy', 'undelete', 'delete-latest-version', 'destroy-all-versions', 'v1'
    if (!deleteType) {
      return;
    }
    if (deleteType === 'destroy-all-versions' || deleteType === 'v1') {
      this.args.model.destroyRecord().then(() => {
        return this.router.transitionTo('vault.cluster.secrets.backend.list-root');
      });
    } else {
      // if they do not have read access on the metadata endpoint we need to pull the version from modelForData so they can perform delete and undelete operations
      // only perform if no access to metadata otherwise it will only delete latest version for any deleteType === delete
      let currentVersionForNoReadMetadata;
      if (!this.args.canReadSecretMetadata) {
        currentVersionForNoReadMetadata = this.args.modelForData?.version;
      }
      return this.store
        .adapterFor('secret-v2-version')
        .v2DeleteOperation(this.store, this.args.modelForData.id, deleteType, currentVersionForNoReadMetadata)
        .then((resp) => {
          if (Ember.testing) {
            this.showDeleteModal = false;
            // we don't want a refresh otherwise test loop will rerun in a loop
            return;
          }
          if (!resp) {
            this.showDeleteModal = false;
            this.args.refresh();
            return;
          }
          if (resp.isAdapterError) {
            const errorMessage = getErrorMessage(resp.errors);
            this.flashMessages.danger(errorMessage);
          } else {
            // not likely to ever get to this situation, but adding just in case.
            location.reload();
          }
        });
    }
  }
}

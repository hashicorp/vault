/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

/**
 * @module KvSecretMetadataDetails renders the details view for kv/metadata and button to delete (which deletes the whole secret) or edit metadata.
 * <Page::Secret::Metadata::Details
 * @backend={{this.model.backend}}
 * @breadcrumbs={{this.breadcrumbs}}
 * @canDeleteMetadata={{this.model.canDeleteMetadata}}
 * @canReadData={{this.model.canReadData}}
 * @canReadMetadata={{this.model.canReadMetadata}}
 * @canUpdateMetadata={{this.model.canUpdateMetadata}}
 * @metadata={{this.model.metadata}}
 * @path={{this.model.path}}
 * />
 *
 * @param {string} backend - The name of the kv secret engine.
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {boolean} canDeleteMetadata - if true, "Permanently delete" action renders in the toolbar
 * @param {boolean} canReadData - if true, user can make a request for custom_metadata if they don't have "read" permissions for metadata
 * @param {boolean} canReadMetadata - if true, secret metadata renders below custom_metadata
 * @param {boolean} canUpdateMetadata - if true, "Edit" action renders in the toolbar
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 *
 *
 */

export default class KvSecretMetadataDetails extends Component {
  @service controlGroup;
  @service flashMessages;
  @service('app-router') router;
  @service store;

  @tracked error = null;
  @tracked customMetadataFromData = null;

  get customMetadata() {
    return this.args.metadata?.customMetadata || this.customMetadataFromData;
  }

  @action
  async onDelete() {
    // The only delete option from this view is delete metadata and all versions
    const { backend, path } = this.args;
    const adapter = this.store.adapterFor('kv/metadata');
    try {
      await adapter.deleteMetadata(backend, path);
      this.store.clearDataset('kv/metadata'); // Clear out the store cache so that the metadata/list view is updated.
      this.flashMessages.success(
        `Successfully deleted the metadata and all version data for the secret ${path}.`
      );
      this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
    } catch (err) {
      this.flashMessages.danger(`There was an issue deleting ${path} metadata. \n ${errorMessage(err)}`);
    }
  }

  @action
  async requestData() {
    const { backend, path } = this.args;
    try {
      const secretData = await this.store.queryRecord('kv/data', { backend, path });
      this.customMetadataFromData = secretData.customMetadata;
    } catch (error) {
      if (error.message === 'Control Group encountered') {
        this.controlGroup.saveTokenFromError(error);
        this.error = this.controlGroup.logFromError(error);
        this.error.isControlGroup = true;
        return;
      }
      this.error.isControlGroup = false;
      this.error = errorMessage(error);
    }
  }
}

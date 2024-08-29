/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

/**
 * @module KvSecretMetadataDetails renders the details view for kv/metadata and button to delete (which deletes the whole secret) or edit metadata.
 * <Page::Secret::Metadata::Details
 * @backend={{this.model.backend}}
 * @breadcrumbs={{this.breadcrumbs}}
 * @canDeleteMetadata={{this.model.permissions.metadata.canDelete}}
 * @canReadMetadata={{this.model.permissions.metadata.canRead}}
 * @canUpdateMetadata={{this.model.permissions.metadata.canUpdate}}
 * @customMetadata={{or this.model.metadata.customMetadata this.model.secret.customMetadata}}
 * @metadata={{this.model.metadata}}
 * @path={{this.model.path}}
 * />
 *
 * @param {string} backend - The name of the kv secret engine.
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {boolean} canDeleteMetadata - if true, "Permanently delete" action renders in the toolbar
 * @param {boolean} canReadMetadata - if true, secret metadata renders below custom_metadata
 * @param {boolean} canUpdateMetadata - if true, "Edit" action renders in the toolbar
 * @param {object} customMetadata - comes from secret metadata or data endpoint. if undefined, user does not have "read" access, if an empty object then there is none
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 *
 *
 */

export default class KvSecretMetadataDetails extends Component {
  @service flashMessages;
  @service router;
  @service store;

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
}

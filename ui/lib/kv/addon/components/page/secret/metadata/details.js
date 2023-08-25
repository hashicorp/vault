/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

/**
 * @module KvSecretMetadataDetails renders the details view for kv/metadata.
 * It also renders a button to delete metadata.
 * <Page::Secret::Metadata::Details
 *  @path={{this.model.path}}
 *  @secret={{this.model.secret}}
 *  @metadata={{this.model.metadata}}
 *  @breadcrumbs={{this.breadcrumbs}}
  />
 *
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 * @param {model} [secret] - Ember data model: 'kv/data'. Param not required for delete-metadata.
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvSecretMetadataDetails extends Component {
  @service flashMessages;
  @service router;
  @service store;

  get customMetadata() {
    // metadata tab is available even if user only has access to kv/data path
    return this.args.metadata?.customMetadata || this.args.secret?.customMetadata;
  }
  @action
  async onDelete() {
    // The only delete option from this view is delete all versions
    const { secret } = this.args;
    try {
      await secret.destroyRecord({
        adapterOptions: { deleteType: 'destroy-all-versions', deleteVersions: this.version },
      });
      this.store.clearDataset('kv/metadata'); // Clear out the store cache so that the metadata/list view is updated.
      this.flashMessages.success(
        `Successfully deleted the metadata and all version data for the secret ${secret.path}.`
      );
      this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
    } catch (err) {
      this.flashMessages.danger(`There was an issue deleting ${secret.path} metadata.`);
    }
  }
}

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
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
  @tracked deleteModalOpen = false;
  @service flashMessages;
  @service router;

  @action async handleDelete() {
    // the only delete option from this view is delete on metadata.
    try {
      await this.args.metadata.destroyRecord({
        adapterOptions: { deleteType: 'delete-metadata' },
      });
      this.flashMessages.success(
        `Successfully deleted the metadata and all version data for the secret ${this.args.metadata.path}.`
      );
      return this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
    } catch (err) {
      this.flashMessages.danger(`There was an issue deleting ${this.args.metadata.path} metadata.`);
    }
  }
}

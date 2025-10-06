/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { waitFor } from '@ember/test-waiters';

/**
 * @module KvSecretMetadataDetails renders the details view for kv metadata and button to delete (which deletes the whole secret) or edit metadata.
 * <Page::Secret::Metadata::Details
 *   @backend={{this.model.backend}}
 *   @breadcrumbs={{this.breadcrumbs}}
 *   @capabilities={{this.model.capabilities}}
 *   @metadata={{this.model.metadata}}
 *   @path={{this.model.path}}
 * />
 *
 * @param {string} backend - The name of the kv secret engine.
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {object} capabilities - capabilities for data, metadata, subkeys, delete and undelete paths
 * @param {object} metadata - kv metadata
 * @param {string} path - path of kv secret 'my/secret' used as the title for the KV page header
 *
 *
 */

export default class KvSecretMetadataDetails extends Component {
  @service controlGroup;
  @service flashMessages;
  @service('app-router') router;
  @service api;

  @tracked error = null;
  @tracked customMetadataFromData = null;
  @tracked didRequestData = false;

  get customMetadata() {
    return this.args.metadata?.custom_metadata || this.customMetadataFromData;
  }

  get canRequestData() {
    const { canReadMetadata, canReadData } = this.args.capabilities;
    return !canReadMetadata && canReadData && !this.didRequestData;
  }

  @action
  async onDelete() {
    // The only delete option from this view is delete metadata and all versions
    const { backend, path } = this.args;
    try {
      await this.api.secrets.kvV2DeleteMetadataAndAllVersions(path, backend);
      this.flashMessages.success(
        `Successfully deleted the metadata and all version data for the secret ${path}.`
      );
      this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
    } catch (err) {
      this.flashMessages.danger(`There was an issue deleting ${path} metadata. \n ${errorMessage(err)}`);
    }
  }

  @action
  @waitFor
  async requestData() {
    const { backend, path } = this.args;
    try {
      const { metadata } = await this.api.secrets.kvV2Read(path, backend);
      this.customMetadataFromData = metadata.custom_metadata;
      this.didRequestData = true;
    } catch (err) {
      const { message, response } = await this.api.parseError(err);
      if (response.isControlGroupError) {
        this.controlGroup.saveTokenFromError(response);
        this.error = this.controlGroup.logFromError(response);
        this.error.isControlGroup = true;
      } else {
        this.error.isControlGroup = false;
        this.error = message;
      }
    }
  }
}

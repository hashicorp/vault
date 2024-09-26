/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { dateFormat } from 'core/helpers/date-format';
import { isDeleted } from 'kv/utils/kv-deleted';

/**
 * @module KvSecretOverview
 * <Page::Secret::Overview
 * @backend={{this.model.backend}}
 * @breadcrumbs={{this.breadcrumbs}}
 * @canReadMetadata={{true}}
 * @canUpdateData={{true}}
 * @isPatchAllowed={{true}}
 * @metadata={{this.model.metadata}}
 * @path={{this.model.path}}
 * @subkeys={{this.model.subkeys}}
 * />
 *
 * @param {string} backend - kv secret mount to make network request
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {boolean} canReadMetadata - permissions to read metadata
 * @param {boolean} canUpdateData - permissions to create a new version of a secret
 * @param {boolean} isPatchAllowed - passed to KvSubkeysCard for rendering patch action. True when: (1) the version is enterprise, (2) a user has "patch" secret + "read" subkeys capabilities, (3) latest secret version is not deleted or destroyed
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {string} path - path to request secret data for selected version
 * @param {object} subkeys - API response from subkeys endpoint, object with "subkeys" and "metadata" keys. This arg is null for community edition
 */

export default class KvSecretOverview extends Component {
  get secretState() {
    if (this.args.metadata) {
      return this.args.metadata.currentSecret.state;
    }
    if (this.args.subkeys?.metadata) {
      const { metadata } = this.args.subkeys;
      const state = metadata.destroyed
        ? 'destroyed'
        : isDeleted(metadata.deletion_time)
        ? 'deleted'
        : 'created';
      return state;
    }
    return 'created';
  }

  get versionSubtext() {
    const state = this.secretState;
    if (state === 'destroyed') {
      return 'The current version of this secret has been permanently deleted and cannot be restored.';
    }
    if (state === 'deleted') {
      const time =
        this.args.metadata?.currentSecret.deletionTime || this.args.subkeys?.metadata.deletion_time;
      const date = dateFormat([time], {});
      return `The current version of this secret was deleted ${date}.`;
    }
    return 'The current version of this secret.';
  }
}

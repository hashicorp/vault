/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { dateFormat } from 'core/helpers/date-format';

/**
 * @module KvSecretOverview
 * <Page::Secret::Overview
 * @backend={{this.model.backend}}
 * @breadcrumbs={{this.breadcrumbs}}
 * @canReadMetadata={{true}}
 * @canUpdateSecret={{true}}
 * @metadata={{this.model.metadata}}
 * @path={{this.model.path}}
 * @secretState="created"
 * @subkeys={{this.model.subkeys}}
 * />
 *
 * @param {string} backend - kv secret mount to make network request
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {boolean} canReadMetadata - permissions to read metadata
 * @param {boolean} canUpdateSecret - permissions to create a new version of a secret
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {string} path - path to request secret data for selected version
 * @param {string} secretState - if a secret has been "destroyed", "deleted" or "created" (still active)
 * @param {object} subkeys - API response from subkeys endpoint, object with "subkeys" and "metadata" keys
 */

export default class KvSecretOverview extends Component {
  get isActive() {
    const state = this.args.secretState;
    return state !== 'destroyed' && state !== 'deleted';
  }

  get versionSubtext() {
    const state = this.args.secretState;
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

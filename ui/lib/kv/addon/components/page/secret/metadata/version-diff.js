/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { kvDataPath } from 'vault/utils/kv-path';

/**
 * @module KvSecretMetadataVersionDiff renders the version diff comparison
 * <Page::Secret::Metadata::VersionDiff
 *  @metadata={{this.model.metadata}}
 *  @path={{this.model.path}}
 *  @backend={{this.model.backend}}
 *  @breadcrumbs={{this.breadcrumbs}}
 * />
 *
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {string} path - path to request secret data for selected version
 * @param {string} backend - kv secret mount to make network request
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

/* eslint-disable no-undef */
export default class KvSecretMetadataVersionDiff extends Component {
  @service store;
  @tracked leftVersion;
  @tracked rightVersion;
  @tracked visualDiff;
  @tracked statesMatch = false;

  constructor() {
    super(...arguments);

    // initialize with most recently (before current), active version on left
    const olderVersions = this.args.metadata.sortedVersions.slice(1);
    const recentlyActive = olderVersions.find((v) => !v.destroyed && !v.isSecretDeleted);
    this.leftVersion = Number(recentlyActive?.version);
    this.rightVersion = this.args.metadata.currentVersion;

    // this diff is from older to newer (current) secret data
    this.createVisualDiff();
  }

  // this can only be true on initialization if the current version is inactive
  // selecting a deleted/destroyed version is otherwise disabled
  get deactivatedState() {
    const { currentVersion, currentSecret } = this.args.metadata;
    return this.rightVersion === currentVersion && currentSecret.isDeactivated ? currentSecret.state : '';
  }

  @action
  handleSelect(side, version, actions) {
    this[side] = Number(version);
    actions.close();
    this.createVisualDiff();
  }

  async createVisualDiff() {
    const leftSecretData = await this.fetchSecretData(this.leftVersion);
    const rightSecretData = await this.fetchSecretData(this.rightVersion);
    const diffpatcher = jsondiffpatch.create({});
    const delta = diffpatcher.diff(leftSecretData, rightSecretData);

    this.statesMatch = !delta;
    this.visualDiff = delta
      ? jsondiffpatch.formatters.html.format(delta, leftSecretData)
      : JSON.stringify(rightSecretData, undefined, 2);
  }

  async fetchSecretData(version) {
    const { backend, path } = this.args;
    // check the store first, avoiding an extra network request if possible.
    const storeData = await this.store.peekRecord('kv/data', kvDataPath(backend, path, version));
    const data = storeData ? storeData : await this.store.queryRecord('kv/data', { backend, path, version });

    return data?.secretData;
  }
}

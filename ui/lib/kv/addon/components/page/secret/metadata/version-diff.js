/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { kvDataPath } from 'vault/utils/kv-path';
import { typeOf } from '@ember/utils';

/**
 * @module KvVersionDiff
 * This component produces a JSON diff view between 2 secret versions. It uses the library jsondiffpatch.
 *
 * @param {string} backend - Backend from the kv/data model.
 * @param {string} path - Backend from the kv/data model.
 * @param {array} metadata - The kv/metadata model. It is version agnostic.
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component.
 * @param {object} currentSecretData - The model class for kv/data. Need version and deleted or destroyed information.
 */

export default class KvVersionDiffComponent extends Component {
  @service store;
  @tracked leftSideVersion;
  @tracked rightSideVersion;
  @tracked visualDiff;
  @tracked statesMatch = false;

  constructor() {
    super(...arguments);
    // tracked properties set here because they use args.
    this.leftSideVersion = this.args.metadata.currentVersion;
    this.rightSideVersion = this.defaultRightSideVersion;
    this.createVisualDiff();
  }

  get defaultRightSideVersion() {
    // unless the version is destroyed or deleted we return the version prior to the current version.
    const versionData = this.args.metadata.sortedVersions.find(
      (version) =>
        version.destroyed === false && version.deletion_time === '' && version.version != this.leftSideVersion
    );
    return versionData ? versionData.version : this.leftSideVersion - 1;
  }

  async fetchSecretData(version) {
    version = typeOf(version) === 'string' ? Number(version) : version;
    const { backend, path } = this.args;
    // check the store first, avoiding an extra network request if possible.
    const storeData = await this.store.peekRecord('kv/data', kvDataPath(backend, path, version));
    const data = storeData ? storeData : await this.store.queryRecord('kv/data', { backend, path, version });

    return data?.secretData;
  }

  async createVisualDiff() {
    /* eslint-disable no-undef */
    const leftSideData = await this.fetchSecretData(this.leftSideVersion);
    const rightSideData = await this.fetchSecretData(this.rightSideVersion);
    const diffpatcher = jsondiffpatch.create({});
    const delta = diffpatcher.diff(rightSideData, leftSideData);

    // ARG TODO account for destroyed, deleted current version
    this.statesMatch = !delta;
    this.visualDiff = delta
      ? jsondiffpatch.formatters.html.format(delta, rightSideData)
      : JSON.stringify(leftSideData, undefined, 2);
  }

  @action selectVersion(version, actions, side) {
    if (side === 'right') {
      this.rightSideVersion = version;
    }
    if (side === 'left') {
      this.leftSideVersion = version;
    }
    // close dropdown menu.
    if (actions) actions.close();
    this.createVisualDiff();
  }
}

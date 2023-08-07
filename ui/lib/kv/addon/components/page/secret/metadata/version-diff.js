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
 * @param {array} metadata - The kv/metadata model. It is version agnostic.
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvVersionDiffComponent extends Component {
  @service store;

  @tracked leftSideVersion = this.args.secret.version;
  @tracked rightSideVersion = this.defaultRightSideVersion;
  @tracked statesMatch = false;
  @tracked visualDiff;

  constructor() {
    super(...arguments);
    this.createVisualDiff();
  }

  get defaultRightSideVersion() {
    // unless the version is destroyed or deleted we return the version prior to the current version.
    const versionData = this.args.metadata.sortedVersions.find(
      (version) =>
        version.destroyed === false && version.deletion_time === '' && version.version != this.leftSideVersion
    );
    // if all versions have been deleted or destroyed return one less than the current version;
    return versionData ? versionData.version : this.leftSideVersion - 1;
  }

  async fetchSecretData(version) {
    // ARG TODO ask if maybe we want to override the adapter peekRecord?
    version = typeOf(version) === 'string' ? Number(version) : version;
    const { backend, path } = this.args.secret;
    // check the store first, avoiding an extra network request if possible.
    const storeData = await this.store.peekRecord('kv/data', kvDataPath(backend, path, version));
    const data = storeData ? storeData : await this.store.queryRecord('kv/data', { backend, path, version });

    return data?.secretData;
  }

  async createVisualDiff() {
    /* eslint-disable no-undef */
    const diffpatcher = jsondiffpatch.create({});
    const leftSideVersionData = await this.fetchSecretData(this.leftSideVersion);
    const rightSideVersionData = await this.fetchSecretData(this.rightSideVersion);
    const delta = diffpatcher.diff(rightSideVersionData, leftSideVersionData);

    this.statesMatch = !delta;
    this.visualDiff = delta
      ? jsondiffpatch.formatters.html.format(delta, rightSideVersionData)
      : JSON.stringify(leftSideVersionData, undefined, 2);
  }

  @action selectVersion(version, actions, side) {
    if (side === 'right') {
      this.rightSideVersion = version;
    }
    if (side === 'left') {
      this.leftSideVersion = version;
    }
    this.createVisualDiff();
    // close dropdown menu.
    actions.close();
  }
}

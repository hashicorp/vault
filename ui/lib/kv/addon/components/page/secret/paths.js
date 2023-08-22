/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { kvMetadataPath, kvDataPath } from 'vault/utils/kv-path';

/**
 * @module KvSecretPaths is used to display copyable secret paths to use on the CLI and API.
 *
 * <Page::Secret::Paths
 *  @path={{this.model.path}}
 *  @secret={{this.model.secret}}
 *  @breadcrumbs={{this.breadcrumbs}}
 * />
 *
 * @param {model} secret - Ember data model: 'kv/data', the new record for the new secret version saved by the form
 * @param {string} path - kv secret path for rendering the page title
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvSecretPaths extends Component {
  @service namespace;

  get paths() {
    const { backend, path } = this.args.secret;
    const namespace = this.namespace.path;
    const cli = `-mount=${backend} "${path}"`;
    const data = kvDataPath(backend, path);
    const metadata = kvMetadataPath(backend, path);

    return {
      'API path': {
        path: namespace ? `/v1/${namespace}/${data}` : `/v1/${data}`,
        text: 'Use this path when referring to this secret in the API.',
      },
      'CLI path': {
        path: namespace ? `-namespace=${namespace} ${cli}` : cli,
        text: 'Use this path when referring to this secret in the CLI.',
      },
      'API path for metadata': {
        path: namespace ? `/v1/${namespace}/${metadata}` : `/v1/${metadata}`,
        text: `Use this path when referring to this secret's metadata in the API and permanent secret deletion.`,
      },
    };
  }
}

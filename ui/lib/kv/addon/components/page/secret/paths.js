/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { kvMetadataPath, kvDataPath } from 'vault/utils/kv-path';

/**
 * @module KvSecretPaths is used to display copyable secret paths for KV v2 for CLI and API use.
 * This view is permission agnostic because args come from the views mount path and url params.
 *
 * <Page::Secret::Paths
 *  @path={{this.model.path}}
 *  @backend={{this.model.backend}}
 *  @breadcrumbs={{this.breadcrumbs}}
 *  @canReadMetadata={{this.model.secret.canReadMetadata}}
 * />
 *
 * @param {string} path - kv secret path for building the CLI and API paths
 * @param {string} backend - the secret engine mount path, comes from the secretMountPath service defined in the route
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * @param {boolean} [canReadMetadata=true] - if true, displays tab for Version History
 */

export default class KvSecretPaths extends Component {
  @service namespace;

  get paths() {
    const { backend, path } = this.args;
    const namespace = this.namespace.path;
    const cli = `-mount="${backend}" "${path}"`;
    const data = kvDataPath(backend, path);
    const metadata = kvMetadataPath(backend, path);

    return [
      {
        label: 'API path',
        snippet: namespace ? `/v1/${namespace}/${data}` : `/v1/${data}`,
        text: 'Use this path when referring to this secret in the API.',
      },
      {
        label: 'CLI path',
        snippet: namespace ? `-namespace=${namespace} ${cli}` : cli,
        text: 'Use this path when referring to this secret in the CLI.',
      },
      {
        label: 'API path for metadata',
        snippet: namespace ? `/v1/${namespace}/${metadata}` : `/v1/${metadata}`,
        text: `Use this path when referring to this secret's metadata in the API and permanent secret deletion.`,
      },
    ];
  }

  get commands() {
    const cliPath = this.paths.findBy('label', 'CLI path').snippet;
    const apiPath = this.paths.findBy('label', 'API path').snippet;
    // as a future improvement, it might be nice to use window.location.protocol here:
    const url = `https://127.0.0.1:8200${apiPath}`;

    return {
      cli: `vault kv get ${cliPath}`,
      /* eslint-disable-next-line no-useless-escape */
      apiCopy: `curl \ --header "X-Vault-Token: ..." \ --request GET \ ${url}`,
      apiDisplay: `curl \\
        --header "X-Vault-Token: ..." \\
        --request GET \\
      ${url}`,
    };
  }
}

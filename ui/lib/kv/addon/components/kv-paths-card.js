/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { kvMetadataPath, kvDataPath } from 'vault/utils/kv-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

/**
 * @module KvPathsCard is used to display copyable secret paths for KV v2 for CLI and API use.
 * This component is permission agnostic because args come from the views mount path and url params.
 *
 * <KvPathsCard
 *  @path={{this.model.path}}
 *  @backend={{this.model.backend}}
 *  @isCondensed={{false}}
 * />
 *
 * @param {string} path - kv secret path for building the CLI and API paths
 * @param {string} backend - the secret engine mount path, comes from the secretMountPath service defined in the route
 * @param {boolean} isCondensed - if true a smaller version displays with no commands section or extra explanatory text
 */

export default class KvPathsCard extends Component {
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
        snippet: namespace ? `/v1/${encodePath(namespace)}/${data}` : `/v1/${data}`,
        text: 'Use this path when referring to this secret in the API.',
      },
      {
        label: 'CLI path',
        snippet: namespace ? `-namespace=${namespace} ${cli}` : cli,
        text: 'Use this path when referring to this secret in the CLI.',
      },
      ...(this.args.isCondensed
        ? []
        : [
            {
              label: 'API path for metadata',
              snippet: namespace ? `/v1/${encodePath(namespace)}/${metadata}` : `/v1/${metadata}`,
              text: `Use this path when referring to this secret's metadata in the API and permanent secret deletion.`,
            },
          ]),
    ];
  }

  get commands() {
    const cliPath = this.paths.find((p) => p.label === 'CLI path').snippet;
    const apiPath = this.paths.find((p) => p.label === 'API path').snippet;
    // as a future improvement, it might be nice to use window.location.protocol here:
    const url = `https://127.0.0.1:8200${apiPath}`;

    return {
      cli: `vault kv get ${cliPath}`,
      api: `curl \\
  --header "X-Vault-Token: ..." \\
  --request GET \\
  ${url}`,
    };
  }
}

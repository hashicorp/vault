/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields()
export default class PkiConfigClusterModel extends Model {
  // This model uses the backend value as the model ID
  get useOpenAPI() {
    return true;
  }

  getHelpUrl(backendPath) {
    return `/v1/${backendPath}/config/cluster?help=1`;
  }

  @attr('string', {
    label: "Mount's API path",
    subText:
      "Specifies the path to this performance replication cluster's API mount path, including any namespaces as path components. This address is used for the ACME directories, which must be served over a TLS-enabled listener.",
  })
  path;
  @attr('string', {
    label: 'AIA path',
    subText:
      "Specifies the path to this performance replication cluster's AIA distribution point; may refer to an external, non-Vault responder.",
  })
  aiaPath;

  // this is for pki-only cluster config, not the universal vault cluster
  @lazyCapabilities(apiPath`${'id'}/config/cluster`, 'id') clusterPath;

  get canSet() {
    return this.clusterPath.get('canUpdate') !== false;
  }
}

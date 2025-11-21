/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { PkiConfigureClusterRequest } from '@hashicorp/vault-client-typescript';

export default class PkiConfigClusterForm extends Form<PkiConfigureClusterRequest> {
  formFields = [
    new FormField('path', 'string', {
      label: "Mount's API path",
      subText:
        "Specifies the path to this performance replication cluster's API mount path, including any namespaces as path components. This address is used for the ACME directories, which must be served over a TLS-enabled listener.",
    }),
    new FormField('aia_path', 'string', {
      label: 'AIA path',
      subText:
        "Specifies the path to this performance replication cluster's AIA distribution point; may refer to an external, non-Vault responder.",
    }),
  ];
}

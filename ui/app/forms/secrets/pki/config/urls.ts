/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';

import type { PkiConfigureUrlsRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';

export default class PkiConfigUrlsForm extends OpenApiForm<PkiConfigureUrlsRequest> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('PkiConfigureUrlsRequest', ...args);
  }
}

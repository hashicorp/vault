/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';

import type { PkiTidyRequest, PkiConfigureAutoTidyRequest } from '@hashicorp/vault-client-typescript';

type PkiTidyFormRequest = PkiTidyRequest | PkiConfigureAutoTidyRequest;

export default class PkiTidyForm extends OpenApiForm<PkiTidyFormRequest> {
  constructor(...args: ConstructorParameters<typeof OpenApiForm>) {
    super(...args);

    // use ttl picker for pause_duration
    const pauseDuration = this.formFields.find((field) => field.name === 'pause_duration');
    if (pauseDuration) {
      pauseDuration.options.editType = 'ttl';
      this.data.pause_duration = '0';
    }
  }
}

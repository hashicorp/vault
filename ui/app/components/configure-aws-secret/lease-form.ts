/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { action } from '@ember/object';
import type SecretEngineModel from 'vault/models/secret-engine';
import type { TtlEvent } from 'vault/app-types';

type LeaseFields = { lease: string; lease_max: string };

interface Args {
  model: SecretEngineModel;
  onSubmit: (data: LeaseFields) => void;
}

export default class ConfigureAwsSecretLeaseFormComponent extends Component<Args> {
  @action
  handleTtlChange(name: string, ttlObj: TtlEvent) {
    // lease values cannot be undefined, set to 0 to use default
    const valueToSet = ttlObj.enabled ? ttlObj.goSafeTimeString : 0;
    this.args.model.set(name, valueToSet);
  }

  @action
  saveLease(data: LeaseFields, event: Event) {
    event.preventDefault();
    this.args.onSubmit(data);
  }
}

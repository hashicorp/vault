/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  lease: [
    {
      validator(model) {
        const { lease, leaseMax } = model;
        return (lease && leaseMax) || (!lease && !leaseMax) ? true : false;
      },
      message: 'Lease TTL and Max Lease TTL are both required if one of them is set.',
    },
  ],
};
@withModelValidations(validations)
export default class AwsLeaseConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr({
    label: 'Max Lease TTL',
    editType: 'ttl',
  })
  leaseMax;
  @attr({
    label: 'Default Lease TTL',
    editType: 'ttl',
  })
  lease;

  get attrs() {
    const keys = ['lease', 'leaseMax'];
    return expandAttributeMeta(this, keys);
  }
}

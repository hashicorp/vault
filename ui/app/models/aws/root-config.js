/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class AwsRootConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string') accessKey;
  @attr('string') region;
  @attr('string', {
    label: 'IAM endpoint',
  })
  iamEndpoint;
  @attr('string', {
    label: 'STS endpoint',
  })
  stsEndpoint;
  @attr('number', {
    defaultValue: -1,
    label: 'Maximum retries',
    subText: 'Number of max retries the client should use for recoverable errors. Default is -1.',
  })
  maxRetries;
  // TODO: there are more options available on the API, but the UI does not support them yet.
  get attrs() {
    const keys = ['accessKey', 'region', 'iamEndpoint', 'stsEndpoint', 'maxRetries'];
    return expandAttributeMeta(this, keys);
  }
}

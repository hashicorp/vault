/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Base } from '../credentials';
import { value, create, fillable, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  userIsPresent: isPresent('[data-test-input="username"]'),
  ipIsPresent: isPresent('[data-test-input="ip"]'),
  user: fillable('[data-test-input="username"]'),
  ip: fillable('[data-test-input="ip"]'),
  warningIsPresent: isPresent('[data-test-warning]'),
  commonNameValue: value('[data-test-input="commonName"]'),
  generateOTP: async function () {},
});

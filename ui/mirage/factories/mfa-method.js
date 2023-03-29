/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  type: 'okta',
  uses_passcode: false,

  afterCreate(mfaMethod) {
    if (mfaMethod.type === 'totp') {
      mfaMethod.uses_passcode = true;
    }
  },
});

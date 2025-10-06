/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  type: 'okta',
  uses_passcode: false,
  self_enrollment_enabled: false,

  afterCreate(mfaMethod) {
    if (mfaMethod.type === 'totp') {
      mfaMethod.uses_passcode = true;
    }
  },
});

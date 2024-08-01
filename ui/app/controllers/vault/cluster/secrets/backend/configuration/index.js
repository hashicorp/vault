/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { computed } from '@ember/object';
import Controller from '@ember/controller';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';

export default Controller.extend({
  isConfigurable: computed('model.type', function () {
    return CONFIGURABLE_SECRET_ENGINES.includes(this.model.type);
  }),
});

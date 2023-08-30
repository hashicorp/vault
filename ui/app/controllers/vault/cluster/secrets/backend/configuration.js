/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { computed } from '@ember/object';
import Controller from '@ember/controller';

export default Controller.extend({
  isConfigurable: computed('model.type', function () {
    const configurableEngines = ['aws', 'ssh'];
    return configurableEngines.includes(this.model.type);
  }),
});

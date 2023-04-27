/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { computed } from '@ember/object';
import Controller from '@ember/controller';

export default Controller.extend({
  isConfigurable: computed('model.type', function () {
    const configurableEngines = ['aws', 'ssh', 'pki'];
    return configurableEngines.includes(this.model.type);
  }),
});

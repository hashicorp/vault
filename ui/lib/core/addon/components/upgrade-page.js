/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/upgrade-page';

export default Component.extend({
  layout,
  title: 'Vault Enterprise',
  featureName: computed('title', function () {
    const title = this.title;
    return title === 'Vault Enterprise' ? 'this feature' : title;
  }),
  minimumEdition: 'Vault Enterprise',
});

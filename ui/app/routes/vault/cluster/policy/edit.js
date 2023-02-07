/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import ShowRoute from './show';
import { inject as service } from '@ember/service';

export default ShowRoute.extend(UnsavedModelRoute, {
  wizard: service(),

  activate() {
    if (this.wizard.featureState === 'details') {
      this.wizard.transitionFeatureMachine('details', 'CONTINUE', this.policyType());
    }
  },
});

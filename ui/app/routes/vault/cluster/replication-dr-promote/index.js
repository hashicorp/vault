/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Base from '../cluster-route-base';

export default Base.extend({
  replicationMode: service(),
  beforeModel() {
    this._super(...arguments);
    this.replicationMode.setMode('dr');
  },
});

/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Base from '../cluster-route-base';

export default Base.extend({
  replicationMode: service(),
  beforeModel() {
    this._super(...arguments);
    this.replicationMode.setMode('dr');
  },
});

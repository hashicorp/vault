/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { alias } from '@ember/object/computed';
import { service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
});

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Actions from './replication-actions-single';
import layout from '../templates/components/replication-action-update-primary';

export default Actions.extend({
  layout,
  tagName: '',
});

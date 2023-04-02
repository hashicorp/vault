/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  address: '127.0.0.1',
  node_id: (i) => `raft_node_${i}`,
  protocol_version: '3',
  voter: true,
  leader: true,
});

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';

//{"node_id":"1249bfbc-b234-96f3-0c66-07078ac3e16e","address":"127.0.0.1:8201","leader":true,"protocol_version":"3","voter":true}
export default Model.extend({
  address: attr('string'),
  nodeId: attr('string'),
  protocolVersion: attr('string'),
  voter: attr('boolean'),
  leader: attr('boolean'),
});

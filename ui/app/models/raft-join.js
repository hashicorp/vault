/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { computed } from '@ember/object';

//leader_api_addr (string: <required>) – Address of the leader node in the Raft cluster to which this node is trying to join.

//retry (bool: false) - Retry joining the Raft cluster in case of failures.

//leader_ca_cert (string: "") - CA certificate used to communicate with Raft's leader node.

//leader_client_cert (string: "") - Client certificate used to communicate with Raft's leader node.

//leader_client_key (string: "") - Client key used to communicate with Raft's leader node.

export default Model.extend({
  leaderApiAddr: attr('string', {
    label: 'Leader API Address',
  }),
  retry: attr('boolean', {
    label: 'Keep retrying to join in case of failures',
  }),
  leaderCaCert: attr('string', {
    label: 'Leader CA Certificate',
    editType: 'file',
  }),
  leaderClientCert: attr('string', {
    label: 'Leader Client Certificate',
    editType: 'file',
  }),
  leaderClientKey: attr('string', {
    label: 'Leader Client Key',
    editType: 'file',
  }),
  fields: computed(function () {
    return expandAttributeMeta(this, [
      'leaderApiAddr',
      'leaderCaCert',
      'leaderClientCert',
      'leaderClientKey',
      'retry',
    ]);
  }),
});

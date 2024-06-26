/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory, trait } from 'miragejs';

export default Factory.extend({
  auth: null,
  data: null, // populated via traits
  lease_duration: 0,
  lease_id: '',
  renewable: true,
  request_id: '22068a49-a504-41ad-b5b0-1eac71659190',
  warnings: null,
  wrap_info: null,

  // add servers to test raft storage configuration
  withRaft: trait({
    afterCreate(config, server) {
      if (!config.data) {
        config.data = {
          config: {
            index: 0,
            servers: server.serializerOrRegistry.serialize(server.createList('server', 2)),
          },
        };
      }
    },
  }),
});

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';

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

import { Factory, trait } from 'ember-cli-mirage';
import faker from 'faker';

export default Factory.extend({
  auth: null,
  data: null, // populated via traits
  lease_duration: 0,
  lease_id: '',
  renewable: () => faker.datatype.boolean(),
  request_id: () => faker.datatype.uuid(),
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

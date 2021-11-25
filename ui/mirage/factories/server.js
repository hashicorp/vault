import { Factory } from 'ember-cli-mirage';
import faker from 'faker';

export default Factory.extend({
  address: () => faker.internet.ip(),
  node_id: i => `raft_node_${i}`,
  protocol_version: '3',
  voter: () => faker.datatype.boolean(),
  leader: () => faker.datatype.boolean(),
});

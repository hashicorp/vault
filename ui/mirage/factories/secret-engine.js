import { Factory } from 'ember-cli-mirage';
import faker from 'faker';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';

export default Factory.extend({
  path: () => faker.system.directoryPath(),
  description: () => faker.git.commitMessage(),
  local: () => faker.datatype.boolean(),
  sealWrap: () => faker.datatype.boolean(),
  // set in afterCreate
  accessor: null,
  type: null,
  options: null,

  afterCreate(secretEngine) {
    if (!secretEngine.type) {
      const type = faker.random.arrayElement(supportedSecretBackends());
      secretEngine.type = type;

      if (!secretEngine.accessor) {
        secretEngine.accessor = `type_${faker.git.shortSha()}`;
      }
    }

    if (!secretEngine.options && ['generic', 'kv'].includes(secretEngine.type)) {
      secretEngine.options = {
        version: faker.random.arrayElement('1', '2'),
      };
    }
  },
});

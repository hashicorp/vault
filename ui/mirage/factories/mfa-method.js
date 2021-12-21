import { Factory } from 'ember-cli-mirage';
import faker from 'faker';

export default Factory.extend({
  type: () => faker.random.arrayElement('duo', 'okta', 'pingid', 'totp'),
  uses_passcode: false,

  afterCreate(mfaMethod) {
    if (mfaMethod.type === 'totp') {
      mfaMethod.uses_passcode = true;
    }
  },
});

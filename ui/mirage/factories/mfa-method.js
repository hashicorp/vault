import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  type: 'okta',
  uses_passcode: false,

  afterCreate(mfaMethod) {
    if (mfaMethod.type === 'totp') {
      mfaMethod.uses_passcode = true;
    }
  },
});

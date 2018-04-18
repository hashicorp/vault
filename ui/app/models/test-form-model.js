// this model is just used for integration tests
//

import AuthMethodModel from './auth-method';
import { fragment } from 'ember-data-model-fragments/attributes';

export default AuthMethodModel.extend({
  otherConfig: fragment('mount-config', { defaultValue: {} }),
});

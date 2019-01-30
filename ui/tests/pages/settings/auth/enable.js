import { create, visitable } from 'ember-cli-page-object';
import backendForm from '../../components/mount-backend-form';
import flashMessages from '../../components/flash-message';
import withFlash from 'vault/tests/helpers/with-flash';

export default create({
  visit: visitable('/vault/settings/auth/enable'),
  ...backendForm,
  flash: flashMessages,
  enable: async function(type, path) {
    await this.visit();
    return withFlash(this.mount(type, path));
  },
});

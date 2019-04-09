import { create, visitable } from 'ember-cli-page-object';
import backendForm from '../../components/mount-backend-form';
import flashMessages from '../../components/flash-message';

export default create({
  visit: visitable('/vault/settings/auth/enable'),
  ...backendForm,
  flash: flashMessages,
  enable: async function(type, path) {
    await this.visit();
    await this.mount(type, path);
  },
});

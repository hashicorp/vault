import { create, visitable } from 'ember-cli-page-object';
import backendForm from '../../components/mount-backend-form';
import flashMessages from '../../components/flash-message';

export default create({
  visit: visitable('/vault/settings/auth/enable'),
  form: backendForm,
  flash: flashMessages,
  enableAuth(type, path) {
    if (path) {
      return this.form.path(path).type(type).submit();
    }
    return this.form.type(type).submit();
  },
});

import { create, clickable, fillable, visitable, selectable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/list'),
  visitShow: visitable('/vault/secrets/:backend/show/:id'),
  visitCreate: visitable('/vault/secrets/:backend/create'),
  createLink: clickable('[data-test-secret-create="true"]'),
  dbPlugin: selectable('[data-test-input="plugin_name"]'),
  name: fillable('[data-test-input="name"]'),
  toggleVerify: clickable('[data-test-input="verify_connection"]'),
  url: fillable('[data-test-input="connection_url"'),
  save: clickable('[data-test-secret-save=""]'),
  enable: clickable('[data-test-enable-connection=""]'),
  edit: clickable('[data-test-edit-link="true"]'),
});

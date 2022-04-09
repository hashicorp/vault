import { create, clickable, fillable, visitable, selectable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/list?itemType=role'),
  visitShow: visitable('/vault/secrets/:backend/show/role/:id'),
  visitCreate: visitable('/vault/secrets/:backend/create?itemType=role'),
  createLink: clickable('[data-test-secret-create]'),
  name: fillable('[data-test-input="name"]'),
  roleType: selectable('[data-test-input="type"'),
  save: clickable('[data-test-secret-save]'),
  edit: clickable('[data-test-edit-link]'),
});

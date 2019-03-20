import { Base } from '../create';
import { clickable, visitable, create, fillable } from 'ember-cli-page-object';

export default create({
  ...Base,
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  toggleDomain: clickable('[data-test-toggle-group="Domain Handling"]'),
  toggleOptions: clickable('[data-test-toggle-group="Options"]'),
  name: fillable('[data-test-input="name"]'),
  allowAnyName: clickable('[data-test-input="allowAnyName"]'),
  allowedDomains: fillable('[data-test-input="allowedDomains"] input'),
  save: clickable('[data-test-role-create]'),
  deleteBtn: clickable('[data-test-role-delete] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  deleteRole() {
    return this.deleteBtn().confirmBtn();
  },

  createRole(name, allowedDomains) {
    return this.toggleDomain()
      .toggleOptions()
      .name(name)
      .allowAnyName()
      .allowedDomains(allowedDomains)
      .save();
  },
});

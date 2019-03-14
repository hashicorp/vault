import { Base } from '../create';
import confirmAction from 'vault/tests/pages/components/confirm-action';
import { clickable, visitable, create, fillable } from 'ember-cli-page-object';

export default create({
  ...Base,
  ...confirmAction,
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  toggleDomain: clickable('[data-test-toggle-group="Domain Handling"]'),
  toggleOptions: clickable('[data-test-toggle-group="Options"]'),
  name: fillable('[data-test-input="name"]'),
  allowAnyName: clickable('[data-test-input="allowAnyName"]'),
  allowedDomains: fillable('[data-test-input="allowedDomains"] input'),
  save: clickable('[data-test-role-create]'),
  async deleteRole() {
    await this.delete();
    await this.confirmDelete();
  },

  async createRole(name, allowedDomains) {
    await this.toggleDomain()
      .toggleOptions()
      .name(name)
      .allowAnyName()
      .allowedDomains(allowedDomains)
      .save();
  },
});

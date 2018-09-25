import { create, visitable } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/identity/edit-form';

export default create({
  visit: visitable('/vault/access/identity/:item_type/create'),
  editForm,
  createItem(item_type, type) {
    if (type) {
      return this.visit({ item_type })
        .editForm.type(type)
        .submit();
    }
    return this.visit({ item_type }).editForm.submit();
  },
});

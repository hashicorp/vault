import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
};

@withModelValidations(validations)
export default class OidcScopeModel extends Model {
  @attr('string', { editDisabled: true }) name;
  @attr('string', { editType: 'textarea' }) description;
  @attr('string', { label: 'JSON Template', editType: 'json', mode: 'ruby' }) template;

  // TODO refactor when field-to-attrs is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'description', 'template']);
    }
    return this._attributeMeta;
  }

  @lazyCapabilities(apiPath`identity/oidc/scope/${'name'}`, 'name') scopePath;
  get canRead() {
    return this.scopePath.get('canRead');
  }
  get canEdit() {
    return this.scopePath.get('canUpdate');
  }
  get canDelete() {
    return this.scopePath.get('canDelete');
  }
}

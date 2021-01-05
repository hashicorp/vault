import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
// import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import attachCapabilities from 'vault/lib/attach-capabilities';

const ModelExport = Model.extend({
  // used for getting appropriate options for backend
  idPrefix: 'role/',
  // the id prefixed with `role/` so we can use it as the *secret param for the secret show route
  idForNav: computed('id', 'idPrefix', function() {
    let modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }),

  name: attr('string', {
    label: 'Name',
    fieldValue: 'id',
    readOnly: true,
    subText: 'TODO add subtext',
  }),
  // TODO ARG SEE API DOCS
  // https://www.vaultproject.io/api-docs/secret/databases#create-role

  // transformations: attr('array', {
  //   editType: 'searchSelect',
  //   fallbackComponent: 'string-list',
  //   label: 'Transformations',
  //   models: ['transform'],
  //   onlyAllowExisting: true,
  //   subLabel: 'Transformations',
  //   subText: 'Select which transformations this role will have access to. It must already exist.',
  // }),

  // attrs: computed('transformations', function() {
  //   let keys = ['name', 'transformations'];
  //   return expandAttributeMeta(this, keys);
  // }),

  backend: attr('string', { readOnly: true }),
});

export default attachCapabilities(ModelExport, {
  // ARG TODO: configures a role
  updatePath: apiPath`${'backend'}/roles/${'id'}`,
});

import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import attachCapabilities from 'vault/lib/attach-capabilities';

const { attr } = DS;

// these arrays define the order in which the fields will be displayed
// see
//https://www.vaultproject.io/api-docs/secret/transform#create-update-transformation
const TYPES = [
  {
    value: 'fpe',
    displayName: 'Format Preserving Encryption (FPE)',
  },
  {
    value: 'masking',
    displayName: 'Masking',
  },
];

const TWEAK_SOURCE = [
  {
    value: 'supplied',
    displayName: 'supplied',
  },
  {
    value: 'generated',
    displayName: 'generated',
  },
  {
    value: 'internal',
    displayName: 'internal',
  },
];

// export default DS.Model.extend({
const Model = DS.Model.extend({
  // TODO: for now, commenting out openApi info, but keeping here just in case we end up using it.
  // useOpenAPI: true,
  // getHelpUrl: function(backend) {
  //   return `/v1/${backend}?help=1`;
  // },
  name: attr('string', {
    // TODO: make this required for making a transformation
    label: 'Name',
    fieldValue: 'id',
    readOnly: true,
  }),
  type: attr('string', {
    defaultValue: 'fpe',
    label: 'Type',
    possibleValues: TYPES,
    subText:
      'Vault provides two types of transformations: Format Preserving Encryption (FPE) is reversible, while Masking is not.',
  }),
  tweak_source: attr('string', {
    defaultValue: 'supplied',
    label: 'Tweak source',
    possibleValues: TWEAK_SOURCE,
    subText: `A tweak value is used when performing FPE transformations. This can be supplied, generated, or internal.`, // TODO: I do not include the link here.  Need to figure out the best way to approach this.
  }),
  masking_character: attr('string', {
    defaultValue: '*',
    label: 'Masking character',
    subText: 'Specify which character you’d like to mask your data.',
  }),
  template: attr('string', {
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    label: 'Template', // TODO: make this required for making a transformation
    models: ['transform/template'],
    subLabel: 'Template Name',
    subText:
      'Templates allow Vault to determine what and how to capture the value to be transformed. Type to use an existing template or create a new one.',
  }),
  templates: attr('array'), // TODO: remove once BE changes the returned property to a singular template on the GET request.
  allowed_roles: attr('string', {
    label: 'Allowed roles',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['transform/role'],
    subText: 'Search for an existing role, type a new role to create it, or use a wildcard (*).',
  }),
  transformAttrs: computed('type', function() {
    // TODO: group them into sections/groups.  Right now, we don't different between required and not required as we do by hiding options.
    // will default to design mocks on how to handle as it will likely be a different pattern using client-side validation, which we have not done before
    if (this.type === 'masking') {
      return ['name', 'type', 'masking_character', 'template', 'templates', 'allowed_roles'];
    }
    return ['name', 'type', 'tweak_source', 'template', 'templates', 'allowed_roles'];
  }),
  transformFieldAttrs: computed('transformAttrs', function() {
    return expandAttributeMeta(this, this.get('transformAttrs'));
  }),
  // zeroAddressPath: lazyCapabilities(apiPath`${'backend'}/config/zeroaddress`, 'backend'),
  // canEditZeroAddress: alias('zeroAddressPath.canUpdate'),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`transform/transformation/${'id'}`,
});

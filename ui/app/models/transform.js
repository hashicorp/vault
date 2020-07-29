import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

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

export default DS.Model.extend({
  // TODO: for now, commenting out openApi info, but keeping here just in case we end up using it.
  // useOpenAPI: true,
  // getHelpUrl: function(backend) {
  //   console.log(backend, 'Backend');
  //   return `/v1/${backend}?help=1`;
  // },
  name: attr('string', {
    label: 'Transformation Name',
  }),
  type: attr('string', {
    defaultValue: 'fpe',
    label: 'Type',
    possibleValues: TYPES,
  }),
  template: attr('string', {
    label: 'Template name',
  }),
  tweak_source: attr('string', {
    defaultValue: 'supplied',
    label: 'Tweak source',
    possibleValues: TWEAK_SOURCE,
  }),
  masking_character: attr('string', {
    label: 'Masking character',
  }),
  allowed_roles: attr('stringArray', {
    label: 'Allowed roles',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['transform/role'],
  }),
  transformAttrs: computed(function() {
    // TODO: group them into sections/groups.  Right now, we don't different between required and not required as we do by hiding options.
    // will default to design mocks on how to handle as it will likely be a different pattern using client-side validation, which we have not done before
    return ['name', 'type', 'template', 'tweak_source', 'masking_characters', 'allowed_roles'];
  }),
  transformFieldAttrs: computed('transformAttrs', function() {
    return expandAttributeMeta(this, this.get('transformAttrs'));
  }),
  updatePath: lazyCapabilities(apiPath`${'backend'}/transforms/${'id'}`, 'backend', 'id'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),

  generatePath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerate: alias('generatePath.canUpdate'),

  signPath: lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id'),
  canSign: alias('signPath.canUpdate'),

  zeroAddressPath: lazyCapabilities(apiPath`${'backend'}/config/zeroaddress`, 'backend'),
  canEditZeroAddress: alias('zeroAddressPath.canUpdate'),
});

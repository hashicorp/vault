import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

// these arrays define the order in which the fields will be displayed
// see
// https://github.com/hashicorp/vault/blob/master/builtin/logical/ssh/path_roles.go#L542 for list of fields for each key type
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
  // useOpenAPI: true,
  // getHelpUrl: function(backend) {
  //   // ARG TODO: this is the open api
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
  allowed_roles: attr('array', {
    defaultValue: function() {
      return [];
    },
    label: 'Allowed roles',
  }),
  transformAttrs: computed(function() {
    // return [{ default: ['name', 'type', 'template', 'tweak_source', 'masking_characters', 'allowed_roles'] }];
    // for looping over form field types
    // TODO: group them into sections
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

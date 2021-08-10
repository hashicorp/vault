import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { fragment } from 'ember-data-model-fragments/attributes';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { validator, buildValidations } from 'ember-cp-validations';

//identity will be managed separately and the inclusion
//of the system backend is an implementation detail
const LIST_EXCLUDED_BACKENDS = ['system', 'identity'];

const Validations = buildValidations({
  path: validator('presence', {
    presence: true,
    message: "Path can't be blank.",
  }),
  maxVersions: [
    validator('number', {
      allowString: false,
      integer: true,
      message: 'Maximum versions must be a number.',
    }),
    validator('length', {
      min: 1,
      max: 16,
      message: 'You cannot go over 16 characters.',
    }),
  ],
});

export default Model.extend(Validations, {
  path: attr('string'),
  accessor: attr('string'),
  name: attr('string'),
  type: attr('string', {
    label: 'Secret engine type',
  }),
  description: attr('string', {
    editType: 'textarea',
  }),
  config: fragment('mount-config', { defaultValue: {} }),
  options: fragment('mount-options', { defaultValue: {} }),
  local: attr('boolean', {
    helpText:
      'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
  }),
  sealWrap: attr('boolean', {
    helpText:
      'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For K/V mounts, all values will be seal wrapped.) This can only be specified at mount time.',
  }),
  // KV 2 additional config default options
  maxVersions: attr('number', {
    defaultValue: 10,
    label: 'Maximum Number of Versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted. This value applies to all keys, but a key’s metadata settings can overwrite this value.',
  }),
  casRequired: attr('boolean', {
    defaultValue: false,
    label: 'Require Check and Set',
    subText:
      'If checked all keys will require the cas parameter to be set on all write request. A key’s metadata settings can overwrite this value.',
  }),
  deleteVersionAfter: attr({
    defaultValue: 0,
    editType: 'ttl',
    label: 'Automate secret deletion',
    helperTextDisabled: 'A secret’s version must be manually deleted.',
    helperTextEnabled: 'Delete all new versions of this secret after',
  }),

  modelTypeForKV: computed('engineType', 'options.version', function() {
    let type = this.engineType;
    let version = this.options?.version;
    let modelType = 'secret';
    if ((type === 'kv' || type === 'generic') && version === 2) {
      modelType = 'secret-v2';
    }
    return modelType;
  }),

  isV2KV: computed.equal('modelTypeForKV', 'secret-v2'),

  formFields: computed('engineType', function() {
    let type = this.engineType;
    let fields = [
      'type',
      'path',
      'description',
      'accessor',
      'local',
      'sealWrap',
      'config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
    ];
    if (type === 'kv' || type === 'generic') {
      fields.push('options.{version}');
    }
    return fields;
  }),

  formFieldGroups: computed('engineType', function() {
    let type = this.engineType;
    let defaultGroup;
    // KV has specific config options it adds on the enable engine. https://www.vaultproject.io/api/secret/kv/kv-v2#configure-the-kv-engine
    if (type === 'kv') {
      defaultGroup = { default: ['path', 'maxVersions', 'casRequired', 'deleteVersionAfter'] };
    } else {
      defaultGroup = { default: ['path'] };
    }
    let optionsGroup = {
      'Method Options': [
        'description',
        'config.listingVisibility',
        'local',
        'sealWrap',
        'config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
      ],
    };
    if (type === 'kv' || type === 'generic') {
      optionsGroup['Method Options'].unshift('options.{version}');
    }
    if (type === 'database') {
      // For the Database Secret Engine we want to highlight the defaultLeaseTtl and maxLeaseTtl, removing them from the options object
      defaultGroup.default.push('config.{defaultLeaseTtl}', 'config.{maxLeaseTtl}');
      return [
        defaultGroup,
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.{auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
          ],
        },
      ];
    }
    return [defaultGroup, optionsGroup];
  }),

  attrs: computed('formFields', function() {
    return expandAttributeMeta(this, this.formFields);
  }),

  fieldGroups: computed('formFieldGroups', function() {
    return fieldToAttrs(this, this.formFieldGroups);
  }),

  // namespaces introduced types with a `ns_` prefix for built-in engines
  // so we need to strip that to normalize the type
  engineType: computed('type', function() {
    return (this.type || '').replace(/^ns_/, '');
  }),

  shouldIncludeInList: computed('engineType', function() {
    return !LIST_EXCLUDED_BACKENDS.includes(this.engineType);
  }),

  localDisplay: computed('local', function() {
    return this.local ? 'local' : 'replicated';
  }),

  // ssh specific ones
  privateKey: attr('string'),
  publicKey: attr('string'),
  generateSigningKey: attr('boolean', {
    defaultValue: true,
  }),

  saveCA(options) {
    if (this.type !== 'ssh') {
      return;
    }
    if (options.isDelete) {
      this.setProperties({
        privateKey: null,
        publicKey: null,
        generateSigningKey: false,
      });
    }
    return this.save({
      adapterOptions: {
        options: options,
        apiPath: 'config/ca',
        attrsToSend: ['privateKey', 'publicKey', 'generateSigningKey'],
      },
    });
  },

  saveZeroAddressConfig() {
    return this.save({
      adapterOptions: {
        adapterMethod: 'saveZeroAddressConfig',
      },
    });
  },

  // aws backend attrs
  lease: attr('string'),
  leaseMax: attr('string'),
});

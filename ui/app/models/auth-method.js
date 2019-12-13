import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import { fragment } from 'ember-data-model-fragments/attributes';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { memberAction } from 'ember-api-actions';

import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';

const { attr, hasMany } = DS;

let Model = DS.Model.extend({
  authConfigs: hasMany('auth-config', { polymorphic: true, inverse: 'backend', async: false }),
  path: attr('string'),
  accessor: attr('string'),
  name: attr('string'),
  type: attr('string'),
  // namespaces introduced types with a `ns_` prefix for built-in engines
  // so we need to strip that to normalize the type
  methodType: computed('type', function() {
    return this.get('type').replace(/^ns_/, '');
  }),
  description: attr('string', {
    editType: 'textarea',
  }),
  config: fragment('mount-config', { defaultValue: {} }),
  local: attr('boolean', {
    helpText:
      'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
  }),
  sealWrap: attr('boolean', {
    helpText:
      'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For K/V mounts, all values will be seal wrapped.) This can only be specified at mount time.',
  }),

  // used when the `auth` prefix is important,
  // currently only when setting perf mount filtering
  apiPath: computed('path', function() {
    return `auth/${this.get('path')}`;
  }),
  localDisplay: computed('local', function() {
    return this.get('local') ? 'local' : 'replicated';
  }),

  tuneAttrs: computed(function() {
    return expandAttributeMeta(this, [
      'description',
      'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,tokenType,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
    ]);
  }),

  //sys/mounts/auth/[auth-path]/tune.
  tune: memberAction({
    path: 'tune',
    type: 'post',
    urlType: 'updateRecord',
  }),

  formFields: computed(function() {
    return [
      'type',
      'path',
      'description',
      'accessor',
      'local',
      'sealWrap',
      'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,tokenType,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
    ];
  }),

  formFieldGroups: computed(function() {
    return [
      { default: ['path'] },
      {
        'Method Options': [
          'description',
          'config.listingVisibility',
          'local',
          'sealWrap',
          'config.{defaultLeaseTtl,maxLeaseTtl,tokenType,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
        ],
      },
    ];
  }),

  attrs: computed('formFields', function() {
    return expandAttributeMeta(this, this.get('formFields'));
  }),

  fieldGroups: computed('formFieldGroups', function() {
    return fieldToAttrs(this, this.get('formFieldGroups'));
  }),
  canDisable: alias('deletePath.canDelete'),
  canEdit: alias('configPath.canUpdate'),
});

export default attachCapabilities(Model, {
  deletePath: apiPath`sys/auth/${'id'}`,
  configPath: function(context) {
    if (context.type === 'aws') {
      return apiPath`auth/${'id'}/config/client`;
    } else {
      return apiPath`auth/${'id'}/config`;
    }
  },
});

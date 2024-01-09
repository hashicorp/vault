/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { belongsTo, hasMany, attr } from '@ember-data/model';
import { alias } from '@ember/object/computed'; // eslint-disable-line
import { computed } from '@ember/object'; // eslint-disable-line
import { inject as service } from '@ember/service';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  path: [
    { type: 'presence', message: "Path can't be blank." },
    {
      type: 'containsWhiteSpace',
      message:
        "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.",
      level: 'warn',
    },
  ],
};

// unsure if ember-api-actions will work on native JS class model
// for now create class to use validations and then use classic extend pattern
@withModelValidations(validations)
class AuthMethodModel extends Model {}
const ModelExport = AuthMethodModel.extend({
  store: service(),

  config: belongsTo('mount-config', { async: false, inverse: null }), // one-to-none that replaces former fragment
  authConfigs: hasMany('auth-config', { polymorphic: true, inverse: 'backend', async: false }),
  path: attr('string'),
  accessor: attr('string'),
  name: attr('string'),
  type: attr('string'),
  // namespaces introduced types with a `ns_` prefix for built-in engines
  // so we need to strip that to normalize the type
  methodType: computed('type', function () {
    return this.type.replace(/^ns_/, '');
  }),
  description: attr('string', {
    editType: 'textarea',
  }),
  local: attr('boolean', {
    helpText:
      'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
  }),
  sealWrap: attr('boolean', {
    helpText:
      'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
  }),

  // used when the `auth` prefix is important,
  // currently only when setting perf mount filtering
  apiPath: computed('path', function () {
    return `auth/${this.path}`;
  }),
  localDisplay: computed('local', function () {
    return this.local ? 'local' : 'replicated';
  }),

  tuneAttrs: computed('path', function () {
    const { methodType } = this;
    let tuneAttrs;
    // token_type should not be tuneable for the token auth method
    if (methodType === 'token') {
      tuneAttrs = [
        'description',
        'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
      ];
    } else {
      tuneAttrs = [
        'description',
        'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,tokenType,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
      ];
    }
    return expandAttributeMeta(this, tuneAttrs);
  }),

  formFields: computed(function () {
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

  formFieldGroups: computed(function () {
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

  attrs: computed('formFields', function () {
    return expandAttributeMeta(this, this.formFields);
  }),

  fieldGroups: computed('formFieldGroups', function () {
    return fieldToAttrs(this, this.formFieldGroups);
  }),
  canDisable: alias('deletePath.canDelete'),
  canEdit: alias('configPath.canUpdate'),

  tune(data) {
    return this.store.adapterFor('auth-method').tune(this.path, data);
  },
});

export default attachCapabilities(ModelExport, {
  deletePath: apiPath`sys/auth/${'id'}`,
  configPath: function (context) {
    if (context.type === 'aws') {
      return apiPath`auth/${'id'}/config/client`.call(this, context);
    } else {
      return apiPath`auth/${'id'}/config`.call(this, context);
    }
  },
});

import Ember from 'ember';
import DS from 'ember-data';
import { fragment } from 'ember-data-model-fragments/attributes';
import { queryRecord } from 'ember-computed-query';
import { methods } from 'vault/helpers/mountable-auth-methods';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { memberAction } from 'ember-api-actions';

const { attr, hasMany } = DS;
const { computed } = Ember;

const METHODS = methods();

const configPath = function configPath(strings, key) {
  return function(...values) {
    return `${strings[0]}${values[key]}${strings[1]}`;
  };
};
export default DS.Model.extend({
  authConfigs: hasMany('auth-config', { polymorphic: true, inverse: 'backend', async: false }),
  path: attr('string', {
    defaultValue: METHODS[0].value,
  }),
  accessor: attr('string'),
  name: attr('string'),
  type: attr('string', {
    defaultValue: METHODS[0].value,
    possibleValues: METHODS,
  }),
  description: attr('string', {
    editType: 'textarea',
  }),
  config: fragment('mount-config', { defaultValue: {} }),
  local: attr('boolean'),
  sealWrap: attr('boolean'),

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
      'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
    ]);
  }),

  //sys/mounts/auth/[auth-path]/tune.
  tune: memberAction({
    path: 'tune',
    type: 'post',
    urlType: 'updateRecord',
  }),

  formFields: [
    'type',
    'path',
    'description',
    'accessor',
    'local',
    'sealWrap',
    'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
  ],

  formFieldGroups: [
    { default: ['type', 'path'] },
    {
      'Method Options': [
        'description',
        'config.listingVisibility',
        'local',
        'sealWrap',
        'config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders}',
      ],
    },
  ],

  attrs: computed('formFields', function() {
    return expandAttributeMeta(this, this.get('formFields'));
  }),

  fieldGroups: computed('formFieldGroups', function() {
    return fieldToAttrs(this, this.get('formFieldGroups'));
  }),

  configPathTmpl: computed('type', function() {
    const type = this.get('type');
    if (type === 'aws') {
      return configPath`auth/${0}/config/client`;
    } else {
      return configPath`auth/${0}/config`;
    }
  }),

  configPath: queryRecord(
    'capabilities',
    context => {
      const { id, configPathTmpl } = context.getProperties('id', 'configPathTmpl');
      return {
        id: configPathTmpl(id),
      };
    },
    'id',
    'configPathTmpl'
  ),
  deletePath: queryRecord(
    'capabilities',
    context => {
      const { id } = context.get('id');
      return {
        id: `sys/auth/${id}`,
      };
    },
    'id'
  ),
  canDisable: computed.alias('deletePath.canDelete'),

  canEdit: computed.alias('configPath.canUpdate'),
});

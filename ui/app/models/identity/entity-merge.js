import { computed } from '@ember/object';
import DS from 'ember-data';
import IdentityModel from './_base';
const { attr } = DS;

export default IdentityModel.extend({
  formFields: computed(function() {
    return ['toEntityId', 'fromEntityIds', 'force'];
  }),
  toEntityId: attr('string', {
    label: 'Entity to merge to',
  }),
  fromEntityIds: attr({
    label: 'Entities to merge from',
    editType: 'stringArray',
  }),
  force: attr('boolean', {
    label: 'Keep MFA secrets from the "to" entity if there are merge conflicts',
    defaultValue: false,
  }),
});

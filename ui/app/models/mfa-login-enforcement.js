import Model, { attr, hasMany } from '@ember-data/model';
import { tracked } from '@glimmer/tracking';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class mfaLoginEnforcementModel extends Model {
  @attr('string', {
    label: 'Enforcement name',
    subText:
      'The name for this enforcement. Giving it a name means that you can refer to it again later. This name will not be editable later.',
  })
  name;
  @attr('string') namespace_id;
  @attr('array') auth_method_accessors; // JLR TODO: find out what this is and if there is a model for it
  @hasMany('auth-method') auth_method_types;
  @hasMany('identity/entity') identity_entity_ids;
  @hasMany('identity/group') identity_group_ids;
  // could be hasMany if global endpoint is created for fetching a method by id
  // currently the endpoints are scoped by type
  @attr('array') mfa_method_ids;
  // for now let's fetch the methods on demand
  @tracked methods = [];
  async fetchMethods() {
    if (this.mfa_method_ids.length) {
      try {
        // we should leverage the method adapter here
        this.methods = await this.store.adapterFor('mfa-method').query({ ids: this.mfa_method_ids });
      } catch (error) {
        this.methods = [];
        throw error;
      }
    }
  }

  // rather than returning an array of form fields handle each one separately
  get nameField() {
    return expandAttributeMeta(this, ['name'])[0];
  }
  // since methods is not an attribute we will fake it
  get methodsField() {
    return {
      name: 'methods',
      type: 'array',
      options: {
        label: 'MFA methods',
        subText: 'The MFA method(s) that this enforcement will apply to.',
        editType: 'searchSelect',
        models: ['mfa-method'],
      },
    };
  }
  // handle the targets field directly in the template since it is an aggregate of 4 attributes
}

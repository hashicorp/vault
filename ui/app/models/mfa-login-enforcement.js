import Model, { attr, hasMany } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class mfaLoginEnforcementModel extends Model {
  @attr('string', {
    label: 'Enforcement name',
    subText:
      'The name for this enforcement. Giving it a name means that you can refer to it again later. This name will not be editable later.',
  })
  name;
  @hasMany('mfa-method', {
    label: 'MFA methods',
    subText: 'The MFA method(s) that this enforcement will apply to.',
    editType: 'searchSelect',
    models: ['mfa-method'],
  })
  mfa_methods;
  @attr('string') namespace_id;
  @attr('array') auth_method_accessors; // ["auth_approle_17a552c6"]
  @attr('array') auth_method_types; // ["userpass"]
  @hasMany('identity/entity') identity_entities;
  @hasMany('identity/group') identity_groups;

  get formFields() {
    // handle the targets field directly in the template since it is an aggregate of 4 attributes
    return expandAttributeMeta(this, ['name', 'mfa_methods']);
  }
}

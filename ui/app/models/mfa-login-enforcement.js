import Model, { attr, hasMany } from '@ember-data/model';

export default class MfaLoginEnforcementModel extends Model {
  @attr('string') name;
  @hasMany('mfa-method') mfa_methods;
  @attr('string') namespace_id;
  @attr('array', { defaultValue: () => [] }) auth_method_accessors; // ["auth_approle_17a552c6"]
  @attr('array', { defaultValue: () => [] }) auth_method_types; // ["userpass"]
  @hasMany('identity/entity') identity_entities;
  @hasMany('identity/group') identity_groups;
}

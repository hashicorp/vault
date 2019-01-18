import { alias } from '@ember/object/computed';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { attr, belongsTo, hasMany } = DS;

export default DS.Model.extend({
  approved: attr('boolean'),
  requestPath: attr('string'),
  requestEntity: belongsTo('identity/entity', { async: false }),
  authorizations: hasMany('identity/entity', { async: false }),

  authorizePath: lazyCapabilities(apiPath`sys/control-group/authorize`),
  canAuthorize: alias('authorizePath.canUpdate'),
  configurePath: lazyCapabilities(apiPath`sys/config/control-group`),
  canConfigure: alias('configurePath.canUpdate'),
});

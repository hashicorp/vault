import Ember from 'ember';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { attr, belongsTo, hasMany } = DS;
const { computed } = Ember;

export default DS.Model.extend({
  approved: attr('boolean'),
  requestPath: attr('string'),
  requestEntity: belongsTo('identity/entity', { async: false }),
  authorizations: hasMany('identity/entity', { async: false }),

  authorizePath: lazyCapabilities(apiPath`sys/control-group/authorize`),
  canAuthorize: computed.alias('authorizePath.canUpdate'),
  configurePath: lazyCapabilities(apiPath`sys/config/control-group`),
  canConfigure: computed.alias('configurePath.canUpdate'),
});

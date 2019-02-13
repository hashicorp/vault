import DS from 'ember-data';
const { belongsTo } = DS;

export default DS.Model.extend({
  backend: belongsTo('auth-method', { inverse: 'authConfigs', readOnly: true, async: false }),
  getOpenApiInfo: function(backend) {
    return {
      helpUrl: `/v1/auth/${backend}/config?help=1`,
      path: `/config`,
    };
  },
});

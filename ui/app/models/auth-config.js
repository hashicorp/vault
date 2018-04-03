import DS from 'ember-data';
const { belongsTo } = DS;

export default DS.Model.extend({
  backend: belongsTo('auth-method', { readOnly: true, async: false }),
});

import DS from 'ember-data';
const { attr, belongsTo } = DS;

export default DS.Model.extend({
  config: belongsTo('kmip/config', { async: false }),
  caPem: attr('string', {
    label: 'CA PEM',
  }),
});

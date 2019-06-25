import DS from 'ember-data';
const { attr } = DS;

export default DS.Model.extend({
  caPem: attr('string', {
    label: 'CA PEM',
  }),
});

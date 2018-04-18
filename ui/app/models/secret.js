import Ember from 'ember';
import DS from 'ember-data';
import KeyMixin from './key-mixin';
const { attr } = DS;
const { computed } = Ember;

export default DS.Model.extend(KeyMixin, {
  auth: attr('string'),
  lease_duration: attr('number'),
  lease_id: attr('string'),
  renewable: attr('boolean'),

  secretData: attr('object'),

  dataAsJSONString: computed('secretData', function() {
    return JSON.stringify(this.get('secretData'), null, 2);
  }),

  isAdvancedFormat: computed('secretData', function() {
    const data = this.get('secretData');
    return Object.keys(data).some(key => typeof data[key] !== 'string');
  }),

  helpText: attr('string'),
  backend: attr('string'),
});

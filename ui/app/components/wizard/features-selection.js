import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  allFeatures: ['secrets', 'authentication', 'policies', 'replication', 'tools'],
  selectedFeatures: null,
  hasFeatures: computed('selectedFeatures', function() {
    return this.get('selectedFeatures') !== null && this.get('selectedFeatures').length > 0;
  }),
  finalFeatures: computed('allFeatures', 'selectedFeatures', function() {
    let features = [];
    let selected = this.get('selectedFeatures');
    this.get('allFeatures').forEach(function(feature) {
      if (selected.includes(feature)) {
        features.push(feature);
      }
    });
    return features;
  }),

  actions: {
    saveFeatures() {
      this.get('wizard').saveFeatures(this.get('finalFeatures'));
      this.get('wizard').transitionTutorialMachine('active.select', 'CONTINUE');
    },
    toggleFeature(event) {
      if (this.get('selectedFeatures') === null) {
        this.set('selectedFeatures', [event.target.value]);
      } else {
        if (this.get('selectedFeatures').includes(event.target.value)) {
          let selected = this.get('selectedFeatures').without(event.target.value);
          this.set('selectedFeatures', selected);
        } else {
          this.set('selectedFeatures', this.get('selectedFeatures').toArray().addObject(event.target.value));
        }
      }
    },
  },
});

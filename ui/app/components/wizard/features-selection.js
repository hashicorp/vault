import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  allFeatures: [
    {
      key: 'secrets',
      name: 'Secrets',
      steps: ['Enabling a secrets engine', 'Entering secrets method details', 'Adding a secret'],
    },
    {
      key: 'authentication',
      name: 'Authentication',
      steps: ['Enabling an auth method', 'Entering auth method details', 'Adding a user'],
    },
    {
      key: 'policies',
      name: 'Policies',
      steps: [],
    },
    {
      key: 'replication',
      name: 'Replication',
      steps: [],
    },
    {
      key: 'tools',
      name: 'Tools',
      steps: [],
    },
  ],
  selectedFeatures: null,
  hasFeatures: computed('selectedFeatures', function() {
    return this.get('selectedFeatures') !== null && this.get('selectedFeatures').length > 0;
  }),
  finalFeatures: computed('selectedFeatures', function() {
    let features = [];
    let selected = this.get('selectedFeatures');
    this.get('allFeatures').forEach(function(feature) {
      if (selected !== null && selected.includes(feature.key)) {
        features.push(feature.key);
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

import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  allFeatures: computed(function() {
    return [
      {
        key: 'secrets',
        name: 'Secrets',
        steps: ['Enabling a secrets engine', 'Entering secrets method details', 'Adding a secret'],
        selected: false,
      },
      {
        key: 'authentication',
        name: 'Authentication',
        steps: ['Enabling an auth method', 'Entering auth method details', 'Adding a user'],
        selected: false,
      },
      {
        key: 'policies',
        name: 'Policies',
        steps: [],
        selected: false,
      },
      {
        key: 'replication',
        name: 'Replication',
        steps: [],
        selected: false,
      },
      {
        key: 'tools',
        name: 'Tools',
        steps: [],
        selected: false,
      },
    ];
  }),

  selectedFeatures: computed('allFeatures.@each.selected', function() {
    return this.get('allFeatures').filterBy('selected').mapBy('key');
  }),

  actions: {
    saveFeatures() {
      let wizard = this.get('wizard');
      wizard.saveFeatures(this.get('selectedFeatures'));
      wizard.transitionTutorialMachine('active.select', 'CONTINUE');
    },
  },
});

import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  featureState: computed.alias('wizard.featureState'),
  secretType: computed.alias('wizard.componentState'),
  fullNextStep: computed.alias('wizard.nextStep'),
  nextStep: computed('fullNextStep', function() {
    return this.get('fullNextStep').split('.').lastObject;
  }),
  stepComponent: computed.alias('wizard.stepComponent'),
  detailsComponent: computed('secretType', function() {
    return this.get('secretType') ? `wizard/${this.get('secretType')}-engine` : 'wizard/ad-engine';
  }),
  onAdvance() {},
  onRepeat() {},
  onReset() {},
});

import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  wizard: inject.service(),
  featureState: computed.alias('wizard.featureState'),
  currentState: computed.alias('wizard.currentState'),
  currentMachine: computed.alias('wizard.currentMachine'),
  mountSubtype: computed.alias('wizard.componentState'),
  fullNextStep: computed.alias('wizard.nextStep'),
  nextFeature: computed.alias('wizard.nextFeature'),
  nextStep: computed('fullNextStep', function() {
    return this.get('fullNextStep').split('.').lastObject;
  }),
  stepComponent: computed.alias('wizard.stepComponent'),
  detailsComponent: computed('mountSubtype', function() {
    let suffix = this.get('currentState').includes('secret') ? 'engine' : 'mount';
    return this.get('mountSubtype') ? `wizard/${this.get('mountSubtype')}-${suffix}` : null;
  }),
  onAdvance() {},
  onRepeat() {},
  onReset() {},
});

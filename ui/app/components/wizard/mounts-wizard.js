import Ember from 'ember';
import { engines } from 'vault/helpers/mountable-secret-engines';
import { methods } from 'vault/helpers/mountable-auth-methods';
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
    let suffix = this.get('currentMachine') === 'secrets' ? 'engine' : 'method';
    return this.get('mountSubtype') ? `wizard/${this.get('mountSubtype')}-${suffix}` : null;
  }),
  mountName: computed('mountSubtype', function() {
    if (this.get('currentMachine') === 'secrets') {
      var secret = engines().find(engine => {
        return engine.type === this.get('mountSubtype');
      });
      if (secret) {
        return secret.displayName;
      }
    } else {
      var auth = methods().find(method => {
        return method.type === this.get('mountSubtype');
      });
      if (auth) {
        return auth.displayName;
      }
    }
    return null;
  }),
  actionText: computed('mountSubtype', function() {
    switch (this.get('mountSubtype')) {
      case 'aws':
        return 'Generate Credential';
      case 'ssh':
        return 'Sign Keys';
      case 'pki':
        return 'Generate Certificate';
      default:
        return null;
    }
  }),

  onAdvance() {},
  onRepeat() {},
  onReset() {},
});

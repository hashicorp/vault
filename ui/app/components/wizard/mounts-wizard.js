import { inject as service } from '@ember/service';
import { alias, equal } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { engines } from 'vault/helpers/mountable-secret-engines';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
const supportedSecrets = supportedSecretBackends();
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
const supportedAuth = supportedAuthBackends();

export default Component.extend({
  wizard: service(),
  featureState: alias('wizard.featureState'),
  currentState: alias('wizard.currentState'),
  currentMachine: alias('wizard.currentMachine'),
  mountSubtype: alias('wizard.componentState'),
  fullNextStep: alias('wizard.nextStep'),
  nextFeature: alias('wizard.nextFeature'),
  nextStep: computed('fullNextStep', function() {
    return this.fullNextStep.split('.').lastObject;
  }),
  needsEncryption: equal('mountSubtype', 'transit'),
  meep: equal('mountSubtype', 'database'),
  stepComponent: alias('wizard.stepComponent'),
  detailsComponent: computed('currentMachine', 'mountSubtype', function() {
    let suffix = this.currentMachine === 'secrets' ? 'engine' : 'method';
    return this.mountSubtype ? `wizard/${this.mountSubtype}-${suffix}` : null;
  }),
  isSupported: computed('currentMachine', 'mountSubtype', function() {
    if (this.currentMachine === 'secrets') {
      return supportedSecrets.includes(this.mountSubtype);
    } else {
      return supportedAuth.includes(this.mountSubtype);
    }
  }),
  mountName: computed('currentMachine', 'mountSubtype', function() {
    if (this.currentMachine === 'secrets') {
      var secret = engines().find(engine => {
        return engine.type === this.mountSubtype;
      });
      if (secret) {
        return secret.displayName;
      }
    } else {
      var auth = methods().find(method => {
        return method.type === this.mountSubtype;
      });
      if (auth) {
        return auth.displayName;
      }
    }
    return null;
  }),
  actionText: computed('mountSubtype', function() {
    switch (this.mountSubtype) {
      case 'aws':
        return 'Generate credential';
      case 'ssh':
        return 'Sign keys';
      case 'pki':
        return 'Generate certificate';
      default:
        return null;
    }
  }),

  onAdvance() {},
  onRepeat() {},
  onReset() {},
  onDone() {},
});

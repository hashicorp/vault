import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
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
    return this.get('fullNextStep').split('.').lastObject;
  }),
  needsEncryption: computed('mountSubtype', function() {
    return this.get('mountSubtype') === 'transit';
  }),
  stepComponent: alias('wizard.stepComponent'),
  detailsComponent: computed('mountSubtype', function() {
    let suffix = this.get('currentMachine') === 'secrets' ? 'engine' : 'method';
    return this.get('mountSubtype') ? `wizard/${this.get('mountSubtype')}-${suffix}` : null;
  }),
  isSupported: computed('mountSubtype', function() {
    if (this.get('currentMachine') === 'secrets') {
      return supportedSecrets.includes(this.get('mountSubtype'));
    } else {
      return supportedAuth.includes(this.get('mountSubtype'));
    }
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

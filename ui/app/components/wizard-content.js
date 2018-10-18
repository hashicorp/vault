import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { FEATURE_MACHINE_STEPS, INIT_STEPS } from 'vault/helpers/wizard-constants';

export default Component.extend({
  wizard: service(),
  classNames: ['ui-wizard'],
  glyph: null,
  headerText: null,
  selectProgress: null,
  currentMachine: computed.alias('wizard.currentMachine'),
  tutorialState: computed.alias('wizard.currentState'),
  tutorialComponent: computed.alias('wizard.tutorialComponent'),
  showProgress: computed('wizard.featureComponent', 'tutorialState', function() {
    return (
      this.tutorialComponent.includes('active') &&
      (this.tutorialState.includes('init.active') ||
        (this.wizard.featureComponent && this.wizard.featureMachineHistory))
    );
  }),
  featureMachineHistory: computed.alias('wizard.featureMachineHistory'),
  totalFeatures: computed('wizard.featureList', function() {
    return this.wizard.featureList.length;
  }),
  completedFeatures: computed('wizard.currentMachine', function() {
    return this.wizard.getCompletedFeatures();
  }),
  currentFeatureProgress: computed('featureMachineHistory.[]', function() {
    if (this.tutorialState.includes('active.feature')) {
      let totalSteps = FEATURE_MACHINE_STEPS[this.currentMachine];
      if (this.currentMachine === 'secrets') {
        if (this.featureMachineHistory.includes('secret')) {
          totalSteps = totalSteps['secret']['secret'];
        }
        if (this.featureMachineHistory.includes('list')) {
          totalSteps = totalSteps['secret']['list'];
        }
        if (this.featureMachineHistory.includes('encryption')) {
          totalSteps = totalSteps['encryption'];
        }
        if (this.featureMachineHistory.includes('role') || typeof totalSteps === 'object') {
          totalSteps = totalSteps['role'];
        }
      }
      return {
        percentage: (this.featureMachineHistory.length / totalSteps) * 100,
        feature: this.currentMachine,
        text: `Step ${this.featureMachineHistory.length} of ${totalSteps}`,
      };
    }
    return null;
  }),
  currentTutorialProgress: computed('tutorialState', function() {
    if (this.tutorialState.includes('init.active')) {
      let currentStepName = this.tutorialState.split('.')[2];
      let currentStepNumber = INIT_STEPS.indexOf(currentStepName) + 1;
      return {
        percentage: (currentStepNumber / INIT_STEPS.length) * 100,
        text: `Step ${currentStepNumber} of ${INIT_STEPS.length}`,
      };
    }
    return null;
  }),
  progressBar: computed('currentFeatureProgress', 'currentFeature', 'currentTutorialProgress', function() {
    let bar = [];
    if (this.currentTutorialProgress) {
      bar.push({
        style: `width:${this.currentTutorialProgress.percentage}%;`,
        completed: false,
        showIcon: true,
      });
    } else {
      if (this.currentFeatureProgress) {
        this.completedFeatures.forEach(feature => {
          bar.push({ style: 'width:100%;', completed: true, feature: feature, showIcon: true });
        });
        this.wizard.featureList.forEach(feature => {
          if (feature === this.currentMachine) {
            bar.push({
              style: `width:${this.currentFeatureProgress.percentage}%;`,
              completed: this.currentFeatureProgress.percentage == 100 ? true : false,
              feature: feature,
              showIcon: true,
            });
          } else {
            bar.push({ style: 'width:0%;', completed: false, feature: feature, showIcon: true });
          }
        });
      }
    }
    return bar;
  }),

  actions: {
    dismissWizard() {
      this.wizard.transitionTutorialMachine(this.wizard.currentState, 'DISMISS');
    },
  },
});

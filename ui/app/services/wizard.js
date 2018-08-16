import Ember from 'ember';
import { Machine } from 'xstate';

const { Service, inject } = Ember;

import getStorage from 'vault/lib/token-storage';

import TutorialMachineConfig from 'vault/machines/tutorial-machine';
import SecretsMachineConfig from 'vault/machines/secrets-machine';

const TutorialMachine = Machine(TutorialMachineConfig);
let FeatureMachine = null;
const TUTORIAL_STATE = 'vault-tutorial-state';
const FEATURE_LIST = 'vault-feature-list';
const FEATURE_STATE = 'vault-feature-state';
const COMPLETED_FEATURES = 'vault-completed-list';
const MACHINES = { secrets: SecretsMachineConfig };

export default Service.extend({
  router: inject.service(),
  currentState: null,
  featureList: null,
  featureState: null,
  currentMachine: null,
  potentialSelection: null,
  componentState: null,

  init() {
    this._super(...arguments);
    if (!this.storageHasKey(TUTORIAL_STATE)) {
      this.saveState('currentState', TutorialMachine.initialState);
      this.saveExtState(TUTORIAL_STATE, this.get('currentState'));
    } else {
      this.saveState('currentState', this.getExtState(TUTORIAL_STATE));
      if (this.storageHasKey(FEATURE_LIST)) {
        this.set('featureList', this.getExtState(FEATURE_LIST));
        if (this.storageHasKey(FEATURE_STATE)) {
          this.saveState('featureState', this.getExtState(FEATURE_STATE));
        } else {
          if (FeatureMachine != null) {
            this.saveState('featureState', FeatureMachine.initialState);
            this.saveExtState(FEATURE_STATE, this.get('featureState'));
          }
        }
        this.buildFeatureMachine();
      }
    }
  },

  saveState(stateType, state) {
    if (state.value) {
      state = state.value;
    }
    let stateKey = '';
    while (Ember.typeOf(state) === 'object') {
      let newState = Object.keys(state);
      stateKey += newState + '.';
      state = state[newState];
    }
    stateKey += state;
    this.set(stateType, stateKey);
  },

  transitionTutorialMachine(currentState, event, extendedState) {
    debugger;
    if (extendedState) {
      this.set('componentState', extendedState);
    }
    let { actions, value } = TutorialMachine.transition(currentState, event);
    this.saveState('currentState', value);
    this.saveExtState(TUTORIAL_STATE, this.get('currentState'));
    for (let action of actions) {
      this.executeAction(action, event);
    }
  },

  transitionFeatureMachine(currentState, event, extendedState) {
    let { actions, value } = FeatureMachine.transition(currentState, event, extendedState);
    this.saveState('featureState', value);
    this.saveExtState(FEATURE_STATE, value);
    for (let action of actions) {
      this.executeAction(action, event);
    }
  },

  saveExtState(key, value) {
    this.storage().setItem(key, value);
  },

  getExtState(key) {
    return this.storage().getItem(key);
  },

  storageHasKey(key) {
    return Boolean(this.getExtState(key));
  },

  executeAction(action, event) {
    let type = action;
    if (action.type) {
      type = action.type;
    }
    switch (type) {
      case 'routeTransition':
        this.get('router').transitionTo(...action.params);
        break;
      case 'saveFeatures':
        this.saveFeatures(event.features);
        break;
      case 'completeFeature':
        this.completeFeature();
        break;
      case 'handleDismissed':
        this.handleDismissed();
        break;
      default:
        break;
    }
  },

  handleDismissed() {
    this.storage().removeItem(FEATURE_STATE);
    this.storage().removeItem(FEATURE_LIST);
    this.storage().removeItem(MACHINES);
    this.storage().removeItem(COMPLETED_FEATURES);
  },

  saveFeatures(features) {
    this.set('featureList', features);
    this.saveExtState(FEATURE_LIST, this.get('featureList'));
    this.buildFeatureMachine();
  },

  buildFeatureMachine() {
    if (this.get('featureList') === null) {
      return;
    }
    const FeatureMachineConfig = MACHINES[this.get('featureList').objectAt(0)];
    FeatureMachine = Machine(FeatureMachineConfig);
    this.set('currentMachine', this.get('featureList').objectAt(0));
    this.saveState('featureState', FeatureMachine.initialState);
    this.saveExtState(FEATURE_STATE, this.get('featureState'));
    for (let action of FeatureMachine.initialState.actions) {
      this.executeAction(action);
    }
  },

  completeFeature() {
    let features = this.get('featureList');
    let completed = features.pop();

    if (!this.getExtState(COMPLETED_FEATURES)) {
      this.saveExtState([completed]);
    } else {
      this.saveExtState(
        COMPLETED_FEATURES,
        this.getExtState(COMPLETED_FEATURES).toArray().addObject(completed)
      );
    }

    this.saveExtState(FEATURE_LIST, this.get('featureList'));
    if (features.length > 0) {
      const FeatureMachineConfig = MACHINES[this.get('featureList').objectAt(0)];
      FeatureMachine = Machine(FeatureMachineConfig);
      this.set('currentMachine', features.objectAt(0));
      this.saveState('featureState', FeatureMachine.initialState);
      this.saveExtState(FEATURE_STATE, this.get('featureState'));
    } else {
      this.completeTutorial();
      FeatureMachine = null;
      TutorialMachine.transition(this.get('currentState'), 'DONE');
    }
  },

  setPotentialSelection(key) {
    this.set('potentialSelection', key);
  },

  storage() {
    return getStorage();
  },
});

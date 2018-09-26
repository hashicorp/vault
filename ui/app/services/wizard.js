import { next } from '@ember/runloop';
import { typeOf } from '@ember/utils';
import Service, { inject as service } from '@ember/service';
import { Machine } from 'xstate';

import getStorage from 'vault/lib/token-storage';

import TutorialMachineConfig from 'vault/machines/tutorial-machine';
import SecretsMachineConfig from 'vault/machines/secrets-machine';
import PoliciesMachineConfig from 'vault/machines/policies-machine';
import ReplicationMachineConfig from 'vault/machines/replication-machine';
import ToolsMachineConfig from 'vault/machines/tools-machine';
import AuthMachineConfig from 'vault/machines/auth-machine';

const TutorialMachine = Machine(TutorialMachineConfig);
let FeatureMachine = null;
const TUTORIAL_STATE = 'vault:ui-tutorial-state';
const FEATURE_LIST = 'vault:ui-feature-list';
const FEATURE_STATE = 'vault:ui-feature-state';
const COMPLETED_FEATURES = 'vault:ui-completed-list';
const COMPONENT_STATE = 'vault:ui-component-state';
const RESUME_URL = 'vault:ui-tutorial-resume-url';
const RESUME_ROUTE = 'vault:ui-tutorial-resume-route';
const MACHINES = {
  secrets: SecretsMachineConfig,
  policies: PoliciesMachineConfig,
  replication: ReplicationMachineConfig,
  tools: ToolsMachineConfig,
  authentication: AuthMachineConfig,
};

const DEFAULTS = {
  currentState: null,
  featureList: null,
  featureState: null,
  currentMachine: null,
  tutorialComponent: null,
  featureComponent: null,
  stepComponent: null,
  detailsComponent: null,
  componentState: null,
  nextFeature: null,
  nextStep: null,
};

export default Service.extend(DEFAULTS, {
  router: service(),
  showWhenUnauthenticated: false,

  init() {
    this._super(...arguments);
    this.initializeMachines();
  },

  initializeMachines() {
    if (!this.storageHasKey(TUTORIAL_STATE)) {
      let state = TutorialMachine.initialState;
      this.saveState('currentState', state.value);
      this.saveExtState(TUTORIAL_STATE, state.value);
    }
    this.saveState('currentState', this.getExtState(TUTORIAL_STATE));
    if (this.storageHasKey(COMPONENT_STATE)) {
      this.set('componentState', this.getExtState(COMPONENT_STATE));
    }
    let stateNodes = TutorialMachine.getStateNodes(this.get('currentState'));
    this.executeActions(stateNodes.reduce((acc, node) => acc.concat(node.onEntry), []), null, 'tutorial');
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
  },

  restartGuide() {
    let storage = this.storage();
    // empty storage
    [
      TUTORIAL_STATE,
      FEATURE_LIST,
      FEATURE_STATE,
      COMPLETED_FEATURES,
      COMPONENT_STATE,
      RESUME_URL,
      RESUME_ROUTE,
    ].forEach(key => storage.removeItem(key));
    // reset wizard state
    this.setProperties(DEFAULTS);
    // restart machines from blank state
    this.initializeMachines();
    // progress machine to 'active.select'
    this.transitionTutorialMachine('idle', 'AUTH');
  },

  saveState(stateType, state) {
    if (state.value) {
      state = state.value;
    }
    let stateKey = '';
    while (typeOf(state) === 'object') {
      let newState = Object.keys(state);
      stateKey += newState + '.';
      state = state[newState];
    }
    stateKey += state;
    this.set(stateType, stateKey);
  },

  transitionTutorialMachine(currentState, event, extendedState) {
    if (extendedState) {
      this.set('componentState', extendedState);
      this.saveExtState(COMPONENT_STATE, extendedState);
    }
    let { actions, value } = TutorialMachine.transition(currentState, event);
    this.saveState('currentState', value);
    this.saveExtState(TUTORIAL_STATE, this.get('currentState'));
    this.executeActions(actions, event, 'tutorial');
  },

  transitionFeatureMachine(currentState, event, extendedState) {
    if (!FeatureMachine || !this.get('currentState').includes('active')) {
      return;
    }
    if (extendedState) {
      this.set('componentState', extendedState);
      this.saveExtState(COMPONENT_STATE, extendedState);
    }

    let { actions, value } = FeatureMachine.transition(currentState, event, this.get('componentState'));
    this.saveState('featureState', value);
    this.saveExtState(FEATURE_STATE, value);
    this.executeActions(actions, event, 'feature');
    // if all features were completed, the FeatureMachine gets nulled
    // out and won't exist here as there is no next step
    if (FeatureMachine) {
      let next;
      if (this.get('currentMachine') === 'secrets' && value === 'display') {
        next = FeatureMachine.transition(value, 'REPEAT', this.get('componentState'));
      } else {
        next = FeatureMachine.transition(value, 'CONTINUE', this.get('componentState'));
      }
      this.saveState('nextStep', next.value);
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

  executeActions(actions, event, machineType) {
    let transitionURL;
    let expectedRouteName;
    let router = this.get('router');

    for (let action of actions) {
      let type = action;
      if (action.type) {
        type = action.type;
      }
      switch (type) {
        case 'render':
          this.set(`${action.level}Component`, action.component);
          break;
        case 'routeTransition':
          expectedRouteName = action.params[0];
          transitionURL = router.urlFor(...action.params).replace(/^\/ui/, '');
          next(() => {
            router.transitionTo(...action.params);
          });
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
        case 'handlePaused':
          this.handlePaused();
          return;
        case 'handleResume':
          this.handleResume();
          break;
        case 'showTutorialWhenAuthenticated':
          this.set('showWhenUnauthenticated', false);
          break;
        case 'showTutorialAlways':
          this.set('showWhenUnauthenticated', true);
          break;
        case 'continueFeature':
          this.transitionFeatureMachine(this.get('featureState'), 'CONTINUE', this.get('componentState'));
          break;
        default:
          break;
      }
    }
    if (machineType === 'tutorial') {
      return;
    }
    // if we're transitioning in the actions, we want that url,
    // else we want the URL we land on in didTransition in the
    // application route - we'll notify the application route to
    // update the route
    if (transitionURL) {
      this.set('expectedURL', transitionURL);
      this.set('expectedRouteName', expectedRouteName);
      this.set('setURLAfterTransition', false);
    } else {
      this.set('setURLAfterTransition', true);
    }
  },

  handlePaused() {
    let expected = this.get('expectedURL');
    if (expected) {
      this.saveExtState(RESUME_URL, this.get('expectedURL'));
      this.saveExtState(RESUME_ROUTE, this.get('expectedRouteName'));
    }
  },

  handleResume() {
    let resumeURL = this.storage().getItem(RESUME_URL);
    if (!resumeURL) {
      return;
    }
    this.get('router').transitionTo(resumeURL).followRedirects().then(() => {
      this.set('expectedRouteName', this.storage().getItem(RESUME_ROUTE));
      this.set('expectedURL', resumeURL);
      this.initializeMachines();
      this.storage().removeItem(RESUME_URL);
    });
  },

  handleDismissed() {
    this.storage().removeItem(FEATURE_STATE);
    this.storage().removeItem(FEATURE_LIST);
    this.storage().removeItem(COMPONENT_STATE);
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
    this.startFeature();
    if (this.storageHasKey(FEATURE_STATE)) {
      this.saveState('featureState', this.getExtState(FEATURE_STATE));
    }
    this.saveExtState(FEATURE_STATE, this.get('featureState'));
    let nextFeature =
      this.get('featureList').length > 1 ? this.get('featureList').objectAt(1).capitalize() : 'Finish';
    this.set('nextFeature', nextFeature);
    let next;
    if (this.get('currentMachine') === 'secrets' && this.get('featureState') === 'display') {
      next = FeatureMachine.transition(this.get('featureState'), 'REPEAT', this.get('componentState'));
    } else {
      next = FeatureMachine.transition(this.get('featureState'), 'CONTINUE', this.get('componentState'));
    }
    this.saveState('nextStep', next.value);
    let stateNodes = FeatureMachine.getStateNodes(this.get('featureState'));
    this.executeActions(stateNodes.reduce((acc, node) => acc.concat(node.onEntry), []), null, 'feature');
  },

  startFeature() {
    const FeatureMachineConfig = MACHINES[this.get('featureList').objectAt(0)];
    FeatureMachine = Machine(FeatureMachineConfig);
    this.set('currentMachine', this.get('featureList').objectAt(0));
    this.saveState('featureState', FeatureMachine.initialState);
  },

  completeFeature() {
    let features = this.get('featureList');
    let done = features.shift();
    if (!this.getExtState(COMPLETED_FEATURES)) {
      let completed = [];
      completed.push(done);
      this.saveExtState(COMPLETED_FEATURES, completed);
    } else {
      this.saveExtState(COMPLETED_FEATURES, this.getExtState(COMPLETED_FEATURES).toArray().addObject(done));
    }

    this.saveExtState(FEATURE_LIST, features.length ? features : null);
    this.storage().removeItem(FEATURE_STATE);
    if (features.length > 0) {
      this.buildFeatureMachine();
    } else {
      this.storage().removeItem(FEATURE_LIST);
      FeatureMachine = null;
      this.transitionTutorialMachine(this.get('currentState'), 'DONE');
    }
  },

  storage() {
    return getStorage();
  },
});

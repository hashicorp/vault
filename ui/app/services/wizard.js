import Ember from 'ember';
import { Machine } from 'xstate';

const { Service } = Ember;

const CubbyholeMachine = Machine({
  key: 'cubbyhole',
  initial: 'idle',
  states: {
    idle: {
      on: { INTERACTION: 'active' },
    },
    active: {
      on: { DISMISS: 'idle' },
      key: 'feature',
      initial: 'create',
      states: {
        create: {
          on: { EDIT: 'edit' },
          onEntry: [{ type: 'render', component: 'cubbyHoleHelp' }],
        },
        edit: {
          on: { SAVE: 'details' },
          onEntry: [{ type: 'render', component: 'cubbyHoleEditHelp' }],
        },
        details: {
          onEntry: [{ type: 'render', component: 'cubbyHoleSuccess' }],
        },
      },
    },
  },
});

export default Service.extend({
  currentState: null,

  init() {
    this._super(...arguments);
    this.set('currentState', CubbyholeMachine.initialState);
  },

  transitionMachine(intendedState, extendedState) {
    let { actions, value } = CubbyholeMachine.transition(state);
    this.set('currentState', value);
  },
});

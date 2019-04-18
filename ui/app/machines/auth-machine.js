export default {
  key: 'auth',
  initial: 'idle',
  on: {
    RESET: 'idle',
    DONE: 'complete',
  },
  states: {
    idle: {
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.settings.auth.enable'] },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/auth-idle' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    enable: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/auth-enable' },
      ],
      on: {
        CONTINUE: 'config',
      },
    },
    config: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/auth-config' },
      ],
      on: {
        CONTINUE: 'details',
      },
    },
    details: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/auth-details' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'complete',
      },
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};

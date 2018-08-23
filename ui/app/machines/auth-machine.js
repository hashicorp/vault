export default {
  key: 'auth',
  initial: 'enable',
  states: {
    enable: {
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.access'] },
        { type: 'render', level: 'feature', component: 'wizard/auth-enable' },
      ],
      on: {
        CONTINUE: {
          appRole: {
            cond: type => type === 'appRole',
          },
        },
      },
    },
    appRole: {
      key: 'appRole',
      initial: 'details',
      states: {
        details: {
          on: { CONTINUE: 'complete' },
        },
      },
    },
    complete: {
      onEntry: ['completeFeature'],
      on: { RESET: 'idle' },
    },
    paused: {},
  },
};

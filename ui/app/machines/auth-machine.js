export default {
  key: 'auth',
  initial: 'enable',
  states: {
    enable: {
      on: {
        CONTINUE: {
          appRole: {
            cond: (extState, event) => event.selected === 'appRole',
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

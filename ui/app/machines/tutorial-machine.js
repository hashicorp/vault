export default {
  key: 'tutorial',
  initial: 'idle',
  states: {
    active: {
      on: {
        DISMISS: 'dismissed',
        AUTH: 'select',
      },
      key: 'feature',
      initial: 'select',
      states: {
        select: {
          on: {
            CONTINUE: 'feature',
          },
        },
        feature: {},
        init: {
          key: 'init',
          initial: 'setup',
          on: { DONE: 'select' },
          states: {
            setup: {
              on: { CONTINUE: 'save' },
            },
            save: {
              on: { CONTINUE: 'unseal' },
            },
            unseal: {},
          },
        },
      },
    },
    idle: {
      on: {
        DISMISS: 'dismissed',
        CONTINUE: 'active',
      },
    },
    dismissed: {
      on: { CONTINUE: 'idle' },
      onEntry: ['handleDismissed'],
    },
    paused: {
      on: { CONTINUE: ['handlePause'] },
    },
    complete: {
      on: {
        CONTINUE: 'idle',
        DISMISS: 'dismissed',
      },
    },
  },
};

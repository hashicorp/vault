export default {
  key: 'tutorial',
  initial: 'idle',
  states: {
    active: {
      on: { DISMISS: 'dismissed' },
      key: 'feature',
      initial: 'select',
      states: {
        select: {
          on: {
            CONTINUE: {
              onExit: ['saveFeatures'],
            },
          },
        },
      },
    },
    idle: {
      on: {
        DISMISS: 'dismissed',
        INTERACTION: 'active',
      },
    },
    dismissed: {
      on: { RESET: 'idle' },
      onEntry: ['saveState'],
    },
    complete: {
      on: {
        RESET: 'idle',
        DISMISS: 'dismissed',
      },
    },
  },
};

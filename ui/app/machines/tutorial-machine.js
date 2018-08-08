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
            CONTINUE: 'feature',
          },
        },
        feature: {},
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
      onEntry: ['saveState'],
    },
    complete: {
      on: {
        CONTINUE: 'idle',
        DISMISS: 'dismissed',
      },
    },
  },
};

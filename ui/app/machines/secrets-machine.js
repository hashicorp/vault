import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
const supportedBackends = supportedSecretBackends();

export default {
  key: 'secrets',
  initial: 'idle',
  on: {
    RESET: 'idle',
    DONE: 'complete',
    ERROR: 'error',
  },
  states: {
    idle: {
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.settings.mount-secret-backend'] },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/secrets-idle' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    enable: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/secrets-enable' },
      ],
      on: {
        CONTINUE: {
          details: { cond: (type) => supportedBackends.includes(type) },
          list: { cond: (type) => !supportedBackends.includes(type) },
        },
      },
    },
    details: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/secrets-details' },
      ],
      on: {
        CONTINUE: {
          connection: {
            cond: (type) => type === 'database',
          },
          role: {
            cond: (type) => ['pki', 'aws', 'ssh'].includes(type),
          },
          secret: {
            cond: (type) => ['kv'].includes(type),
          },
          encryption: {
            cond: (type) => type === 'transit',
          },
          provider: {
            cond: (type) => type === 'keymgmt',
          },
        },
      },
    },
    connection: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-connection' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'displayConnection',
      },
    },
    encryption: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/secrets-encryption' },
      ],
      on: {
        CONTINUE: 'display',
      },
    },
    credentials: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-credentials' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'display',
      },
    },
    role: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-role' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'displayRole',
      },
    },
    displayRole: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-display-role' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'credentials',
      },
    },
    displayConnection: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-connection-show' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'displayRoleDatabase',
      },
    },
    displayRoleDatabase: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-display-database-role' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'display',
      },
    },
    secret: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-secret' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'display',
      },
    },
    provider: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-keymgmt' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'displayProvider',
      },
    },
    displayProvider: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-keymgmt' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'distribute',
      },
    },
    distribute: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-keymgmt' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'display',
      },
    },
    display: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-display' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        REPEAT: {
          connection: {
            cond: (type) => type === 'database',
            actions: [{ type: 'routeTransition', params: ['vault.cluster.secrets.backend.create-root'] }],
          },
          role: {
            cond: (type) => ['pki', 'aws', 'ssh'].includes(type),
            actions: [{ type: 'routeTransition', params: ['vault.cluster.secrets.backend.create-root'] }],
          },
          secret: {
            cond: (type) => ['kv'].includes(type),
            actions: [{ type: 'routeTransition', params: ['vault.cluster.secrets.backend.create-root'] }],
          },
          encryption: {
            cond: (type) => type === 'transit',
            actions: [{ type: 'routeTransition', params: ['vault.cluster.secrets.backend.create-root'] }],
          },
          provider: {
            cond: (type) => type === 'keymgmt',
            actions: [
              {
                type: 'routeTransition',
                params: [
                  'vault.cluster.secrets.backend.create-root',
                  { queryParams: { itemType: 'provider' } },
                ],
              },
            ],
          },
        },
      },
    },
    list: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-list' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'display',
      },
    },
    error: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/tutorial-error' },
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

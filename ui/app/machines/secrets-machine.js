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
        CONTINUE: {
          ad: {
            cond: type => type === 'ad',
          },
          aws: {
            cond: type => type === 'aws',
          },
          consul: {
            cond: type => type === 'aws',
          },
          cubbyhole: {
            cond: type => type === 'cubbyhole',
          },
          database: {
            cond: type => type === 'database',
          },
          gcp: {
            cond: type => type === 'gcp',
          },
          kv: {
            cond: type => type === 'kv',
          },
          nomad: {
            cond: type => type === 'nomad',
          },
          pki: {
            cond: type => type === 'pki',
          },
          rabbitmq: {
            cond: type => type === 'rabbitmq',
          },
          ssh: {
            cond: type => type === 'ssh',
          },
          totp: {
            cond: type => type === 'totp',
          },
          transit: {
            cond: type => type === 'transit',
          },
        },
      },
    },
    enable: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/secrets-enable' },
      ],
      on: {
        CONTINUE: {
          details: { cond: type => supportedBackends.includes(type) },
          display: {
            cond: type => ['ad', 'consul', 'database', 'gcp', 'nomad', 'rabbitmq', 'totp'].includes(type),
          },
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
          role: {
            cond: type => ['pki', 'aws', 'ssh'].includes(type),
          },
          secret: {
            cond: type =>
              ['cubbyhole', 'database', 'gcp', 'kv', 'nomad', 'rabbitmq', 'totp', 'transit'].includes(type),
          },
        },
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
    secret: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/secrets-secret' },
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
      REPEAT: {
        role: {
          cond: type => ['pki', 'aws', 'ssh'].includes(type),
        },
        secret: {
          cond: type =>
            ['cubbyhole', 'database', 'gcp', 'kv', 'nomad', 'rabbitmq', 'totp', 'transit'].includes(type),
        },
      },
    },
    ad: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/ad-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    aws: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/aws-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    consul: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/consul-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    cubbyhole: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/ch-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    database: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/database-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    gcp: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/gcp-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    kv: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/kv-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    nomad: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/nomad-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    pki: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/pki-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    rabbitmq: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/rabbitmq-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    ssh: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/ssh-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    totp: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/totp-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    transit: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/transit-engine' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    error: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/tutorial-error' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'idle',
      },
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};

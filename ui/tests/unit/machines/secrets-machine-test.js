/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { Machine } from 'xstate';
import SecretsMachineConfig from 'vault/machines/secrets-machine';

module('Unit | Machine | secrets-machine', function () {
  const secretsMachine = Machine(SecretsMachineConfig);

  const testCases = [
    {
      currentState: secretsMachine.initialState,
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'enable',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-enable', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'aws',
      expectedResults: {
        value: 'details',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-details', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      params: 'aws',
      expectedResults: {
        value: 'role',
        actions: [
          { component: 'wizard/secrets-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'role',
      event: 'CONTINUE',
      params: 'aws',
      expectedResults: {
        value: 'displayRole',
        actions: [
          { component: 'wizard/secrets-display-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'displayRole',
      event: 'CONTINUE',
      params: 'aws',
      expectedResults: {
        value: 'credentials',
        actions: [
          { component: 'wizard/secrets-credentials', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'credentials',
      event: 'CONTINUE',
      params: 'aws',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'REPEAT',
      params: 'aws',
      expectedResults: {
        value: 'role',
        actions: [
          {
            params: ['vault.cluster.secrets.backend.create-root'],
            type: 'routeTransition',
          },
          { component: 'wizard/secrets-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'aws',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'aws',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'aws',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'pki',
      expectedResults: {
        value: 'details',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-details', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      params: 'pki',
      expectedResults: {
        value: 'role',
        actions: [
          { component: 'wizard/secrets-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'role',
      event: 'CONTINUE',
      params: 'pki',
      expectedResults: {
        value: 'displayRole',
        actions: [
          { component: 'wizard/secrets-display-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'displayRole',
      event: 'CONTINUE',
      params: 'pki',
      expectedResults: {
        value: 'credentials',
        actions: [
          { component: 'wizard/secrets-credentials', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'credentials',
      event: 'CONTINUE',
      params: 'pki',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'REPEAT',
      params: 'pki',
      expectedResults: {
        value: 'role',
        actions: [
          {
            params: ['vault.cluster.secrets.backend.create-root'],
            type: 'routeTransition',
          },
          { component: 'wizard/secrets-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'pki',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'pki',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'pki',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'ssh',
      expectedResults: {
        value: 'details',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-details', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      params: 'ssh',
      expectedResults: {
        value: 'role',
        actions: [
          { component: 'wizard/secrets-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'role',
      event: 'CONTINUE',
      params: 'ssh',
      expectedResults: {
        value: 'displayRole',
        actions: [
          { component: 'wizard/secrets-display-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'displayRole',
      event: 'CONTINUE',
      params: 'ssh',
      expectedResults: {
        value: 'credentials',
        actions: [
          { component: 'wizard/secrets-credentials', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'credentials',
      event: 'CONTINUE',
      params: 'ssh',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'REPEAT',
      params: 'ssh',
      expectedResults: {
        value: 'role',
        actions: [
          {
            params: ['vault.cluster.secrets.backend.create-root'],
            type: 'routeTransition',
          },
          { component: 'wizard/secrets-role', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'ssh',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'ssh',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'ssh',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'consul',
      expectedResults: {
        value: 'list',
        actions: [
          { type: 'render', level: 'step', component: 'wizard/secrets-list' },
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        ],
      },
    },
    {
      currentState: 'list',
      event: 'CONTINUE',
      params: 'consul',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'consul',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'consul',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'consul',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'database',
      expectedResults: {
        value: 'details',
        actions: [
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
          { type: 'render', level: 'step', component: 'wizard/secrets-details' },
        ],
      },
    },
    {
      currentState: 'list',
      event: 'CONTINUE',
      params: 'database',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'database',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'database',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'database',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'gcp',
      expectedResults: {
        value: 'list',
        actions: [
          { type: 'render', level: 'step', component: 'wizard/secrets-list' },
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        ],
      },
    },
    {
      currentState: 'list',
      event: 'CONTINUE',
      params: 'gcp',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'gcp',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'gcp',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'gcp',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'nomad',
      expectedResults: {
        value: 'list',
        actions: [
          { type: 'render', level: 'step', component: 'wizard/secrets-list' },
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        ],
      },
    },
    {
      currentState: 'list',
      event: 'CONTINUE',
      params: 'nomad',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'nomad',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'nomad',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'nomad',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'rabbitmq',
      expectedResults: {
        value: 'list',
        actions: [
          { type: 'render', level: 'step', component: 'wizard/secrets-list' },
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        ],
      },
    },
    {
      currentState: 'list',
      event: 'CONTINUE',
      params: 'rabbitmq',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'rabbitmq',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'rabbitmq',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'rabbitmq',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'totp',
      expectedResults: {
        value: 'list',
        actions: [
          { type: 'render', level: 'step', component: 'wizard/secrets-list' },
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        ],
      },
    },
    {
      currentState: 'list',
      event: 'CONTINUE',
      params: 'totp',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'totp',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'totp',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'totp',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'kv',
      expectedResults: {
        value: 'details',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-details', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      params: 'kv',
      expectedResults: {
        value: 'secret',
        actions: [
          { component: 'wizard/secrets-secret', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'secret',
      event: 'CONTINUE',
      params: 'kv',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'REPEAT',
      params: 'kv',
      expectedResults: {
        value: 'secret',
        actions: [
          {
            params: ['vault.cluster.secrets.backend.create-root'],
            type: 'routeTransition',
          },
          { component: 'wizard/secrets-secret', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'kv',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'kv',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'kv',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: 'transit',
      expectedResults: {
        value: 'details',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-details', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      params: 'transit',
      expectedResults: {
        value: 'encryption',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-encryption', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'encryption',
      event: 'CONTINUE',
      params: 'transit',
      expectedResults: {
        value: 'display',
        actions: [
          { component: 'wizard/secrets-display', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'REPEAT',
      params: 'transit',
      expectedResults: {
        value: 'encryption',
        actions: [
          {
            params: ['vault.cluster.secrets.backend.create-root'],
            type: 'routeTransition',
          },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/secrets-encryption', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'RESET',
      params: 'transit',
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.mount-secret-backend'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/secrets-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
    {
      currentState: 'display',
      event: 'DONE',
      params: 'transit',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'display',
      event: 'ERROR',
      params: 'transit',
      expectedResults: {
        value: 'error',
        actions: [
          { component: 'wizard/tutorial-error', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
  ];

  testCases.forEach((testCase) => {
    test(`transition: ${testCase.event} for currentState ${testCase.currentState} and componentState ${testCase.params}`, function (assert) {
      const result = secretsMachine.transition(testCase.currentState, testCase.event, testCase.params);
      assert.strictEqual(result.value, testCase.expectedResults.value);
      assert.deepEqual(result.actions, testCase.expectedResults.actions);
    });
  });
});

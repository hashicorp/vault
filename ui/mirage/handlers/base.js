/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// base handlers used in mirage config when a specific handler is not specified
const EXPIRY_DATE = '2021-05-12T23:20:50.52Z';

export default function (server) {
  server.get('/sys/internal/ui/feature-flags', (db) => {
    const featuresResponse = db.features.first();
    return {
      data: {
        feature_flags: featuresResponse ? featuresResponse.feature_flags : null,
      },
    };
  });

  server.get('/sys/health', function (schema) {
    return schema.db.healths.firstOrCreate({}, server.create('health'));
  });

  server.get('/sys/activation-flags', () => {
    return {
      data: {
        activated: [''],
        unactivated: ['secrets-sync'],
      },
    };
  });

  server.get('/sys/license/features', function () {
    return {
      features: [
        [
          'HSM',
          'Performance Replication',
          'DR Replication',
          'MFA',
          'Sentinel',
          'Seal Wrapping',
          'Control Groups',
          'Performance Standby',
          'Namespaces',
          'KMIP',
          'Entropy Augmentation',
          'Transform Secrets Engine',
          'Lease Count Quotas',
          'Key Management Secrets Engine',
          'Automated Snapshots',
          'Key Management Transparent Data Encryption',
          'Secrets Sync',
          'Secrets Import',
          'Oracle Database Secrets Engine',
        ],
      ],
    };
  });

  server.get('/sys/license/status', function () {
    return {
      data: {
        autoloading_used: false,
        persisted_autoload: {
          expiration_time: EXPIRY_DATE,
          features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
          license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
          performance_standby_count: 0,
          start_time: '2020-04-28T00:00:00Z',
        },
        autoloaded: {
          expiration_time: EXPIRY_DATE,
          features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
          license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
          performance_standby_count: 0,
          start_time: '2020-04-28T00:00:00Z',
        },
      },
    };
  });

  server.get('/sys/seal-status', function (schema) {
    return schema.db.sealStatuses.firstOrCreate({}, server.create('seal-status'));
  });

  server.get('/sys/replication/status', function (schema) {
    return schema.db.replicationStatuses.firstOrCreate({}, server.create('replication-status'));
  });

  server.get('sys/namespaces', function () {
    return {
      data: {
        keys: [
          'ns1/',
          'ns2/',
          'ns3/',
          'ns4/',
          'ns5/',
          'ns6/',
          'ns7/',
          'ns8/',
          'ns9/',
          'ns10/',
          'ns11/',
          'ns12/',
          'ns13/',
          'ns14/',
          'ns15/',
          'ns16/',
          'ns17/',
          'ns18/',
        ],
      },
    };
  });

  server.get('/sys/internal/ui/unauthenticated-messages', function () {
    return { data: {} };
  });

  server.get('/sys/internal/ui/authenticated-messages', function () {
    return { data: {} };
  });

  server.get('/sys/internal/ui/default-auth-methods', function () {
    return { data: null };
  });

  // defaults to root token
  server.get('/auth/token/lookup-self', function (schema) {
    return schema.db.tokens.firstOrCreate({}, server.create('token'));
  });

  server.post('/sys/capabilities-self', function (schema, req) {
    const { paths } = JSON.parse(req.requestBody);
    const response = {
      data: {
        capabilities: ['root'],
      },
    };

    paths.forEach((path) => {
      response.data[path] = ['root'];
    });

    return response;
  });

  server.get('/sys/internal/ui/resultant-acl', function () {
    return {
      data: { chroot_namespace: '', root: true },
    };
  });

  server.get('/sys/internal/ui/mounts', function () {
    return { auth: {}, data: {} };
  });
}

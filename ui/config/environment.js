/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
'use strict';

module.exports = function (environment) {
  const ENV = {
    modulePrefix: 'vault',
    environment,
    rootURL: '/ui/',
    serviceWorkerScope: '/v1/sys/storage/raft/snapshot',
    locationType: 'history',
    EmberENV: {
      FEATURES: {
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      },
      EXTEND_PROTOTYPES: {
        Date: false, // Prevent Ember Data from overriding Date.parse.
        String: false, // Prevent user from using an Ember string method on string. ex: "foo".capitalize();
      },
    },

    APP: {
      // endpoints that the UI polls
      POLLING_URLS: ['sys/health', 'sys/replication/status', 'sys/seal-status'],
      // endpoints that UI uses to determine the cluster state
      // calls to these endpoints will always go to the root namespace
      // these also need to be updated in the open-api-explorer engine
      NAMESPACE_ROOT_URLS: [
        'sys/health',
        'sys/seal-status',
        'sys/license/features',
        'sys/internal/counters/config',
      ],
      // number of records to show on a single page by default - this is used by the client-side pagination
      DEFAULT_PAGE_SIZE: 100,

      ANALYTICS_CONFIG: { enabled: false },
    },
    flashMessageDefaults: {
      timeout: 7000,
      sticky: false,
    },
  };

  if (environment === 'development') {
    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
    if (process.env.MIRAGE_DEV_HANDLER !== undefined) {
      ENV['ember-cli-mirage'] = {
        enabled: true,
        handler: process.env.MIRAGE_DEV_HANDLER,
      };
    }

    if (process.env.ENABLE_POSTHOG) {
      ENV.APP.ANALYTICS_CONFIG = {
        provider: 'posthog',
        enabled: true,
        project_id: 'phc_zPQ9fPlFj4ZTYKJmThG1C8AE4J4RgPQx8dJJ7agg4SG',
        api_host: 'https://eu.i.posthog.com',
      };
    }
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.locationType = 'none';
    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;
    ENV.APP.rootElement = '#ember-testing';
    ENV.APP.autoboot = false;
    ENV.flashMessageDefaults.timeout = 50;
    ENV['ember-cli-mirage'] = {
      enabled: false,
    };
    ENV.APP.ANALYTICS_CONFIG = { enabled: false };
  }

  if (environment !== 'production') {
    ENV.APP.DEFAULT_PAGE_SIZE = 15;
  }

  if (environment === 'production') {
    ENV.APP.ANALYTICS_CONFIG = {
      provider: 'posthog',
      enabled: true,
      project_id: 'phc_pIw6t5numW5jDram4dnJjSnwDOorf9IGd1MmlFp0dHh',
      api_host: 'https://eu.i.posthog.com',
    };
  }

  ENV.welcomeMessage = process.env.UI_AUTH_WELCOME;

  return ENV;
};

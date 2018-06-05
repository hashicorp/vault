/* jshint node: true */

module.exports = function(environment) {
  var ENV = {
    modulePrefix: 'vault',
    environment: environment,
    rootURL: '/ui/',
    locationType: 'auto',
    EmberENV: {
      FEATURES: {
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      },
      EXTEND_PROTOTYPES: {
        // Prevent Ember Data from overriding Date.parse.
        Date: false,
      },
    },

    APP: {
      // Here you can pass flags/options to your application instance
      // when it is created
    },
    flashMessageDefaults: {
      timeout: 7000,
      sticky: false,
      preventDuplicates: true,
    },
  };
  if (environment === 'development') {
    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
    //ENV['ember-cli-mirage'] = {
    //enabled: true
    //};
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.locationType = 'none';

    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;

    ENV.APP.rootElement = '#ember-testing';

    ENV['ember-cli-mirage'] = {
      enabled: false,
    };
  }
  if (environment !== 'production') {
    ENV.contentSecurityPolicyHeader = 'Content-Security-Policy';
    ENV.contentSecurityPolicyMeta = true;
    ENV.contentSecurityPolicy = {
      'connect-src': ["'self'"],
      'img-src': ["'self'", 'data:'],
      'form-action': ["'none'"],
      'script-src': ["'self'"],
      'style-src': ["'unsafe-inline'", "'self'"],
    };
  }

  if (environment === 'production') {
  }

  return ENV;
};

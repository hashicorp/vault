module.exports = function (environment) {
  return {
    delivery: ['header', 'meta'],
    enabled: environment !== 'production',
    failTests: true,
    policy: {
      'default-src': ["'none'"],
      'script-src': ["'self'"],
      'font-src': ["'self'"],
      'connect-src': ["'self'", 'ws://127.0.0.1:9201'],
      'img-src': ["'self'", 'data:'],
      'style-src': ["'unsafe-inline'", "'self'"],
      'media-src': ["'self'"],
      'form-action': ["'none'"],
    },
    reportOnly: false,
  };
};

/* eslint-env node */
'use strict';

module.exports = function (environment) {
  const ENV = {
    modulePrefix: 'kubernetes',
    environment,
  };

  return ENV;
};

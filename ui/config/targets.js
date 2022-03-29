'use strict';

const browsers = ['last 1 Chrome versions', 'last 1 Firefox versions', 'last 1 Safari versions'];

// Ember's browser support policy is changing, and IE11 support will end in
// v4.0 onwards.
//
// See https://deprecations.emberjs.com/v3.x#toc_3-0-browser-support-policy
//
// If you need IE11 support on a version of Ember that still offers support
// for it, uncomment the code block below.
//
// const isCI = Boolean(process.env.CI);
// const isProduction = process.env.EMBER_ENV === 'production';
//
// if (isCI || isProduction) {
//   browsers.push('ie 11');
// }

module.exports = {
  browsers,
};

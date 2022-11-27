/* eslint-env node */

'use strict';

const vault_addr = process.env.VAULT_ADDR;
console.log('VAULT_ADDR=' + vault_addr); // eslint-disable-line

module.exports = {
  test_page: 'tests/index.html?hidepassed',
  tap_quiet_logs: true,
  tap_failed_tests_only: true,
  disable_watching: true,
  launch_in_ci: ['Chrome'],
  browser_start_timeout: 120,
  browser_args: {
    Chrome: {
      ci: [
        // --no-sandbox is needed when running Chrome inside a container
        process.env.CI ? '--no-sandbox' : null,
        '--headless',
        '--disable-dev-shm-usage',
        '--disable-software-rasterizer',
        '--mute-audio',
        '--remote-debugging-port=0',
        '--window-size=1440,900',
      ].filter(Boolean),
    },
  },
  proxies: {
    '/v1': {
      target: vault_addr,
    },
  },
};

if (process.env.CI) {
  module.exports.reporter = 'xunit';
  module.exports.report_file = 'test-results/qunit/results.xml';
  module.exports.xunit_intermediate_output = true;
}

/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

'use strict';

module.exports = {
  test_page: 'tests/index.html?hidepassed&enableA11yAudit',
  tap_quiet_logs: true,
  tap_failed_tests_only: true,
  disable_watching: true,
  launch_in_ci: ['Chrome'],
  browser_start_timeout: 120,
  socket_server_options: {
    maxHttpBufferSize: 1e8, // set socket.io server max buffer size to 100MB
  },
  browser_args: {
    Chrome: {
      ci: [
        // --no-sandbox is needed when running Chrome inside a container
        process.env.CI ? '--no-sandbox' : null,
        '--headless=new',
        '--disable-dev-shm-usage',
        '--disable-software-rasterizer',
        '--disable-search-engine-choice-screen', // needed from Chrome 127
        '--mute-audio',
        '--remote-debugging-port=0',
        '--window-size=1440,900',
      ].filter(Boolean),
    },
    Firefox: {
      ci: ['--headless'],
    },
  },
  proxies: {
    '/v1': {
      target: 'http://127.0.0.1:9200',
    },
  },
  parallel: process.env.EMBER_EXAM_SPLIT_COUNT || 1,
};

if (process.env.CI) {
  module.exports.reporter = 'xunit';
  module.exports.report_file = 'test-results/qunit/results.xml';
  module.exports.xunit_intermediate_output = true;
}

const config = {
  framework: 'qunit',
  test_page: 'tests/index.html?hidepassed',
  tap_quiet_logs: true,
  disable_watching: true,
  launch_in_ci: ['Chrome'],
  browser_args: {
    Chrome: {
      ci: [
        // --no-sandbox is needed when running Chrome inside a container
        process.env.CI ? '--no-sandbox' : null,
        '--headless',
        '--disable-gpu',
        '--disable-software-rasterizer',
        '--mute-audio',
        '--remote-debugging-port=0',
        '--window-size=1440,900',
      ].filter(Boolean),
    },
  },
  on_exit:
    '[ -e ../../vault-ui-integration-server.pid ] && node ../../scripts/start-vault.js `cat ../../vault-ui-integration-server.pid`; [ -e ../../vault-ui-integration-server.pid ] && rm ../../vault-ui-integration-server.pid',

  proxies: {
    '/v1': {
      target: 'http://localhost:9200',
    },
  },
};

if (process.env.CI) {
  config.reporter = 'xunit';
  config.report_file = 'test-reports/ember.xml';
  config.xunit_intermediate_output = true;
}

module.exports = config;

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
        // as per https://github.com/ember-cli/ember-cli/pull/8774
        '--disable-software-rasterizer',
        '--mute-audio',
        '--remote-debugging-port=0',
        '--window-size=1440,900',
      ].filter(Boolean),
    },
  },
  proxies: {
    '/v1': {
      target: 'http://localhost:9200',
    },
  },
};

if (process.env.CI) {
  config.reporter = 'xunit';
  config.report_file = 'test-results/qunit/results.xml';
  config.xunit_intermediate_output = true;
}

module.exports = config;

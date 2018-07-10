/*jshint node:true*/

module.exports = {
  framework: 'qunit',
  test_page: 'tests/index.html?hidepassed',
  tap_quiet_logs: true,
  disable_watching: true,
  launch_in_ci: ['Chrome'],
  browser_args: {
    Chrome: {
      mode: 'ci',
      args: [
        // --no-sandbox is needed when running Chrome inside a container
        process.env.TRAVIS ? '--no-sandbox' : null,

        '--disable-gpu',
        '--headless',
        '--remote-debugging-port=0',
        '--window-size=1440,900',
      ].filter(Boolean),
    },
  },
  launch_in_dev: ['Chrome'],
  on_exit:
    '[ -e ../../vault-ui-integration-server.pid ] && node ../../scripts/start-vault.js `cat ../../vault-ui-integration-server.pid`; [ -e ../../vault-ui-integration-server.pid ] && rm ../../vault-ui-integration-server.pid',

  proxies: {
    '/v1': {
      target: 'http://localhost:9200',
    },
  },
};

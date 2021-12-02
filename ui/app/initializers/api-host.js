import config from '../config/environment';

export function initialize(application) {
  if (config.environment === 'electron') {
    const { ipcRenderer } = window.requireNode('electron');
    ipcRenderer.on('api-host', (event, host) => {
      config.apiHost = host;
      application.advanceReadiness();
    });
    ipcRenderer.send('get-api-host');
    application.deferReadiness();
  }
}

export default {
  before: ['ember-inspect-disable', 'enable-engines'],
  initialize,
};

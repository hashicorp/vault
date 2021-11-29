const {
  default: installExtension,
  EMBER_INSPECTOR,
} = require('electron-devtools-installer');
const path = require('path');
const { app } = require('electron');
const handleFileUrls = require('../src/handle-file-urls');
const {
  setupTestem,
  openTestWindow,
} = require('ember-electron/lib/test-support');

const emberAppDir = path.resolve(__dirname, '..', 'ember-test');

app.on('ready', async function onReady() {
  try {
    require('devtron').install();
  } catch (err) {
    console.log('Failed to install Devtron: ', err);
  }
  try {
    await installExtension(EMBER_INSPECTOR, {
      loadExtensionOptions: { allowFileAccess: true },
    });
  } catch (err) {
    console.log('Failed to install Ember Inspector: ', err);
  }

  await handleFileUrls(emberAppDir);
  setupTestem();
  openTestWindow(emberAppDir);
});

app.on('window-all-closed', function onWindowAllClosed() {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

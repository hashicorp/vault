#!/usr/bin/env node
/* eslint-env node */
/* eslint-disable no-console */
const {
  default: installExtension,
  EMBER_INSPECTOR,
} = require('electron-devtools-installer');
const { pathToFileURL } = require('url');
const { app, BrowserWindow, screen } = require('electron');
const path = require('path'); 
const handleFileUrls = require('./handle-file-urls');
const isDev = require('electron-is-dev');
const windowStateKeeper = require('electron-window-state');

const emberAppDir = path.resolve(__dirname, '..', 'ember-dist');
const emberAppURL = pathToFileURL(
  path.join(emberAppDir, 'index.html')
).toString();

// Uncomment the lines below to enable Electron's crash reporter
// For more information, see http://electron.atom.io/docs/api/crash-reporter/
// electron.crashReporter.start({
//     productName: 'YourName',
//     companyName: 'YourCompany',
//     submitURL: 'https://your-domain.com/url-to-submit',
//     autoSubmit: true
// });

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('ready', async () => {
  if (isDev) {
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
  }

  await handleFileUrls(emberAppDir);

  const primaryDisplay = screen.getPrimaryDisplay();
  const mainWindowState = windowStateKeeper({
    defaultWidth: primaryDisplay.size.width / 2,
    defaultHeight: primaryDisplay.size.height
  }); 
  let mainWindow = new BrowserWindow({
    x: mainWindowState.x,
    y: mainWindowState.y,
    width: primaryDisplay.size.width / 2,
    height: primaryDisplay.size.height,
    useContentSize: true
  });

  // register listeners on the window to update the state
  // listeners will be removed automatically when the window is closed and restore the maximized or full screen state
  mainWindowState.manage(mainWindow);

  // If you want to open up dev tools programmatically, call
  // mainWindow.openDevTools();

  // Load the ember application
  mainWindow.loadURL(emberAppURL);

  // If a loading operation goes wrong, we'll send Electron back to
  // Ember App entry point
  mainWindow.webContents.on('did-fail-load', () => {
    mainWindow.loadURL(emberAppURL);
  });

  mainWindow.webContents.on('render-process-gone', (_event, details) => {
    if (details.reason === 'killed' || details.reason === 'clean-exit') {
      return;
    }
    console.log(
      'Your main window process has exited unexpectedly -- see https://www.electronjs.org/docs/api/web-contents#event-render-process-gone'
    );
    console.log('Reason: ' + details.reason);
  });

  mainWindow.on('unresponsive', () => {
    console.log(
      'Your Ember app (or other code) has made the window unresponsive.'
    );
  });

  mainWindow.on('responsive', () => {
    console.log('The main window has become responsive again.');
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
});

// Handle an unhandled error in the main thread
//
// Note that 'uncaughtException' is a crude mechanism for exception handling intended to
// be used only as a last resort. The event should not be used as an equivalent to
// "On Error Resume Next". Unhandled exceptions inherently mean that an application is in
// an undefined state. Attempting to resume application code without properly recovering
// from the exception can cause additional unforeseen and unpredictable issues.
//
// Attempting to resume normally after an uncaught exception can be similar to pulling out
// of the power cord when upgrading a computer -- nine out of ten times nothing happens -
// but the 10th time, the system becomes corrupted.
//
// The correct use of 'uncaughtException' is to perform synchronous cleanup of allocated
// resources (e.g. file descriptors, handles, etc) before shutting down the process. It is
// not safe to resume normal operation after 'uncaughtException'.
process.on('uncaughtException', (err) => {
  console.log('An exception in the main thread was not handled.');
  console.log(
    'This is a serious issue that needs to be handled and/or debugged.'
  );
  console.log(`Exception: ${err}`);
});

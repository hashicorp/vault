#!/usr/bin/env node
/* eslint-env node */
/* eslint-disable no-console */
const {
  default: installExtension,
  EMBER_INSPECTOR,
} = require('electron-devtools-installer');
const { pathToFileURL } = require('url');
const { app, BrowserWindow, screen, globalShortcut, ipcMain } = require('electron');
const path = require('path'); 
const handleFileUrls = require('./handle-file-urls');
const isDev = require('electron-is-dev');
const windowStateKeeper = require('electron-window-state');
const setApplicationMenu = require('./app-menu');

// Uncomment the lines below to enable Electron's crash reporter
// For more information, see http://electron.atom.io/docs/api/crash-reporter/
// electron.crashReporter.start({
//     productName: 'YourName',
//     companyName: 'YourCompany',
//     submitURL: 'https://your-domain.com/url-to-submit',
//     autoSubmit: true
// });

const getWindowState = () => {
  const primaryDisplay = screen.getPrimaryDisplay();
  const windowState = windowStateKeeper({
    defaultWidth: primaryDisplay.size.width / 2,
    defaultHeight: primaryDisplay.size.height
  });
  const { x, y, width, height } = windowState;
  return [
    windowState,
    { x, y, width, height }
  ];
};

const setupBrowserWindow = (options, contentPath) => {
  const [windowState, windowStateProps] = getWindowState();
  let appWindow = new BrowserWindow({
    ...windowStateProps,
    ...options
  });
  // register listeners on the window to update the state
  // listeners will be removed automatically when the window is closed and restore the maximized or full screen state
  windowState.manage(appWindow);
  // Load the ember application or html content
  const loadMethod = contentPath.includes('file://') ? 'loadURL' : 'loadFile';
  appWindow[loadMethod](contentPath);
  // If a loading operation goes wrong, we'll send Electron back to
  // Ember App entry point
  appWindow.webContents.on('did-fail-load', () => {
    appWindow.loadURL(contentPath);
  });
  appWindow.webContents.on('render-process-gone', (_event, details) => {
    if (details.reason === 'killed' || details.reason === 'clean-exit') {
      return;
    }
    console.log(
      'Your main window process has exited unexpectedly -- see https://www.electronjs.org/docs/api/web-contents#event-render-process-gone'
    );
    console.log('Reason: ' + details.reason);
  });
  appWindow.on('unresponsive', () => {
    console.log(
      'Your Ember app (or other code) has made the window unresponsive.'
    );
  });
  appWindow.on('responsive', () => {
    console.log('The main window has become responsive again.');
  });
  appWindow.on('closed', () => {
    appWindow = null;
    globalShortcut.unregister('CommandOrControl+R');
  });
  // setup reload shortcut
  if (isDev) {
    globalShortcut.register('CommandOrControl+R', () => {
      appWindow.reload();
    });
  }
  this.appWindow = appWindow;
};

const renderEmberApp = async serverAddress => {
  const emberAppDir = path.resolve(__dirname, '..', 'ember-dist');
  const emberAppURL = pathToFileURL(path.join(emberAppDir, 'index.html')).toString();
  await handleFileUrls(emberAppDir);
  const options = {
    useContentSize: true,
    // there are secure ways to access node in renderers
    // below should be replaced with preload and contextBridge strategy
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
      // webSecurity: false
    }
  };
  setupBrowserWindow(options, emberAppURL);
  // provide server address to ember app
  ipcMain.on('get-api-host', event => {
    event.reply('api-host', serverAddress);
  });
};

const renderServerSelect = () => {
  const options = {
    // there are secure ways to access node in renderers
    // below should be replaced with preload and contextBridge strategy
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  };
  const contentPath = path.join(__dirname, 'server-select.html');
  setupBrowserWindow(options, contentPath);
};

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

  // setup custom application menu
  setApplicationMenu();

  // create setting to skip server select screen
  // for now render by default
  renderServerSelect();

  // communication with renderers
  // listen for server selected event to switch window to ember app
  ipcMain.on('server-selected', (event, serverAddress) => {
    this.appWindow.close();
    renderEmberApp(serverAddress);
  });
  // show server select -- triggered via menu
  ipcMain.on('show-servers', () => {
    this.appWindow.close();
    renderServerSelect();
  });
  ipcMain.on('get-preferences-path', event => {
    const filePath = path.join(app.getPath('userData'), 'preferences.json');
    event.reply('preferences-path', filePath);
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
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

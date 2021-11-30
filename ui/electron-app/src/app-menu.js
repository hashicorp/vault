#!/usr/bin/env node
/* eslint-env node */
const { app, Menu, shell } = require('electron');

const isMac = process.platform === 'darwin';
const template = [
  ...(isMac ? [{
    label: app.name,
    submenu: [
      { role: 'about' },
      { type: 'separator' },
      { role: 'services' },
      { type: 'separator' },
      { role: 'hide' },
      { role: 'hideOthers' },
      { role: 'unhide' },
      { type: 'separator' },
      { role: 'quit' }
    ]
  }] : []),
  {
    label: 'File',
    submenu: [
      isMac ? { role: 'close' } : { role: 'quit' }
    ]
  },
  {
    label: 'View',
    submenu: [
      { role: 'reload' },
      { role: 'forceReload' },
      { role: 'toggleDevTools' },
      { type: 'separator' },
      { role: 'resetZoom' },
      { role: 'zoomIn' },
      { role: 'zoomOut' },
      { type: 'separator' },
      { role: 'togglefullscreen' }
    ]
  },
  {
    label: 'Window',
    submenu: [
      { role: 'minimize' },
      { role: 'zoom' },
      ...(isMac ? [
        { type: 'separator' },
        { role: 'front' },
        { type: 'separator' },
        { role: 'window' }
      ] : [
        { role: 'close' }
      ])
    ]
  },
  {
    role: 'help',
    submenu: [
      {
        label: 'Reference Guide',
        click: async () => {
          await shell.openExternal('https://www.vaultproject.io/docs');
        }
      },
      {
        label: 'Tutorials',
        click: async () => {
          await shell.openExternal('https://learn.hashicorp.com/vault');
        }
      },
      {
        label: 'API Documentation',
        click: async () => {
          await shell.openExternal('https://www.vaultproject.io/api-docs');
        }
      }
    ]
  }
];
module.exports = () => Menu.setApplicationMenu(Menu.buildFromTemplate(template));

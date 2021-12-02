#!/usr/bin/env node
/* eslint-env node */
const { ipcRenderer } = require('electron');
const fs = require('fs');

const writePreferences = () => {
  try {
    fs.writeFileSync(this.filePath, JSON.stringify(this.preferences));
  } catch (error) {
    console.log(error);
    // silently fail persisting value for now
  }
};
const createServerElements = () => {
  const markup = this.preferences.servers.map(({ address }, index) => `
    <div id="server-${index}" class="server row" onclick="selectServer(${index})">
      <div class="row flex">
        <span class="material-icons right-margin">storage</span>
        <input
          id="server-input-${index}"
          value="${address}"
          placeholder="Enter Vault Server Address"
          spellcheck="false"
          class="flex"
          ${address ? "disabled" : ''}
        >
      </div>
      <div class="row left-margin">
        <button class="icon right-margin" onclick="editServer(event, ${index})">
          <span id="server-edit-icon-${index}" class="material-icons">${address ? 'edit' : 'save'}</span>
        </button>
        <button id="server-remove-${index}" class="icon" onclick="removeServer(event, ${index})">
          <span class="material-icons">remove_circle</span>
        </button>
      </div>
    </div>
  `).join('');
  document.getElementById('server-list').innerHTML = markup;
};
const addServer = index => { // eslint-disable-line no-unused-vars
  this.preferences.servers.unshift({ address: '' });
  createServerElements();
};
const editServer = (event, index) => { // eslint-disable-line no-unused-vars
  event.stopPropagation();
  const input = document.getElementById(`server-input-${index}`);
  const isSaving = !input.disabled;
  document.getElementById(`server-${index}`).classList[isSaving ? 'add' : 'remove']('pointer');
  document.getElementById(`server-edit-icon-${index}`).innerHTML = isSaving ? 'edit' : 'save';
  input.disabled = isSaving;
  if (isSaving) {
    this.preferences.servers[index].address = input.value;
    writePreferences();
  }
};
const removeServer = (event, index) => { // eslint-disable-line no-unused-vars
  event.stopPropagation();
  document.getElementById(`server-${index}`).remove();
  this.preferences.servers.splice(index, 1);
  writePreferences();
};
const selectServer = index => { // eslint-disable-line no-unused-vars
  const input = document.getElementById(`server-input-${index}`);
  if (input.disabled) {
    ipcRenderer.send('server-selected', input.value);
  }
};

ipcRenderer.on('preferences-path', (event, filePath) => {
  try {
    this.filePath = filePath;
    this.preferences = JSON.parse(fs.readFileSync(filePath));
  } catch (error) {
    console.log(error);
    this.preferences = {};
  }
  if (!this.preferences.servers) {
    this.preferences.servers = [];
  }
  if (!this.preferences.servers.length) {
    this.preferences.servers.push({ address: '' });
  }
  createServerElements();
});
ipcRenderer.send('get-preferences-path');

import Ember from 'ember';
import argTokenizer from 'yargs-parser-tokenizer';
import keys from 'vault/lib/keycodes';

const { inject, computed } = Ember;
const supportedCommands = ['read', 'write', 'list', 'delete'];

export default Ember.Component.extend({
  console: inject.service(),
  inputValue: null,
  commandHistory: computed('log.[]', function() {
    return this.get('log').filterBy('type', 'command');
  }),
  log: computed(function() {
    return [];
  }),
  commandIndex: null,


  handleServiceError(command, method, vaultPath, error) {
    this.pushCommand(command);

    let content;
    let { httpStatus, path } = error;
    let verbClause = {
      'read': 'reading from',
      'write': 'writing to',
      'list': 'listing',
      'delete': 'deleting at'
    }[method];

    content = `Error ${verbClause}: ${vaultPath}.\nURL: ${path}\nCode: ${httpStatus}`;

    if(typeof error.errors[0] === 'string'){
      content = `${content}\nErrors:\n  ${error.errors.join('\n  ')}`;
    }

    this.appendToLog({ type: 'error', content });
  },

  processResponse(response, command, path, method, flags) {
    this.pushCommand(command);
    if (!response) {
      let message = method === 'write' ?
        `Success! Data written to: ${path}` :
        `Success! Data deleted (if it existed) at: ${path}`;

      // print something here
      this.appendToLog({type: 'text', content: message});
      return;
    }
    let { format, field } = flags;
    let secret = response.auth || response.data || response.wrap_info;

    if (field) {
      let fieldValue = secret[field];
      if (fieldValue) {
        switch (typeof fieldValue) {
          case 'string':
            this.appendToLog({type: 'text', content: fieldValue});
            break;
          default:
            this.appendToLog({type: 'object', content: fieldValue});
            break;
        }
      } else {
        this.appendToLog({type: 'error', content: `Field "${field}" not present in secret`});
      }
      return;
    }

    if (format && format === 'json') {
      // just print whole response
      this.appendToLog({type: 'json', content: response});
      return;
    }

    if(method === 'list'){
      this.appendToLog({type: 'list', content: secret});
      return;
    }

    this.appendToLog({type: 'object', content: secret });
  },

  parseCommand(command, shouldThrow) {
    let args = argTokenizer(command);
    if (args[0] === 'vault') {
      args.shift();
    }

    let [method, ...rest] = args;
    let path;
    let flags = [];
    let data = [];

    rest.forEach((arg) => {
      if(arg.startsWith('-')){
        flags.push(arg);
      }
      else{
        if(path){
          data.push(arg);
        }
        else{
          path = arg;
        }
      }
    });

    if(!supportedCommands.includes(method)) {
      if(shouldThrow) {
        throw new Error('invalid command');
      }
      return false;
    }
    return [method, flags, path, data];

  },

  extractDataAndFlags(data, flags) {
    return data.concat(flags).reduce((accumulator, val) => {
      // will be "key=value" or "-flag=value" or "foo=bar=baz"
      // split on the first =
      let [ item, value ] = val.split(/=(.+)/);
      if (item.startsWith('-')) {
        let flagName = item.replace(/^-/, '');
        if (flagName === 'wrap-ttl') {
          flagName = 'wrapTTL';
        }
        accumulator.flags[flagName] = value || true;
        return accumulator;
      }
      // if it exists in data already, then we have multiple
      // foo=bar in the list and need to make it an array
      if (accumulator.data[item]) {
        accumulator.data[item] = [].concat(accumulator.data[item], value);
        return accumulator;
      }
      accumulator.data[item] = value;

      return accumulator;
    }, {data: {}, flags: {}});
  },

  appendToLog(logItem){
    this.get('log').pushObject(logItem);
  },

  pushCommand(command){
    this.set('inputValue', '');
    this.appendToLog({type: 'command', content: command});
    this.set('commandIndex', null)
  },

  executeCommand(command, shouldThrow=false) {
    let serviceArgs;
    // parse to verify it's valid
    try {
      serviceArgs = this.parseCommand(command, shouldThrow);
    } catch (e) {
      this.pushCommand(command);
      this.appendToLog({type: 'help'});
      return;
    }
    // we have a invalid command but don't want to throw
    if (serviceArgs === false) {
      return;
    }

    let [method, flagArray, path, dataArray] = serviceArgs;

    if(path === undefined){
      this.pushCommand(command);
      this.appendToLog({type: 'error', content: 'A path is required to make a request.'});
      return;
    }

    if(dataArray || flagArray) {
      var {data, flags} = this.extractDataAndFlags(dataArray, flagArray);
    }

    if(method === 'write' && !flags.force && dataArray.length === 0){
      this.pushCommand(command);
      this.appendToLog({type: 'error', content: 'Must supply data or use -force'});
      return;
    }
    this.get('console')[method](path, data, flags.wrapTTL)
      .then(resp => this.processResponse(resp, command, path, method, flags))
      .catch((error) => this.handleServiceError(command, method, path, error));
  },

  shiftCommandIndex(keyCode){
    let newInputValue;
    let commandHistory = this.get('commandHistory');
    let commandHistoryLength = commandHistory.length;
    let index = this.get('commandIndex');

    if(keyCode === keys.UP){
      index -= 1;
      if(index < 0){
        index = commandHistoryLength - 1;
      }
    }
    else{
      index += 1;
      if(index === commandHistoryLength){
        newInputValue = "";
      }
      if(index > commandHistoryLength){
        index -= 1;
      }
    }

    if(newInputValue !== ""){
      newInputValue = commandHistory.objectAt(index).content;
    }

    this.set('commandIndex', index);
    this.set('inputValue', newInputValue);
  },

  actions: {
    setValue(val){
      this.set('inputValue', val);
    },
    executeCommand(val){
      this.executeCommand(val, true);
    },
    shiftCommandIndex(direction){
      this.shiftCommandIndex(direction);
    }
  },




});

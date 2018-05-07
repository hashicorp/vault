import Ember from 'ember';
import cliArgParser from 'vault/lib/tokenize-cli-arg-string';


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

  executeCommand(command, shouldThrow=false) {
    let serviceArgs;
    // parse to verify it's valid
    try {
      serviceArgs = this.parseCommand(command, shouldThrow);
    } catch (e) {
      this.set('inputValue', '');
      this.appendToLog({type: 'command', content: command});
      this.appendToLog({type: 'help'});
    }
    // we have a invalid command but don't want to throw
    if (serviceArgs === false) {
      return;
    }
    let [method, path, dataAndFlags] = serviceArgs;
    if (dataAndFlags) {
      var {data, flags} = this.extractDataAndFlags(dataAndFlags);
    }
    this.get('console')[method](path, data, flags.wrapTTL)
      .then(resp => this.processResponse(resp, command, path, method, flags))
      .catch(this.handleServiceError);
  },

  handleServiceError(error) {
    //TODO
    throw error;
  },

  processResponse(response, command, path, method, flags) {
    this.set('inputValue', '');
    this.appendToLog({type: 'command', content: command});
    if (!response) {
      let message = method === 'write' ?
        `Success! Data written to: ${path}` :
        `Success! Data deleted (if it existed) at: ${path}`;

      // print something here
      this.appendToLog({type: 'text', content: message});
      return;
    }
    let { wrapTTL, format, field } = flags;
    let secret = response.data || response.wrap_info;

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

    if (wrapTTL) {
      this.appendToLog({type: 'object', content: response.wrap_info });
      return;
    }

    this.appendToLog({type: 'object', content: response.data });
  },

  parseCommand(command, shouldThrow) {
    let args = cliArgParser(command);
    if (args[0] === 'vault') {
      args.shift();
    }
    let [method, path, ...dataAndFlags] = args;

    if(!supportedCommands.includes(method)) {
      if(shouldThrow) {
        throw new Error('invalid command');
      }
      return false;
    }
    return [method, path, dataAndFlags];

  },

  extractDataAndFlags(dataAndFlags) {
    return dataAndFlags.reduce((accumulator, val) => {
      // will be "key=value" or "-flag=value" or "foo=bar=baz"
      // split on the first =
      let [ item, value ] = val.split(/=(.+)/);
      if (item.startsWith('-')) {
        let flagName = item.replace(/^-/, '');
        if (flagName === 'wrap-ttl') {
          flagName = 'wrapTTL';
        }
        accumulator.flags[flagName] = value;
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

  actions: {
    setValue(val){
      this.set('inputValue', val);
    }
  },




});

import Ember from 'ember';
import cliArgParser from 'vault/lib/tokenize-cli-arg-string';


const { inject } = Ember;
const supportedCommands = ['read', 'write', 'list', 'delete'];

export default Ember.Component.extend({
  console: inject.service(),
  inputValue: null,
  commandHistory: computed(function() {
    return []
  }),

  executeCommand(command) {
    let serviceArgs;
    try {
      serviceArgs = this.parseCommand(command);
    } catch (e) {
      this.renderLog(e);
    }
    // parse to verify it's valid
    // if no error, call console service with the command + flags
    // if no error append to history
    // if no error clear the value
    // render to log (error or command + output)
  },

  parseCommand(command) {
    let args = cliArgParser(parsed);
    if (args[0] === 'vault') {
      args.shift();
    }
    let [method, path, ...dataAndFlags] = args;

    if(!supportedCommands.includes(method)) {
      throw new Error(`${method} is not a supported command, pl`);
    }

  },

  formatResponse(response) {
  },



});

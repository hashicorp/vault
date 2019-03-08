import keys from 'vault/lib/keycodes';
import argTokenizer from 'yargs-parser/lib/tokenize-arg-string.js';

const supportedCommands = ['read', 'write', 'list', 'delete'];
const uiCommands = ['clearall', 'clear', 'fullscreen', 'refresh'];

export function extractDataAndFlags(data, flags) {
  return data.concat(flags).reduce(
    (accumulator, val) => {
      // will be "key=value" or "-flag=value" or "foo=bar=baz"
      // split on the first =
      let [item, value] = val.split(/=(.+)/);
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
    },
    { data: {}, flags: {} }
  );
}

export function executeUICommand(command, logAndOutput, clearLog, toggleFullscreen, refreshFn) {
  const isUICommand = uiCommands.includes(command);
  if (isUICommand) {
    logAndOutput(command);
  }
  switch (command) {
    case 'clearall':
      clearLog(true);
      break;
    case 'clear':
      clearLog();
      break;
    case 'fullscreen':
      toggleFullscreen();
      break;
    case 'refresh':
      refreshFn();
      break;
  }

  return isUICommand;
}

export function parseCommand(command, shouldThrow) {
  // encode everything but spaces
  let cmd = encodeURIComponent(command).replace(/%20/g, decodeURIComponent);
  let args = argTokenizer(cmd);
  if (args[0] === 'vault') {
    args.shift();
  }

  args = args.map(decodeURIComponent);
  let [method, ...rest] = args;
  let path;
  let flags = [];
  let data = [];

  rest.forEach(arg => {
    if (arg.startsWith('-')) {
      flags.push(arg);
    } else {
      if (path) {
        data.push(arg);
      } else {
        path = arg;
      }
    }
  });

  if (!supportedCommands.includes(method)) {
    if (shouldThrow) {
      throw new Error('invalid command');
    }
    return false;
  }
  return [method, flags, path, data];
}

export function logFromResponse(response, path, method, flags) {
  let { format, field } = flags;
  let secret = response && (response.auth || response.data || response.wrap_info);
  if (!secret) {
    if (method === 'write') {
      return { type: 'success', content: `Success! Data written to: ${path}` };
    } else if (method === 'delete') {
      return { type: 'success', content: `Success! Data deleted (if it existed) at: ${path}` };
    } else {
      secret = response;
    }
  }

  if (field) {
    let fieldValue = secret[field];
    let response;
    if (fieldValue) {
      if (format && format === 'json') {
        return { type: 'json', content: fieldValue };
      }
      if (typeof fieldValue == 'string') {
        response = { type: 'text', content: fieldValue };
      } else if (typeof fieldValue == 'number') {
        response = { type: 'text', content: JSON.stringify(fieldValue) };
      } else if (typeof fieldValue == 'boolean') {
        response = { type: 'text', content: JSON.stringify(fieldValue) };
      } else if (Array.isArray(fieldValue)) {
        response = { type: 'text', content: JSON.stringify(fieldValue) };
      } else {
        response = { type: 'object', content: fieldValue };
      }
    } else {
      response = { type: 'error', content: `Field "${field}" not present in secret` };
    }
    return response;
  }

  if (format && format === 'json') {
    // just print whole response
    return { type: 'json', content: response };
  }

  if (method === 'list') {
    return { type: 'list', content: secret };
  }

  return { type: 'object', content: secret };
}

export function logFromError(error, vaultPath, method) {
  let content;
  let { httpStatus, path } = error;
  let verbClause = {
    read: 'reading from',
    write: 'writing to',
    list: 'listing',
    delete: 'deleting at',
  }[method];

  content = `Error ${verbClause}: ${vaultPath}.\nURL: ${path}\nCode: ${httpStatus}`;

  if (typeof error.errors[0] === 'string') {
    content = `${content}\nErrors:\n  ${error.errors.join('\n  ')}`;
  }

  return { type: 'error', content };
}

export function shiftCommandIndex(keyCode, history, index) {
  let newInputValue;
  let commandHistoryLength = history.length;

  if (!commandHistoryLength) {
    return [];
  }

  if (keyCode === keys.UP) {
    index -= 1;
    if (index < 0) {
      index = commandHistoryLength - 1;
    }
  } else {
    index += 1;
    if (index === commandHistoryLength) {
      newInputValue = '';
    }
    if (index > commandHistoryLength) {
      index -= 1;
    }
  }

  if (newInputValue !== '') {
    newInputValue = history.objectAt(index).content;
  }

  return [index, newInputValue];
}

export function logErrorFromInput(path, method, flags, dataArray) {
  if (path === undefined) {
    return { type: 'error', content: 'A path is required to make a request.' };
  }
  if (method === 'write' && !flags.force && dataArray.length === 0) {
    return { type: 'error', content: 'Must supply data or use -force' };
  }
}

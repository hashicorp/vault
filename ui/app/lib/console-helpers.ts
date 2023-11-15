/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import keys from 'core/utils/key-codes';
import AdapterError from '@ember-data/adapter/error';
import { parse } from 'shell-quote';

import argTokenizer from './arg-tokenizer';
import { StringMap } from 'vault/vault/app-types';

// Add new commands to `log-help` component for visibility
const supportedCommands = ['read', 'write', 'list', 'delete', 'kv-get'];
const uiCommands = ['api', 'clearall', 'clear', 'fullscreen', 'refresh'];

interface DataObj {
  [key: string]: string | string[];
}

export function extractDataFromStrings(dataArray: string[]): DataObj {
  if (!dataArray) return {};
  return dataArray.reduce((accumulator: DataObj, val: string) => {
    // will be "key=value" or "foo=bar=baz"
    // split on the first =
    // default to value of empty string
    const [item = '', value = ''] = val.split(/=(.+)?/);
    if (!item) return accumulator;

    // if it exists in data already, then we have multiple
    // foo=bar in the list and need to make it an array
    const existingValue = accumulator[item];
    if (existingValue) {
      accumulator[item] = Array.isArray(existingValue) ? [...existingValue, value] : [existingValue, value];
      return accumulator;
    }
    accumulator[item] = value;
    return accumulator;
  }, {});
}

interface Flags {
  field?: string;
  format?: string;
  force?: boolean;
  wrapTTL?: boolean;
  [key: string]: string | boolean | undefined;
}
export function extractFlagsFromStrings(flagArray: string[], method: string): Flags {
  if (!flagArray) return {};
  return flagArray.reduce((accumulator: Flags, val: string) => {
    // val will be "-flag=value" or "--force"
    // split on the first =
    // default to value or true
    const [item, value] = val.split(/=(.+)?/);
    if (!item) return accumulator;

    let flagName = item.replace(/^-/, '');
    if (flagName === 'wrap-ttl') {
      flagName = 'wrapTTL';
    } else if (method === 'write') {
      if (flagName === 'f' || flagName === '-force') {
        flagName = 'force';
      }
    }
    accumulator[flagName] = value || true;
    return accumulator;
  }, {});
}

interface CommandFns {
  [key: string]: CallableFunction;
}

export function executeUICommand(
  command: string,
  logAndOutput: CallableFunction,
  commandFns: CommandFns
): boolean {
  const cmd = command.startsWith('api') ? 'api' : command;
  const isUICommand = uiCommands.includes(cmd);
  if (isUICommand) {
    logAndOutput(command);
  }
  const execCommand = commandFns[cmd];
  if (execCommand && typeof execCommand === 'function') {
    execCommand();
  }
  return isUICommand;
}

interface ParsedCommand {
  method: string;
  path: string;
  flagArray: string[];
  dataArray: string[];
}
export function parseCommand(command: string): ParsedCommand {
  const args: string[] = argTokenizer(parse(command));
  if (args[0] === 'vault') {
    args.shift();
  }

  const [method = '', ...rest] = args;
  let path = '';
  const flags: string[] = [];
  const data: string[] = [];

  rest.forEach((arg) => {
    if (arg.startsWith('-')) {
      flags.push(arg);
    } else {
      if (path) {
        const strippedArg = arg
          // we'll have arg=something or arg="lol I need spaces", so need to split on the first =
          .split(/=(.+)/)
          // if there were quotes, there's an empty string as the last member in the array that we don't want,
          // so filter it out
          .filter((str) => str !== '')
          // glue the data back together
          .join('=');
        data.push(strippedArg);
      } else {
        path = arg;
      }
    }
  });

  if (!supportedCommands.includes(method)) {
    throw new Error('invalid command');
  }
  return { method, flagArray: flags, path, dataArray: data };
}

interface LogResponse {
  auth?: StringMap;
  data?: StringMap;
  wrap_info?: StringMap;
  [key: string]: unknown;
}

export function logFromResponse(response: LogResponse, path: string, method: string, flags: Flags) {
  const { format, field } = flags;
  const respData: StringMap | undefined = response && (response.auth || response.data || response.wrap_info);
  const secret: StringMap | LogResponse = respData || response;

  if (!respData) {
    if (method === 'write') {
      return { type: 'success', content: `Success! Data written to: ${path}` };
    } else if (method === 'delete') {
      return { type: 'success', content: `Success! Data deleted (if it existed) at: ${path}` };
    }
  }

  if (field) {
    const fieldValue = secret[field];
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

interface CustomError extends AdapterError {
  httpStatus: number;
  path: string;
  errors: string[];
}
export function logFromError(error: CustomError, vaultPath: string, method: string) {
  let content;
  const { httpStatus, path } = error;
  const verbClause = {
    read: 'reading from',
    'kv-get': 'reading secret',
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

interface CommandLog {
  type: string;
  content?: string;
}
export function shiftCommandIndex(keyCode: number, history: CommandLog[], index: number) {
  let newInputValue;
  const commandHistoryLength = history.length;

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
    newInputValue = history.objectAt(index)?.content;
  }

  return [index, newInputValue];
}

export function formattedErrorFromInput(path: string, method: string, flags: Flags, dataArray: string[]) {
  if (path === undefined) {
    return { type: 'error', content: 'A path is required to make a request.' };
  }
  if (method === 'write' && !flags.force && dataArray.length === 0) {
    return { type: 'error', content: 'Must supply data or use -force' };
  }
  return;
}

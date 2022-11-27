#!/usr/bin/env node
/* eslint-env node */

// Read more: https://github.com/facebook/jscodeshift#parser

const babylon = require('@babel/parser');

const parserConfig = {
  sourceType: 'module',
  allowImportExportEverywhere: true,
  allowReturnOutsideFunction: true,
  startLine: 1,
  tokens: true,
  plugins: [
    ['flow', { all: true }],
    'flowComments',
    'jsx',
    'asyncGenerators',
    'bigInt',
    'classProperties',
    'classPrivateProperties',
    'classPrivateMethods',
    'decorators-legacy', // allows decorator to come before export statement
    'doExpressions',
    'dynamicImport',
    'exportDefaultFrom',
    'exportNamespaceFrom',
    'functionBind',
    'functionSent',
    'importMeta',
    'logicalAssignment',
    'nullishCoalescingOperator',
    'numericSeparator',
    'objectRestSpread',
    'optionalCatchBinding',
    'optionalChaining',
    ['pipelineOperator', { proposal: 'minimal' }],
    'throwExpressions',
  ],
};

export default {
  parse: function (source) {
    return babylon.parse(source, parserConfig);
  },
};

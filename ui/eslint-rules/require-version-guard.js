/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

'use strict';

/**
 * ESLint rule: require-version-guard
 *
 * Enforces that every call to `.versioned()` on the api service is wrapped
 * inside an `if`/`else if` block whose condition contains a call to
 * `.hasFeature(...)` on the version service.
 *
 * This prevents the `any`-typed escape hatch exposed by `ApiService.versioned()`
 * from being used outside a runtime version gate, where it would be unsafe.
 *
 * Valid:
 *   if (this.version.hasFeature('agents')) {
 *     this.api.versioned().sys.newMethod();
 *   }
 *
 * Invalid:
 *   this.api.versioned().sys.newMethod();           // no guard at all
 *   if (someOtherCondition) {
 *     this.api.versioned().sys.newMethod();         // wrong guard function
 *   }
 */

/** @type {import('eslint').Rule.RuleModule} */
module.exports = {
  meta: {
    type: 'problem',
    docs: {
      description: 'Require that calls to .versioned() are inside a this.version.hasFeature(...) guard',
      category: 'Best Practices',
      recommended: false,
      url: 'ui/docs/versioned-api-calls.md',
    },
    messages: {
      missingGuard:
        '`.versioned()` must be called inside an `if (this.version.hasFeature(...))` block. ' +
        'See ui/docs/versioned-api-calls.md for the version-gating pattern.',
    },
    schema: [],
  },

  create(context) {
    /**
     * Recursively walk a node and its descendants to determine whether any
     * CallExpression calls a function named `hasFeature`.
     */
    function containsHasFeature(node) {
      if (!node) return false;

      if (
        node.type === 'CallExpression' &&
        node.callee.type === 'MemberExpression' &&
        node.callee.property.type === 'Identifier' &&
        node.callee.property.name === 'hasFeature'
      ) {
        return true;
      }

      // Recurse into the child nodes most likely to contain a nested call
      for (const key of ['left', 'right', 'test', 'callee', 'object', 'argument', 'expressions']) {
        const child = node[key];
        if (!child) continue;
        if (Array.isArray(child)) {
          if (child.some(containsHasFeature)) return true;
        } else if (typeof child === 'object' && child.type) {
          if (containsHasFeature(child)) return true;
        }
      }

      return false;
    }

    /**
     * Walk up the ancestor chain. Return true if any enclosing IfStatement's
     * test (or a logical sub-expression thereof) contains `hasFeature`.
     */
    function isInsideHasFeatureGuard(node) {
      let current = node.parent;

      while (current) {
        if (current.type === 'IfStatement') {
          if (containsHasFeature(current.test)) {
            return true;
          }
        }
        current = current.parent;
      }

      return false;
    }

    return {
      CallExpression(node) {
        // Match any <expr>.versioned() call
        if (
          node.callee.type === 'MemberExpression' &&
          node.callee.property.type === 'Identifier' &&
          node.callee.property.name === 'versioned'
        ) {
          if (!isInsideHasFeatureGuard(node)) {
            context.report({
              node,
              messageId: 'missingGuard',
            });
          }
        }
      },
    };
  },
};

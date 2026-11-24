/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

'use strict';

// eslint-disable-next-line n/no-extraneous-require
const { RuleTester } = require('eslint');
const rule = require('../require-version-guard');

const ruleTester = new RuleTester({
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
  },
});

ruleTester.run('require-version-guard', rule, {
  valid: [
    // Direct guard at the top of the if block
    {
      code: `
        if (this.version.hasFeature('agents')) {
          this.api.versioned().sys.newMethod();
        }
      `,
    },
    // Guard on a standalone variable reference
    {
      code: `
        if (this.version.hasFeature('agents')) {
          const v = this.api.versioned();
          v.sys.newMethod();
        }
      `,
    },
    // versioned() inside a nested block within the gate
    {
      code: `
        if (this.version.hasFeature('agents')) {
          try {
            this.api.versioned().sys.newMethod();
          } catch(e) {}
        }
      `,
    },
    // hasFeature is part of a logical && expression
    {
      code: `
        if (this.isEnterprise && this.version.hasFeature('agents')) {
          this.api.versioned().sys.newMethod();
        }
      `,
    },
    // hasFeature on the left side of &&
    {
      code: `
        if (this.version.hasFeature('agents') && someOtherCondition) {
          this.api.versioned().sys.newMethod();
        }
      `,
    },
    // versioned() used inside an else-if that has hasFeature
    {
      code: `
        if (something) {
          doSomething();
        } else if (this.version.hasFeature('agents')) {
          this.api.versioned().sys.newMethod();
        }
      `,
    },
    // Non-api versioned() calls (e.g. a different object's method named versioned) are fine
    // as long as they are inside a hasFeature guard — rule only errors when NOT in guard
    {
      code: `
        if (this.version.hasFeature('agents')) {
          someOtherObject.versioned().doSomething();
        }
      `,
    },
  ],

  invalid: [
    // Called at module top level — no guard at all
    {
      code: `this.api.versioned().sys.newMethod();`,
      errors: [{ messageId: 'missingGuard' }],
    },
    // Called inside an if with an unrelated condition
    {
      code: `
        if (someOtherCondition) {
          this.api.versioned().sys.newMethod();
        }
      `,
      errors: [{ messageId: 'missingGuard' }],
    },
    // Called inside an if (isEnterprise) — not a hasFeature guard
    {
      code: `
        if (this.version.isEnterprise) {
          this.api.versioned().sys.newMethod();
        }
      `,
      errors: [{ messageId: 'missingGuard' }],
    },
    // Called inside an if with a different method named similarly
    {
      code: `
        if (this.version.hasFlag('agents')) {
          this.api.versioned().sys.newMethod();
        }
      `,
      errors: [{ messageId: 'missingGuard' }],
    },
    // versioned() outside any if block — assigned to a variable
    {
      code: `const v = this.api.versioned();`,
      errors: [{ messageId: 'missingGuard' }],
    },
    // versioned() in a function body that itself is not in a hasFeature guard
    {
      code: `
        function doWork() {
          this.api.versioned().sys.newMethod();
        }
      `,
      errors: [{ messageId: 'missingGuard' }],
    },
  ],
});

// eslint-disable-next-line no-console
console.log('require-version-guard: all tests passed');

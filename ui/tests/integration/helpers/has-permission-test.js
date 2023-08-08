/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';
import sinon from 'sinon';

const Permissions = Service.extend({
  globPaths: null,
  hasNavPermission() {
    return this.globPaths ? true : false;
  },
});

module('Integration | Helper | has-permission', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:permissions', Permissions);
    this.permissions = this.owner.lookup('service:permissions');
  });

  test('it renders', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);

    assert.dom(this.element).hasText('No');
    await run(() => {
      this.permissions.set('globPaths', { 'test/': { capabilities: ['update'] } });
    });
    await settled();
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when globPaths changes');
  });

  test('it should pass args from helper to service method', async function (assert) {
    const stub = sinon.stub(this.permissions, 'hasNavPermission').returns(true);
    this.permissions.set('exactPaths', {
      'sys/auth': {
        capabilities: ['read'],
      },
      'identity/mfa/method': {
        capabilities: ['read'],
      },
    });
    this.routeParams = ['methods', 'mfa'];

    await render(hbs`
      {{#if (has-permission "access" routeParams=this.routeParams requireAll=true)}}
        Yes
      {{else}}
        No
      {{/if}}
    `);

    assert.true(
      stub.withArgs('access', this.routeParams, true).calledOnce,
      'Args are passed from helper to service'
    );
    assert.dom(this.element).hasText('Yes', 'Helper returns value from service method');
    stub.restore();
  });
});

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { tracked } from '@glimmer/tracking';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';
import sinon from 'sinon';

class PermissionsService extends Service {
  @tracked globPaths = null;
  @tracked exactPaths = null;
  @tracked canViewAll = false;
  @tracked chrootNamespace = null;

  hasNavPermission(...args) {
    const [route, routeParams, requireAll] = args;
    void (typeof route === 'string');
    void (routeParams === undefined || Array.isArray(routeParams));
    void (requireAll === undefined || typeof requireAll === 'boolean');

    if (this.canViewAll) return true;
    if (this.globPaths || this.exactPaths) return true;
    return false;
  }
}

class NamespaceService extends Service {
  @tracked path = '';
}

module('Integration | Helper | has-permission', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:permissions', PermissionsService);
    this.owner.register('service:namespace', NamespaceService);
    this.permissions = this.owner.lookup('service:permissions');
    this.namespace = this.owner.lookup('service:namespace');
  });

  test('it recomputes when globPaths change', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No');
    this.permissions.globPaths = { 'test/': { capabilities: ['update'] } };
    await settled();
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when globPaths changes');
  });

  test('it recomputes when exactPaths changes', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No');
    this.permissions.exactPaths = { 'test/': { capabilities: ['update'] } };
    await settled();
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when exactPaths changes');
  });

  test('it recomputes when canViewAll changes', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No');
    this.permissions.canViewAll = true;
    await settled();
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when canViewAll changes');
  });

  test('it recomputes when chrootNamespace changes', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No');
    this.permissions.chrootNamespace = 'admin';
    this.permissions.globPaths = { 'test/': { capabilities: ['update'] } };
    await settled();
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when chrootNamespace changes');
  });

  test('it recomputes when namespace.path changes', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No');
    this.namespace.path = 'new-ns';
    this.permissions.globPaths = { 'test/': { capabilities: ['update'] } };
    await settled();
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when namespace.path changes');
  });

  test('it should pass args from helper to service method', async function (assert) {
    const stub = sinon.stub(this.permissions, 'hasNavPermission').returns(true);

    // seed some ACL so the helper would logically return true
    this.permissions.exactPaths = {
      'sys/auth': { capabilities: ['read'] },
      'identity/mfa/method': { capabilities: ['read'] },
    };

    // strict-mode: define on context and use `this.` in template
    this.routeParams = ['methods', 'mfa'];

    await render(hbs`
    {{if (has-permission "access" routeParams=this.routeParams requireAll=true) "Yes" "No"}}
  `);

    // use deep/loose match for the array arg
    assert.true(
      stub.calledWithMatch('access', ['methods', 'mfa'], true),
      'Args are passed from helper to service'
    );

    assert.dom(this.element).hasText('Yes', 'Helper returns value from service method');
    stub.restore();
  });

  test('returns false for missing route', async function (assert) {
    const stub = sinon.stub(this.permissions, 'hasNavPermission').returns(false);
    await render(hbs`{{#if (has-permission "missing")}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No', 'Helper returns false for missing route');
    stub.restore();
  });

  test('returns false for undefined params', async function (assert) {
    const stub = sinon.stub(this.permissions, 'hasNavPermission').returns(false);
    await render(hbs`{{#if (has-permission "access" routeParams=undefined)}}Yes{{else}}No{{/if}}`);
    assert.dom(this.element).hasText('No', 'Helper returns false for undefined params');
    stub.restore();
  });
});

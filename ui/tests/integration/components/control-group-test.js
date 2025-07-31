/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { CONTROL_GROUP } from 'vault/tests/helpers/components/control-group-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const controlGroupService = Service.extend({
  init() {
    this._super(...arguments);
    this.set('wrapInfo', null);
  },
  wrapInfoForAccessor() {
    return this.wrapInfo;
  },
});

const authService = Service.extend();

module('Integration | Component | control group', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:auth', authService);
    this.owner.register('service:control-group', controlGroupService);
    this.controlGroup = this.owner.lookup('service:controlGroup');
    this.auth = this.owner.lookup('service:auth');
  });

  const setup = (modelData = {}, authData = {}) => {
    const modelDefaults = {
      approved: false,
      requestPath: 'foo/bar',
      id: 'accessor',
      requestEntity: { id: 'requestor', name: 'entity8509' },
      reload: sinon.stub(),
    };
    const authDataDefaults = { entityId: 'requestor' };

    return {
      model: {
        ...modelDefaults,
        ...modelData,
      },
      authData: {
        ...authDataDefaults,
        ...authData,
      },
    };
  };

  test('requestor rendering', async function (assert) {
    const { model, authData } = setup();
    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);
    assert.dom(GENERAL.inlineAlert).exists('shows accessor callout');
    assert.dom(CONTROL_GROUP.bannerPrefix).hasText('Locked');
    assert.dom(CONTROL_GROUP.bannerText).hasText('The path you requested is locked by a Control Group');
    assert.dom(CONTROL_GROUP.requestorText).hasText(`You are requesting access to ${model.requestPath}`);
    assert.dom(CONTROL_GROUP.tokenValue).doesNotExist('does not show token message when there is no token');
    assert.dom(GENERAL.button('Refresh')).exists('shows refresh button');
    assert.dom(CONTROL_GROUP.authorizations).hasText('Awaiting authorization.');
  });

  test('requestor rendering: with token', async function (assert) {
    const { model, authData } = setup();
    this.set('model', model);
    this.set('auth.authData', authData);
    this.set('controlGroup.wrapInfo', { token: 'token' });
    await render(hbs`<ControlGroup @model={{this.model}} />`);
    assert.dom(CONTROL_GROUP.tokenValue).exists('shows token message');
    assert.dom(CONTROL_GROUP.tokenValue).hasText('token', 'shows token value');
  });

  test('requestor rendering: some approvals', async function (assert) {
    const { model, authData } = setup({ authorizations: [{ name: 'manager 1' }, { name: 'manager 2' }] });
    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);
    assert.dom(CONTROL_GROUP.authorizations).hasText('Already approved by manager 1, manager 2');
  });

  test('requestor rendering: approved with no token', async function (assert) {
    const { model, authData } = setup({ approved: true });
    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);

    assert.dom(CONTROL_GROUP.bannerPrefix).hasText('Success!');
    assert.dom(CONTROL_GROUP.bannerText).hasText('You have been given authorization');
    assert.dom(CONTROL_GROUP.tokenValue).doesNotExist('does not show token message when there is no token');
    assert.dom(GENERAL.button('Refresh')).doesNotExist('does not shows refresh button');
    assert.dom(CONTROL_GROUP.successComponent).exists('renders control group success');
  });

  test('requestor rendering: approved with token', async function (assert) {
    const { model, authData } = setup({ approved: true });
    this.set('model', model);
    this.set('auth.authData', authData);
    this.set('controlGroup.wrapInfo', { token: 'token' });
    await render(hbs`<ControlGroup @model={{this.model}} />`);
    assert.dom(CONTROL_GROUP.tokenValue).exists('shows token');
    assert.dom(GENERAL.button('Refresh')).doesNotExist('does not shows refresh button');
    assert.dom(CONTROL_GROUP.successComponent).exists('renders control group success');
  });

  test('authorizer rendering', async function (assert) {
    const { model, authData } = setup({ canAuthorize: true }, { entityId: 'manager' });

    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);

    assert.dom(CONTROL_GROUP.bannerPrefix).hasText('Locked');
    assert
      .dom(CONTROL_GROUP.bannerText)
      .hasText('Someone is requesting access to a path locked by a Control Group');
    assert
      .dom(CONTROL_GROUP.requestorText)
      .hasText(`${model.requestEntity.name} is requesting access to ${model.requestPath}`);
    assert.dom(CONTROL_GROUP.tokenValue).doesNotExist('does not show token message when there is no token');

    assert.dom(GENERAL.button('Authorize')).exists('shows authorize button');
  });

  test('authorizer rendering:authorized', async function (assert) {
    const { model, authData } = setup(
      { canAuthorize: true, authorizations: [{ id: 'manager', name: 'manager' }] },
      { entityId: 'manager' }
    );

    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);

    assert.dom(CONTROL_GROUP.bannerPrefix).hasText('Thanks!');
    assert.dom(CONTROL_GROUP.bannerText).hasText('You have given authorization');
    assert.dom(GENERAL.backButton).exists('back link is visible');
  });

  test('authorizer rendering: authorized and success', async function (assert) {
    const { model, authData } = setup(
      { approved: true, canAuthorize: true, authorizations: [{ id: 'manager', name: 'manager' }] },
      { entityId: 'manager' }
    );

    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);

    assert.dom(CONTROL_GROUP.bannerPrefix).hasText('Thanks!');
    assert.dom(CONTROL_GROUP.bannerText).hasText('You have given authorization');
    assert.dom(GENERAL.backButton).exists('back link is visible');
    assert
      .dom(CONTROL_GROUP.requestorText)
      .hasText(`${model.requestEntity.name} is authorized to access ${model.requestPath}`);
    assert.dom(CONTROL_GROUP.successComponent).doesNotExist('does not render control group success');
  });

  test('third-party: success', async function (assert) {
    const { model, authData } = setup(
      { approved: true, canAuthorize: true, authorizations: [{ id: 'foo', name: 'foo' }] },
      { entityId: 'manager' }
    );

    this.set('model', model);
    this.set('auth.authData', authData);
    await render(hbs`<ControlGroup @model={{this.model}} />`);
    assert.dom(CONTROL_GROUP.bannerPrefix).hasText('Success!');
    assert.dom(CONTROL_GROUP.bannerText).hasText('This Control Group has been authorized');
    assert.dom(GENERAL.backButton).exists('back link is visible');
    assert.dom(CONTROL_GROUP.successComponent).doesNotExist('does not render control group success');
  });
});

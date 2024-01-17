/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later, run, _cancelTimers as cancelTimers } from '@ember/runloop';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import controlGroupSuccess from '../../pages/components/control-group-success';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const component = create(controlGroupSuccess);

const controlGroupService = Service.extend({
  deleteControlGroupToken: sinon.stub(),
  markTokenForUnwrap: sinon.stub(),
});

const storeService = Service.extend({
  adapterFor() {
    return {
      toolAction() {
        return resolve({ data: { foo: 'bar' } });
      },
    };
  },
});

module('Integration | Component | control group success', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:control-group', controlGroupService);
      this.controlGroup = this.owner.lookup('service:control-group');
      this.owner.register('service:store', storeService);
      this.router = this.owner.lookup('service:router');
      this.router.reopen({
        transitionTo: sinon.stub().returns(resolve()),
      });
      setRunOptions({
        rules: {
          // TODO: swap out JsonEditor with Hds::CodeBlock for display
          'color-contrast': { enabled: false },
          label: { enabled: false },
        },
      });
    });
  });

  const MODEL = {
    approved: false,
    requestPath: 'foo/bar',
    id: 'accessor',
    requestEntity: { id: 'requestor', name: 'entity8509' },
    reload: sinon.stub(),
  };
  test('render with saved token', async function (assert) {
    assert.expect(3);
    const response = {
      uiParams: { url: '/foo' },
      token: 'token',
    };
    this.set('model', MODEL);
    this.set('response', response);
    await render(hbs`<ControlGroupSuccess @model={{this.model}} @controlGroupResponse={{this.response}} />`);
    assert.ok(component.showsNavigateMessage, 'shows unwrap message');
    await component.navigate();
    later(() => cancelTimers(), 50);
    return settled().then(() => {
      assert.ok(this.controlGroup.markTokenForUnwrap.calledOnce, 'marks token for unwrap');
      assert.ok(this.router.transitionTo.calledOnce, 'calls router transition');
    });
  });

  test('render without token', async function (assert) {
    assert.expect(2);
    this.set('model', MODEL);
    await render(hbs`<ControlGroupSuccess @model={{this.model}} />`);
    assert.ok(component.showsUnwrapForm, 'shows unwrap form');
    await component.token('token');
    component.unwrap();
    later(() => cancelTimers(), 50);
    return settled().then(() => {
      assert.ok(component.showsJsonViewer, 'shows unwrapped data');
    });
  });
});

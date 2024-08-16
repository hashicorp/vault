/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, typeIn, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const handler = (e) => {
  // required because filter-input-explicit passes handleSearch on form submit
  if (e && e.preventDefault) e.preventDefault();
  return;
};

module('Integration | Component | filter-input-explicit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.handleSearch = sinon.spy(handler);
    this.handleInput = sinon.spy();
    this.handleKeyDown = sinon.spy();
    this.query = '';
    this.placeholder = 'Filter roles';

    this.renderComponent = () => {
      return render(
        hbs`<FilterInputExplicit aria-label="test-component" @placeholder={{this.placeholder}} @query={{this.query}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`
      );
    };
  });

  test('it renders', async function (assert) {
    this.query = 'foo';
    await this.renderComponent();

    assert
      .dom(GENERAL.filterInputExplicit)
      .hasAttribute('placeholder', 'Filter roles', 'Placeholder passed to input element');
    assert.dom(GENERAL.filterInputExplicit).hasValue('foo', 'Value passed to input element');
  });

  test('it should call handleSearch on submit', async function (assert) {
    await this.renderComponent();
    await typeIn(GENERAL.filterInputExplicit, 'bar');
    await click(GENERAL.filterInputExplicitSearch);
    assert.ok(this.handleSearch.calledOnce, 'handleSearch was called once');
  });

  test('it should send keydown event on keydown', async function (assert) {
    await this.renderComponent();
    await typeIn(GENERAL.filterInputExplicit, 'a');
    await typeIn(GENERAL.filterInputExplicit, 'b');

    assert.ok(this.handleKeyDown.calledTwice, 'handle keydown was called twice');
    assert.ok(this.handleSearch.notCalled, 'handleSearch was not called on a keydown event');
  });
});

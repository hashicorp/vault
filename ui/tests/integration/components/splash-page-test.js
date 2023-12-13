/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | splash-page', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render', async function (assert) {
    assert.expect(4);
    await render(hbs`<SplashPage>
    <:header>
      Header here
    </:header>
    <:subHeader>
      sub header
    </:subHeader>
    <:content>
      content
    </:content>
    <:footer>
      <div data-test-footer>footer</div>
    </:footer>

    </SplashPage>
      `);
    assert.dom('[data-test-splash-page-header]').includesText('Header here', 'Header renders');
    assert.dom('[data-test-splash-page-sub-header]').includesText('sub header', 'SubHeader renders');
    assert.dom('[data-test-splash-page-content]').includesText('content', 'Content renders');
    assert.dom('[data-test-footer]').includesText('footer', 'Footer renders');
  });
});

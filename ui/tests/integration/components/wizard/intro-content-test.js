/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | wizard/intro-content', function (hooks) {
  setupRenderingTest(hooks);

  test('it uses the dark image in dark mode', async function (assert) {
    const theme = this.owner.lookup('service:theme');
    sinon.stub(theme, 'isDarkMode').value(true);

    await render(hbs`
      <Wizard::IntroContent
        @description="Description"
        @imageAlt="Diagram"
        @imageCaption="Caption"
        @imageSrc="light.png"
      />
    `);

    assert.dom('[data-test-intro-content-image]').hasAttribute('src', 'light-dark.png');
  });

  test('it uses the light image outside dark mode', async function (assert) {
    const theme = this.owner.lookup('service:theme');
    sinon.stub(theme, 'isDarkMode').value(false);

    await render(hbs`
      <Wizard::IntroContent
        @description="Description"
        @imageAlt="Diagram"
        @imageCaption="Caption"
        @imageSrc="light.png"
      />
    `);

    assert.dom('[data-test-intro-content-image]').hasAttribute('src', 'light.png');
  });
});

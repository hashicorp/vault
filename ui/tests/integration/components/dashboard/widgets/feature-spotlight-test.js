/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | dashboard/widgets/feature-spotlight', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.title = 'New Agent Registry in Vault';
    this.description =
      'Centralized backend for registering AI agents and their associated human owners in Vault.';
    this.link = '/vault/docs/agent-registry';
    this.imageSrc = '/ui/images/agent-registry-dashboard.png';

    this.renderComponent = async () => {
      return render(hbs`
        <Dashboard::Widgets::FeatureSpotlight
          @title={{this.title}}
          @description={{this.description}}
          @link={{this.link}}
          @imageSrc={{this.imageSrc}}
        />
      `);
    };
  });

  test('it renders the title, description, learn more link, and image', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.textDisplay(this.title)).hasText(this.title);

    assert.dom(GENERAL.textBody('feature-spotlight-description')).hasText(this.description);

    assert.dom(GENERAL.linkTo('feature-spotlight-learn-more')).exists('Learn more link is rendered');
    assert
      .dom(GENERAL.linkTo('feature-spotlight-learn-more'))
      .hasText('Learn more', 'learn more link has expected text');

    assert
      .dom(`${GENERAL.cardContainer('feature-spotlight')} img`)
      .exists('image is rendered')
      .hasAttribute('alt', this.title);
  });

  test('it does not render an image when @imageSrc is not provided', async function (assert) {
    this.imageSrc = null;
    await this.renderComponent();
    assert
      .dom(`${GENERAL.cardContainer('feature-spotlight')} img`)
      .doesNotExist('image is not rendered without @imageSrc');
  });

  test('it uses the @link arg to build the Learn more href', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.linkTo('feature-spotlight-learn-more'))
      .hasAttribute(
        'href',
        'https://developer.hashicorp.com/vault/docs/agent-registry',
        'href is the full doc-link URL'
      );
  });

  test('it renders with custom @title and @description args', async function (assert) {
    this.title = 'My Custom Feature';
    this.description = 'My custom description.';
    this.link = '/vault/docs/my-feature';
    await this.renderComponent();

    assert.dom(GENERAL.textDisplay('My Custom Feature')).hasText('My Custom Feature');
    assert.dom(GENERAL.textBody('feature-spotlight-description')).hasText('My custom description.');
  });

  test('it renders the card container', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('feature-spotlight')).exists('card container is rendered');
  });
});

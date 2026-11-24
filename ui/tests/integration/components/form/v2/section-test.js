/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';

module('Integration | Component | form/v2/section', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders section title, description, and yielded content', async function (assert) {
    this.section = {
      name: 'test-section',
      title: 'User Settings',
      description: 'Configure your user preferences',
    };

    await render(hbs`
      <Hds::Form as |Form|>
        <Form::V2::Section @section={{this.section}} @formComponent={{Form}}>
          <div data-test-content>Settings fields</div>
        </Form::V2::Section>
      </Hds::Form>
    `);

    assert.dom('h2').hasText('User Settings', 'renders section title');
    assert.dom('p').includesText('Configure your user preferences', 'renders section description');
    assert.dom('[data-test-content]').hasText('Settings fields', 'renders yielded content');
  });

  test('it renders without section header when title is not provided', async function (assert) {
    this.section = {
      name: 'no-title-section',
    };

    await render(hbs`
      <Hds::Form as |Form|>
        <Form::V2::Section @section={{this.section}} @formComponent={{Form}}>
          <div data-test-content>Content without title</div>
        </Form::V2::Section>
      </Hds::Form>
    `);

    assert.dom('h2').doesNotExist('does not render title');
    assert.dom('p').doesNotExist('does not render description');
    assert.dom('[data-test-content]').hasText('Content without title', 'renders yielded content');
  });

  test('it renders an empty section', async function (assert) {
    this.section = {
      name: 'empty-section',
      title: 'Empty Section',
    };

    await render(hbs`
      <Hds::Form as |Form|>
        <Form::V2::Section @section={{this.section}} @formComponent={{Form}} />
      </Hds::Form>
    `);

    assert.dom('h2').hasText('Empty Section', 'renders section title');
  });
});

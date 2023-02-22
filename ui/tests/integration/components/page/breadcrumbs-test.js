/* eslint-disable ember/no-private-routing-service */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | page/breadcrumbs', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.breadcrumbs = [
      { label: 'Home', route: 'home', linkExternal: true },
      { label: 'Details', route: 'home.details' },
      { label: 'Edit item' },
    ];

    await render(hbs`<Page::Breadcrumbs @breadcrumbs={{this.breadcrumbs}} />`);
    assert.dom('[data-test-breadcrumbs]').exists('renders passed in breadcrumbs');
    assert.strictEqual(findAll('[data-test-breadcrumbs] li').length, 3, 'it renders 3 breadcrumbs');
    assert.strictEqual(
      findAll('[data-test-breadcrumbs] a').length,
      2,
      'it does not render a link if no path'
    );
  });
});

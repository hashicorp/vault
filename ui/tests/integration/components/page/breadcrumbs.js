import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | breadcrumbs', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.breadcrumbs = [
      { label: 'Home', path: 'home', linkExternal: true },
      { label: 'Details', path: 'overview' },
      { label: 'Edit item' },
    ];
    await render(hbs`<Page::BreadcrumbHeader @breadcrumbs={{this.breadcrumbs}} />`);

    assert.dom('[data-test-breadcrumbs]').exists('renders passed in breadcrumbs');
  });
});

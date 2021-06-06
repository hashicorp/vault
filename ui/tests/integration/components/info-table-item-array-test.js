import { module, test } from 'qunit';
import Service from '@ember/service';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';

const DISPLAY_ARRAY = ['role-1', 'role-2', 'role-3', 'role-4', 'role-5'];

const storeService = Service.extend({
  query() {
    return new Promise(resolve => {
      resolve([
        { id: 'role-1' },
        { id: 'role-2' },
        { id: 'role-3' },
        { id: 'role-4' },
        { id: 'role-5' },
        { id: 'role-6' },
      ]);
    });
  },
});

module('Integration | Component | InfoTableItemArray', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('displayArray', DISPLAY_ARRAY);
    this.set('isLink', true);
    this.set('modelType', 'transform/role');
    this.set('queryParam', 'role');
    this.set('backend', 'transform');
    this.set('wildcardLabel', 'role');
    this.set('viewAll', 'roles');
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
    });
  });

  hooks.afterEach(function() {
    this.owner.unregister('service:store');
  });

  test('it renders', async function(assert) {
    await render(hbs`<InfoTableItemArray
        @displayArray={{displayArray}}
      />`);

    assert.dom('[data-test-info-table-item-array]').exists();
    let noLinkString = document.querySelector('code').textContent;

    assert.equal(
      noLinkString.length,
      DISPLAY_ARRAY.toString().length,
      'renders a string of the array if isLink is not provided'
    );
  });

  test('it renders links if isLink is true', async function(assert) {
    await render(hbs`<InfoTableItemArray
      @displayArray={{displayArray}}
      @isLink={{isLink}}
      @modelType={{modelType}}
      @queryParam={{queryParam}}
      @backend={{backend}}
    />`);
    assert.equal(
      document.querySelectorAll('a > span').length,
      DISPLAY_ARRAY.length,
      'renders each item in array with link'
    );
  });

  test('it renders a badge and view all if wildcard in display array && < 10', async function(assert) {
    const displayArrayWithWildcard = ['role-1', 'role-2', 'role-3', 'r*'];
    this.set('displayArrayWithWildcard', displayArrayWithWildcard);
    await render(hbs`<InfoTableItemArray
      @displayArray={{displayArrayWithWildcard}}
      @isLink={{isLink}}
      @modelType={{modelType}}
      @queryParam={{queryParam}}
      @backend={{backend}}
      @viewAll={{viewAll}}
    />`);

    assert.equal(
      document.querySelectorAll('a > span').length,
      DISPLAY_ARRAY.length - 1,
      'renders each item in array with link'
    );
    // 6 here comes from the six roles setup in the store service.
    assert.dom('[data-test-count="6"]').exists('correctly counts with wildcard filter and shows the count');
    assert.dom('[data-test-view-all="roles"]').exists({ count: 1 }, 'renders 1 view all roles');
  });

  test('it renders a badge and view all if wildcard in display array && >= 10', async function(assert) {
    const displayArrayWithWildcard = [
      'role-1',
      'role-2',
      'role-3',
      'r*',
      'role-4',
      'role-5',
      'role-6',
      'role-7',
      'role-8',
      'role-9',
      'role-10',
    ];
    this.set('displayArrayWithWildcard', displayArrayWithWildcard);
    await render(hbs`<InfoTableItemArray
      @displayArray={{displayArrayWithWildcard}}
      @isLink={{isLink}}
      @modelType={{modelType}}
      @queryParam={{queryParam}}
      @backend={{backend}}
      @viewAll={{viewAll}}
    />`);
    const numberCutOffTruncatedArray = displayArrayWithWildcard.length - 5;
    assert.equal(document.querySelectorAll('a > span').length, 5, 'renders truncated array of five');
    assert
      .dom(`[data-test-and="${numberCutOffTruncatedArray}"]`)
      .exists('correctly counts with wildcard filter and shows the count');
  });
});

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, triggerEvent, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

let file;
const data = JSON.stringify({ some: 'content' }, null, 2);
const fileEvent = () => {
  file = new Blob([data], { type: 'application/json' });
  file.name = 'file.json';
  return ['change', { files: [file] }];
};

module('Integration | Component | pgp list', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.set('listLength', 0);
    await render(hbs`<PgpList @listLength={{this.listLength}} />`);
    assert.dom('[data-test-empty-text]').exists('shows the empty state');
    this.set('listLength', 1);
    assert
      .dom('[data-test-pgp-file]')
      .exists({ count: 1 }, 'renders pgp-file one component when length is updated');
    this.set('listLength', 2);
    assert
      .dom('[data-test-pgp-file]')
      .exists({ count: 2 }, 'renders multiple pgp-file components when length is updated');
  });

  test('onDataUpdate is called properly', async function (assert) {
    this.set('spy', sinon.spy());
    const event = fileEvent();

    await render(hbs`<PgpList @listLength={{1}} @onDataUpdate={{this.spy}} />`);
    triggerEvent('[data-test-pgp-file-input]', ...event);

    // FileReader is async, but then we need extra run loop wait to re-render
    await waitUntil(() => {
      return !!this.spy.calledOnce;
    });
    const expected = [btoa(data)];
    assert.deepEqual(
      this.spy.getCall(0).args[0],
      expected,
      'calls onchange with an array of base64 converted files'
    );
  });

  test('sparse filling of multiple files, then shortening', async function (assert) {
    this.set('spy', sinon.spy());
    this.set('listLength', 3);

    await render(hbs`<PgpList @listLength={{this.listLength}} @onDataUpdate={{this.spy}} />`);

    // add a file to the first input
    triggerEvent('[data-test-pgp-file]:nth-child(1) [data-test-pgp-file-input]', ...fileEvent());
    await waitUntil(() => {
      return !!this.spy.calledOnce;
    });
    const expected = [btoa(data), '', ''];
    assert.deepEqual(
      this.spy.getCall(0).args[0],
      expected,
      'calls onchange with an array of base64 converted files'
    );

    // add a file to the third input
    triggerEvent('[data-test-pgp-file]:nth-child(3) [data-test-pgp-file-input]', ...fileEvent());
    await waitUntil(() => {
      return !!this.spy.calledTwice;
    });
    const expected2 = [btoa(data), '', btoa(data)];
    assert.deepEqual(
      this.spy.getCall(1).args[0],
      expected2,
      'calls onchange with an array 2 base64 converted files'
    );

    // this will trim off the last input which was filled, so we should only have one file in the array
    // now
    this.set('listLength', 2);
    const expected3 = [btoa(data), ''];
    assert.deepEqual(
      this.spy.getCall(2).args[0],
      expected3,
      'shortens the list with an array with one base64 converted files'
    );
  });

  test('sparse filling of multiple files, then lengthening', async function (assert) {
    this.set('spy', sinon.spy());
    this.set('listLength', 3);

    await render(hbs`<PgpList @listLength={{this.listLength}} @onDataUpdate={{this.spy}} />`);

    // add a file to the first input
    triggerEvent('[data-test-pgp-file]:nth-child(1) [data-test-pgp-file-input]', ...fileEvent());
    await waitUntil(() => {
      return !!this.spy.calledOnce;
    });
    const expected = [btoa(data), '', ''];
    assert.deepEqual(
      this.spy.getCall(0).args[0],
      expected,
      'calls onchange with an array of base64 converted files'
    );

    // add a file to the third input
    triggerEvent('[data-test-pgp-file]:nth-child(3) [data-test-pgp-file-input]', ...fileEvent());
    await waitUntil(() => {
      return !!this.spy.calledTwice;
    });
    const expected2 = [btoa(data), '', btoa(data)];
    assert.deepEqual(
      this.spy.getCall(1).args[0],
      expected2,
      'calls onchange with an array 2 base64 converted files'
    );

    // this will add a couple of inputs but should keep the existing 3
    this.set('listLength', 5);
    const expected3 = [btoa(data), '', btoa(data), '', ''];
    assert.deepEqual(
      this.spy.getCall(2).args[0],
      expected3,
      'lengthening the list with an array with one base64 converted files'
    );
  });

  test('it should render correct amount of file components on listLength change', async function (assert) {
    assert.expect(4);

    this.set('listLength', null);

    await render(hbs`
      <PgpList
        @listLength={{this.listLength}}
      />
    `);
    [1, 5, 3, 0].forEach((count) => {
      this.set('listLength', count);
      if (count) {
        assert
          .dom('[data-test-pgp-file]')
          .exists({ count }, `Correct number of file inputs render when listLength is updated to ${count}`);
      } else {
        assert.dom('[data-test-empty-text]').exists('Placeholder renders when list length is zero');
      }
    });
  });
});

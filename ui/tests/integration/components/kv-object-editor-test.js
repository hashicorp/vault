import { moduleForComponent, test } from 'ember-qunit';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';

import { create } from 'ember-cli-page-object';
import kvObjectEditor from '../../pages/components/kv-object-editor';

import sinon from 'sinon';
const component = create(kvObjectEditor);

moduleForComponent('kv-object-editor', 'Integration | Component | kv object editor', {
  integration: true,
  beforeEach() {
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

test('it renders with no initial value', function(assert) {
  let spy = sinon.spy();
  this.set('onChange', spy);
  this.render(hbs`{{kv-object-editor onChange=onChange}}`);
  assert.equal(component.rows().count, 1, 'renders a single row');
  component.addRow();
  return wait().then(() => {
    assert.equal(component.rows().count, 1, 'will only render row with a blank key');
  });
});

test('it calls onChange when the val changes', function(assert) {
  let spy = sinon.spy();
  this.set('onChange', spy);
  this.render(hbs`{{kv-object-editor onChange=onChange}}`);
  component.rows(0).kvKey('foo').kvVal('bar');
  wait().then(() => {
    assert.equal(spy.callCount, 2, 'calls onChange each time change is triggered');
    assert.deepEqual(
      spy.lastCall.args[0],
      { foo: 'bar' },
      'calls onChange with the JSON respresentation of the data'
    );
  });
  component.addRow();
  return wait().then(() => {
    assert.equal(component.rows().count, 2, 'adds a row when there is no blank one');
  });
});

test('it renders passed data', function(assert) {
  let metadata = { foo: 'bar', baz: 'bop' };
  this.set('value', metadata);
  this.render(hbs`{{kv-object-editor value=value}}`);
  assert.equal(
    component.rows().count,
    Object.keys(metadata).length + 1,
    'renders both rows of the metadata, plus an empty one'
  );
});

test('it deletes a row', function(assert) {
  let spy = sinon.spy();
  this.set('onChange', spy);
  this.render(hbs`{{kv-object-editor onChange=onChange}}`);
  component.rows(0).kvKey('foo').kvVal('bar');
  component.addRow();
  wait().then(() => {
    assert.equal(component.rows().count, 2);
    assert.equal(spy.callCount, 2, 'calls onChange for editing');
    component.rows(0).deleteRow();
  });

  return wait().then(() => {
    assert.equal(component.rows().count, 1, 'only the blank row left');
    assert.equal(spy.callCount, 3, 'calls onChange deleting row');
    assert.deepEqual(spy.lastCall.args[0], {}, 'last call to onChange is an empty object');
  });
});

test('it shows a warning if there are duplicate keys', function(assert) {
  let metadata = { foo: 'bar', baz: 'bop' };
  this.set('value', metadata);
  this.render(hbs`{{kv-object-editor value=value}}`);
  component.rows(0).kvKey('foo');

  return wait().then(() => {
    assert.ok(component.showsDuplicateError, 'duplicate keys are allowed but an error message is shown');
  });
});

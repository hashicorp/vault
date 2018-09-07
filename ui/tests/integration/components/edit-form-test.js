import { run } from '@ember/runloop';
import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import Service from '@ember/service';
import { moduleForComponent, test } from 'ember-qunit';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/edit-form';

const component = create(editForm);

const flash = Service.extend({
  success: sinon.stub(),
});

const createModel = (canDelete = true) => {
  return EmberObject.create({
    fields: [{ name: 'one', type: 'string' }, { name: 'two', type: 'boolean' }],
    canDelete,
    destroyRecord() {
      return resolve();
    },
    save() {
      return resolve();
    },
    rollbackAttributes() {},
  });
};

moduleForComponent('edit-form', 'Integration | Component | edit form', {
  integration: true,
  beforeEach() {
    this.register('service:flash-messages', flash);
    this.inject.service('flashMessages');
  },
});

test('it renders', function(assert) {
  let model = createModel();
  this.set('model', model);
  this.render(hbs`{{edit-form model=model}}`);

  assert.ok(component.fields.length, 2);
  assert.ok(component.showsDelete);
  assert.equal(component.deleteText, 'Delete');
});

test('it renders: custom deleteButton', function(assert) {
  let model = createModel();
  let delText = 'Exterminate';
  this.set('model', model);
  this.set('deleteButtonText', delText);
  this.render(hbs`{{edit-form model=model deleteButtonText=deleteButtonText}}`);

  assert.ok(component.showsDelete);
  assert.equal(component.deleteText, delText);
});

test('it calls flash message fns on save', function(assert) {
  let model = createModel();
  let onSave = () => {
    return resolve();
  };
  this.set('model', model);
  this.set('onSave', onSave);
  let saveSpy = sinon.spy(this, 'onSave');

  this.render(hbs`{{edit-form model=model onSave=onSave}}`);

  component.submit();
  return wait().then(() => {
    assert.ok(saveSpy.calledOnce, 'calls passed onSave');
    assert.equal(saveSpy.getCall(0).args[0].saveType, 'save');
    assert.deepEqual(saveSpy.getCall(0).args[0].model, model, 'passes model to onSave');
    assert.equal(this.flashMessages.success.callCount, 1, 'calls flash message success');
  });
});

test('it calls flash message fns on delete', function(assert) {
  let model = createModel();
  let onSave = () => {
    return resolve();
  };
  this.set('model', model);
  this.set('onSave', onSave);
  let saveSpy = sinon.spy(this, 'onSave');

  this.render(hbs`{{edit-form model=model onSave=onSave}}`);
  component.deleteButton();
  wait().then(() => {
    run(() => component.deleteConfirm());
  });

  return wait().then(() => {
    return wait().then(() => {
      assert.ok(saveSpy.calledOnce, 'calls onSave');
      assert.equal(saveSpy.getCall(0).args[0].saveType, 'destroyRecord');
      assert.deepEqual(saveSpy.getCall(0).args[0].model, model, 'passes model to onSave');
    });
  });
});

import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import itemDetails from 'vault/tests/pages/components/identity/item-details';

const component = create(itemDetails);

module('Integration | Component | identity/item details', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it renders the disabled warning', async function(assert) {
    let model = EmberObject.create({
      save() {
        return resolve();
      },
      disabled: true,
      canEdit: true,
    });
    sinon.spy(model, 'save');
    this.set('model', model);
    await render(hbs`{{identity/item-details model=model}}`);
    assert.dom('[data-test-disabled-warning]').exists();
    await component.enable();

    assert.ok(model.save.calledOnce, 'clicking enable calls model save');
  });

  test('it does not render the button if canEdit is false', async function(assert) {
    let model = EmberObject.create({
      disabled: true,
    });

    this.set('model', model);
    await render(hbs`{{identity/item-details model=model}}`);
    assert.dom('[data-test-disabled-warning]').exists('shows the warning banner');
    assert.dom('[data-test-enable]').doesNotExist('does not show the enable button');
  });

  test('it does not render the banner when item is enabled', async function(assert) {
    let model = EmberObject.create();
    this.set('model', model);

    await render(hbs`{{identity/item-details model=model}}`);
    assert.dom('[data-test-disabled-warning]').doesNotExist('does not show the warning banner');
  });
});

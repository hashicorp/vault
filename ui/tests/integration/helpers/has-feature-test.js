import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import waitForError from 'vault/tests/helpers/wait-for-error';

const versionStub = Service.extend({
  features: null,
});

module('helper:has-feature', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.register('service:version', versionStub);
    this.versionService = this.owner.lookup('service:version');
  });

  test('it asserts on unknown features', async function(assert) {
    let promise = waitForError();
    render(hbs`{{has-feature 'New Feature'}}`);
    let err = await promise;
    assert.ok(
      err.message.includes('New Feature is not one of the available values for Vault Enterprise features.'),
      'asserts when an unknown feature is passed as an arg'
    );
  });

  test('it is true with existing features', async function(assert) {
    this.set('versionService.features', ['HSM']);
    await render(hbs`{{if (has-feature 'HSM') 'It works' null}}`);
    assert.dom(this.element).hasText('It works', 'present features evaluate to true');
  });

  test('it is false with missing features', async function(assert) {
    this.set('versionService.features', ['MFA']);
    await render(hbs`{{if (has-feature 'HSM') 'It works' null}}`);
    assert.dom(this.element).hasText('', 'missing features evaluate to false');
  });
});

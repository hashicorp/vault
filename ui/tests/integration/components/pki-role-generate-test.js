import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import Sinon from 'sinon';

module('Integration | Component | pki-role-generate', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.model = this.store.createRecord('pki/certificate/generate', { name: 'role-name' });
    this.onSuccess = Sinon.spy();
  });

  test('it renders a form by default', async function (assert) {
    await render(hbs`<PkiRoleGenerate @model={{this.model}} @onSuccess={{this.onSuccess}} />`);

    assert
      .dom('[data-test-pki-generate-cert-form]')
      .exists({ count: 1 }, 'PKI Generate Certificate form exists');

    // Template block usage:
    await render(hbs`
      <PkiRoleGenerate>
        template block text
      </PkiRoleGenerate>
    `);

    assert.dom(this.element).hasText('template block text');
  });
});

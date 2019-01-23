import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import featuresSelection from 'vault/tests/pages/components/wizard/features-selection';
import hbs from 'htmlbars-inline-precompile';

const component = create(featuresSelection);

module('Integration | Component | features-selection', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it disables and enables wizard items according to user permissions', async function(assert) {
    const enabled = { Secrets: true, Authentication: false, Policies: true, Tools: false };
    await render(
      hbs`{{wizard/features-selection hasSecretsPermission=true hasAuthenticationPermission=false hasPoliciesPermission=true}}`
    );

    component.wizardItems.forEach(i => {
      assert.equal(
        i.hasDisabledTooltip,
        !enabled[i.text],
        'shows a tooltip only when the wizard item is not enabled'
      );
    });
  });

  test('it disables the start button if no wizard items are checked', async function(assert) {
    await render(hbs`{{wizard/features-selection}}`);
    assert.equal(component.hasDisabledStartButton, true);
  });

  test('it enables the start button when user has permission and wizard items are checked', async function(assert) {
    await render(hbs`{{wizard/features-selection hasSecretsPermission=true}}`);
    await component.selectSecrets();

    assert.equal(component.hasDisabledStartButton, false);
  });
});

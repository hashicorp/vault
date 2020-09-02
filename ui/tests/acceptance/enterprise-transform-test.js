import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';

import authPage from 'vault/tests/pages/auth';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import transformationsPage from 'vault/tests/pages/secrets/backend/transform/transformations';
import searchSelect from 'vault/tests/pages/components/search-select';

const searchSelectComponent = create(searchSelect);

const mount = async () => {
  let path = `transform-${Date.now()}`;
  await mountSecrets.enable('transform', path);
  return path;
};

module('Acceptance | Enterprise | Transform secrets', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it enables Transform secrets engine and shows tabs', async function(assert) {
    let path = `transform-${Date.now()}`;
    await mountSecrets.enable('transform', path);
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/list`,
      'mounts and redirects to the transformations list page'
    );
    assert.ok(transformationsPage.isEmpty, 'renders empty state');
    assert
      .dom('.is-active[data-test-tab="Transformations"]')
      .exists('Has Transformations tab which is active');
    assert.dom('[data-test-tab="Roles"]').exists('Has Roles tab');
    assert.dom('[data-test-tab="Templates"]').exists('Has Templates tab');
    assert.dom('[data-test-tab="Alphabets"]').exists('Has Alphabets tab');
  });

  test('it can create a transformation', async function(assert) {
    let path = await mount();
    await transformationsPage.createLink();
    assert.equal(currentURL(), `/vault/secrets/${path}/create`, 'redirects to create transformation page');
    await transformationsPage.name('foo');

    assert.dom('[data-test-input="type"').hasValue('fpe', 'Has type FPE by default');
    assert.dom('[data-test-input="tweak_source"]').exists('Shows tweak source when FPE');
    await transformationsPage.type('masking');
    assert
      .dom('[data-test-input="masking_character"]')
      .exists('Shows masking character input when changed to masking type');
    assert.dom('[data-test-input="tweak_source"]').doesNotExist('Does not show tweak source when masking');

    await clickTrigger('#template');
    assert.equal(searchSelectComponent.options.length, 2, 'list shows two builtin options by default');
    await selectChoose('#template', '.ember-power-select-option', 0);
    await transformationsPage.submit();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/show/foo`,
      'redirects to show transformation page after submit'
    );
  });
});

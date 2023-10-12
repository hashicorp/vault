import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | resultant-acl-banner', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders correctly by default', async function (assert) {
    await render(hbs`<ResultantAclBanner />`);

    assert.dom('[data-test-resultant-acl-banner] .hds-alert__title').hasText('Resultant ACL check failed');
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__description')
      .hasText(
        "Links might be shown that you don't have access to. Contact your administrator to update your policy."
      );
    assert.dom('[data-test-resultant-acl-reauthenticate]').doesNotExist('Does not show reauth link');
  });

  test('it renders correctly for namespaces', async function (assert) {
    const nsService = this.owner.lookup('service:namespace');
    nsService.setNamespace('my-ns');
    await render(hbs`<ResultantAclBanner @isEnterprise={{true}} />`);

    assert.dom('[data-test-resultant-acl-banner] .hds-alert__title').hasText('Resultant ACL check failed');
    assert
      .dom('[data-test-resultant-acl-banner] .hds-alert__description')
      .hasText("You may be in the wrong namespace, so links might be shown that you don't have access to.");
    assert.dom('[data-test-resultant-acl-reauthenticate]').exists('Shows reauth link');
    // unset
    nsService.setNamespace();
  });
});

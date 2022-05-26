import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';

module('Acceptance | mfa-login-enforcement', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'mfaConfig';
  });
  hooks.beforeEach(function () {
    return authPage.login();
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should create login enforcement', async function (assert) {
    await visit('/ui/vault/access');
    await click('[data-test-link="mfa"]');
    await click('[data-test-tab="enforcements"]');
    await click('[data-test-enforcement-create]');

    assert.dom('[data-test-mleh-title]').hasText('New enforcement', 'Title renders');
    await click('[data-test-mlef-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .exists({ count: 3 }, 'Validation error messages are displayed');

    await click('[data-test-mlef-cancel]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.index',
      'Cancel transitions to enforcements list'
    );
    await click('[data-test-enforcement-create]');
    await click('.breadcrumb a');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.index',
      'Breadcrumb transitions to enforcements list'
    );
    await click('[data-test-enforcement-create]');

    await fillIn('[data-test-mlef-input="name"]', 'foo');
    await click('[data-test-component="search-select"] .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await fillIn('[data-test-mount-accessor-select]', 'auth_userpass_bb95c2b1');
    await click('[data-test-mlef-add-target]');
    await click('[data-test-mlef-save]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.enforcement.index',
      'Route transitions to enforcement on save success'
    );
  });

  test('it should list login enforcements', async function (assert) {
    await visit('/vault/access/mfa/enforcements');
    assert.dom('[data-test-tab="enforcements"]').hasClass('active', 'Enforcements tab is active');
    assert.dom('.toolbar-link').exists({ count: 1 }, 'Correct number of toolbar links render');
    assert
      .dom('[data-test-enforcement-create]')
      .includesText('New enforcement', 'New enforcement link renders');

    await click('[data-test-enforcement-create]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.create',
      'New enforcement link transitions to create route'
    );
    await click('[data-test-mlef-cancel]');

    const enforcements = this.server.db.mfaLoginEnforcements.where({});
    const item = enforcements[0];
    assert.dom('[data-test-list-item]').exists({ count: enforcements.length }, 'Enforcements list renders');
    assert
      .dom(`[data-test-list-item="${item.name}"] svg`)
      .hasClass('flight-icon-lock', 'Lock icon renders for list item');
    assert.dom(`[data-test-list-item-title="${item.name}"]`).hasText(item.name, 'Enforcement name renders');

    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-list-item-link="details"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.enforcement.index',
      'Details more menu action transitions to enforcement route'
    );
    await click('.breadcrumb a');
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-list-item-link="edit"]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.enforcement.edit',
      'Edit more menu action transitions to enforcement edit route'
    );
  });

  test('it should display login enforcement', async function (assert) {
    await visit('/vault/access/mfa/enforcements');
    const enforcement = this.server.db.mfaLoginEnforcements.where({})[0];
    await click(`[data-test-list-item="${enforcement.name}"]`);
    // heading
    assert.dom('h1').includesText(enforcement.name, 'Name renders in title');
    assert.dom('h1 svg').hasClass('flight-icon-lock', 'Lock icon renders in title');
    assert.dom('[data-test-tab="targets"]').hasClass('active', 'Targets tab is active by default');
    assert.dom('[data-test-target]').exists({ count: 4 }, 'Targets render in list');
    // targets tab
    const targets = {
      accessor: ['userpass/', 'auth_userpass_bb95c2b1', '/ui/vault/access/userpass'],
      method: ['userpass', 'All userpass mounts (1)'],
      entity: [
        'test-entity',
        'f831667b-7392-7a1c-c0fc-33d48cb1c57d',
        '/ui/vault/access/identity/entities/f831667b-7392-7a1c-c0fc-33d48cb1c57d/details',
      ],
      group: [
        'test-group',
        '34db6b52-591e-bc22-8af0-4add5e167326',
        '/ui/vault/access/identity/groups/34db6b52-591e-bc22-8af0-4add5e167326/details',
      ],
    };
    for (const key in targets) {
      const t = targets[key];
      const selector = `[data-test-target="${t[0]}"]`;
      assert.dom(selector).includesText(`${t[0]} ${t[1]}`, `Target text renders for ${key} type`);
      if (key !== 'method') {
        await click(`${selector} [data-test-popup-menu-trigger]`);
        assert
          .dom(`[data-test-target-link="${t[0]}"]`)
          .hasAttribute('href', t[2], `Details link renders for ${key} type`);
      } else {
        assert.dom(`${selector} [data-test-popup-menu-trigger]`).doesNotExist('Method type has no link');
      }
    }
    // methods tab
    await click('[data-test-tab="methods"]');
    assert.dom('[data-test-tab="methods"]').hasClass('active', 'Methods tab is active');
    const method = this.owner.lookup('service:store').peekRecord('mfa-method', enforcement.mfa_method_ids[0]);
    assert
      .dom(`[data-test-mfa-method-list-item="${method.id}"]`)
      .includesText(
        `${method.name} ${method.id} Namespace: ${method.namespace_id}`,
        'Method list item renders'
      );
    await click('[data-test-popup-menu-trigger]');
    assert
      .dom(`[data-test-mfa-method-menu-link="details"]`)
      .hasAttribute('href', `/ui/vault/access/mfa/methods/${method.id}`, `Details link renders for method`);
    assert
      .dom(`[data-test-mfa-method-menu-link="edit"]`)
      .hasAttribute('href', `/ui/vault/access/mfa/methods/${method.id}/edit`, `Edit link renders for method`);
    // toolbar
    assert
      .dom('[data-test-enforcement-edit]')
      .hasAttribute(
        'href',
        `/ui/vault/access/mfa/enforcements/${enforcement.name}/edit`,
        'Toolbar edit action has link to edit route'
      );
    await click('[data-test-enforcement-delete]');
    assert.dom('[data-test-confirm-button]').isDisabled('Delete button disabled with no confirmation');
    await fillIn('[data-test-confirmation-modal-input]', enforcement.name);
    await click('[data-test-confirm-button]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.index',
      'Route transitions to enforcements list on delete success'
    );
  });

  test('it should edit login enforcement', async function (assert) {
    await visit('/vault/access/mfa/enforcements');
    const enforcement = this.server.db.mfaLoginEnforcements.where({})[0];
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-list-item-link="edit"]');

    assert.dom('h1').hasText('Update enforcement', 'Title renders');
    assert.dom('[data-test-mlef-input="name"]').hasValue(enforcement.name, 'Name input is populated');
    assert.dom('[data-test-mlef-input="name"]').isDisabled('Name is disabled and cannot be changed');

    const method = this.owner.lookup('service:store').peekRecord('mfa-method', enforcement.mfa_method_ids[0]);
    assert
      .dom('[data-test-selected-option]')
      .includesText(`${method.name} ${method.id}`, 'Selected mfa method renders');
    assert
      .dom('[data-test-mlef-target="Authentication mount"]')
      .includesText('Authentication mount auth_userpass_bb95c2b1', 'Accessor target populates');
    assert
      .dom('[data-test-mlef-target="Authentication method"]')
      .includesText('Authentication method userpass', 'Method target populates');
    assert
      .dom('[data-test-mlef-target="Group"]')
      .includesText('Group test-group 34db6b52-591e-bc22-8af0-4add5e167326', 'Group target populates');
    assert
      .dom('[data-test-mlef-target="Entity"]')
      .includesText('Entity test-entity f831667b-7392-7a1c-c0fc-33d48cb1c57d', 'Entity target populates');

    await click('[data-test-mlef-remove-target="Entity"]');
    await click('[data-test-mlef-remove-target="Group"]');
    await click('[data-test-mlef-remove-target="Authentication method"]');
    await click('[data-test-mlef-save]');

    assert.equal(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.enforcement.index',
      'Route transitions to enforcement on save success'
    );
    assert.dom('[data-test-target]').exists({ count: 1 }, 'Targets were successfully removed on save');
  });
});

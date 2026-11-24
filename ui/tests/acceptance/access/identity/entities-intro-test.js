/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

const ROUTE = '/vault/access/identity/entities';

module('Acceptance | identity | entities | intro page', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    return login();
  });

  test('Mode 1 (Full-Screen Card): shows intro card in initial empty state, dismisses on Skip, and reveals header CTA', async function (assert) {
    await visit(ROUTE);

    // 1. Verify full-screen intro card is rendered initially
    assert.dom(GENERAL.wizardIntro).exists('full-screen wizard intro card is shown on initial visit');
    assert.dom(GENERAL.button('Skip')).exists('"Skip" button is present in full-screen intro card');
    assert
      .dom(GENERAL.button('create-entity-link'))
      .doesNotExist('"Create new entity"  button is not present');
    assert.dom(GENERAL.button('intro')).doesNotExist('"New to Entities" intro button is not present');
    assert.dom(GENERAL.button('entity-merge-link')).doesNotExist('"Merge Entities"  button is not present');

    // 2. Dismiss via Skip
    await click(GENERAL.button('Skip'));

    // 3. Wizard intro is hidden, empty state is displayed, and page header shows "New to Entities?" and Create new entity buttons
    assert.dom(GENERAL.wizardIntro).doesNotExist('wizard intro is hidden after dismissal');
    assert.dom(GENERAL.emptyStateTitle).hasText('No entities yet', 'empty state title is rendered');
    assert.dom(GENERAL.button('intro')).exists('"New to Entities?" button is present in page header');
    assert.dom(GENERAL.button('entity-create-link')).exists('"Create new entity" button is in page header');
  });

  test('Mode 2 (Modal View): re-opens intro as modal via header CTA and closes via Close button', async function (assert) {
    // Start with wizard already dismissed
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.identityEntities);
    await visit(ROUTE);

    assert.dom(GENERAL.wizardIntro).doesNotExist('intro is initially dismissed');

    // 1. Click header CTA to open modal
    await click(GENERAL.button('intro'));

    // 2. Verify modal presentation
    assert.dom(GENERAL.modal.container('Entities')).exists('modal dialog is rendered');
    assert.dom(GENERAL.wizardIntro).exists('wizard intro is visible inside modal');
    assert.dom(GENERAL.button('Close')).exists('"Close" button is present in modal footer');

    // 3. Close the modal
    await click(GENERAL.button('Close'));
    assert.dom(GENERAL.modal.container('Entities')).doesNotExist('modal dialog is closed');
  });

  test('it does not show intro card or header CTA when entities exist', async function (assert) {
    const entityId = 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee';

    this.server.get('/identity/entity/id', () => ({
      data: {
        key_info: { [entityId]: { name: 'test-entity' } },
        keys: [entityId],
      },
    }));
    this.server.get(`/identity/entity/id/${entityId}`, () => ({
      data: { id: entityId, name: 'test-entity', aliases: [] },
    }));
    await visit(ROUTE);

    assert.dom(GENERAL.wizardIntro).doesNotExist('wizard intro is not shown when entities exist');
    assert
      .dom(GENERAL.button('intro'))
      .doesNotExist('"New to Entities?" button is not shown when entities exist');
    assert
      .dom(GENERAL.button('entity-merge-link'))
      .exists('"Merge Entities" button is present when not in empty state');

    assert
      .dom(GENERAL.button('entity-create-link'))
      .exists('"Create new entity" button is present when not in empty state');
  });
});

/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, findAll, settled, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

const SECRET_ENGINE_MANAGE_DROPDOWN_ROUTING_CASES = [
  {
    key: 'alicloud',
    type: 'alicloud',
    isEnginePathClickable: false,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'azure',
    type: 'azure',
    isEnginePathClickable: true,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'gcp',
    type: 'gcp',
    isEnginePathClickable: true,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'gcpkms',
    type: 'gcpkms',
    isEnginePathClickable: false,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'keymgmt',
    type: 'keymgmt',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'kubernetes',
    type: 'kubernetes',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
    expectedActionConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.kubernetes.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.kubernetes.configure',
    ],
    expectedLandingConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.kubernetes.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.kubernetes.configure',
    ],
  },
  {
    key: 'kvv1',
    type: 'kv',
    version: 1,
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'kvv2',
    type: 'kv',
    version: 2,
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: true,
    showConfigure: true,
    showDelete: true,
    expectedLandingConfigureRoutesOverride: ['vault.cluster.secrets.backend.kv.configuration'],
  },
  {
    key: 'transform',
    type: 'transform',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'transit',
    type: 'transit',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'kmip',
    type: 'kmip',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
    expectedActionConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.kmip.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.kmip.configure',
    ],
    expectedLandingConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.kmip.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.kmip.configure',
    ],
  },
  {
    key: 'ldap',
    type: 'ldap',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
    expectedActionConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.ldap.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.ldap.configure',
    ],
    expectedLandingConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.ldap.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.ldap.configure',
    ],
  },
  {
    key: 'pki',
    type: 'pki',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
    expectedActionConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.pki.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.pki.configuration.create',
    ],
    expectedLandingConfigureRoutesOverride: [
      // if the engine is configured
      'vault.cluster.secrets.backend.pki.configuration',
      // if the engine is not configured
      'vault.cluster.secrets.backend.pki.configuration.create',
    ],
  },
  {
    key: 'ssh',
    type: 'ssh',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'totp',
    type: 'totp',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'aws',
    type: 'aws',
    isEnginePathClickable: true,
    showManageDropdown: true,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'consul',
    type: 'consul',
    isEnginePathClickable: false,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'nomad',
    type: 'nomad',
    isEnginePathClickable: false,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'rabbitmq',
    type: 'rabbitmq',
    isEnginePathClickable: false,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
  },
  {
    key: 'database',
    type: 'database',
    isEnginePathClickable: true,
    showManageDropdown: false,
    showGeneratePolicy: false,
    showConfigure: true,
    showDelete: true,
    expectedLandingRouteOverride: 'vault.cluster.secrets.backend.overview',
    expectedLandingConfigureRoutesOverride: [],
  },
];

const secretsEngineListRoute = '/vault/secrets-engines';

const mountEngine = async ({ type, version }, path) => {
  await mountSecrets.visit();
  await click(GENERAL.cardContainer(type));
  await fillIn(GENERAL.inputByAttr('path'), path);
  if (type === 'kv' && version === 1) {
    await click(GENERAL.button('Method Options'));
    await mountSecrets.version(1);
  }
  await click(GENERAL.submitButton);
};

const filterEngineRowByPath = async (path) => {
  await visit(secretsEngineListRoute);
  const searchInputSelector = GENERAL.inputSearch('secret-engine-path');
  if (findAll(searchInputSelector).length) {
    await fillIn(searchInputSelector, path);
  }
};

const clickVisibleMenuItem = async (name) => {
  const visibleItem = findAll(GENERAL.menuItem(name)).find((el) => el.offsetParent !== null);
  if (!visibleItem) {
    throw new Error(`No visible menu item found for: ${name}`);
  }
  await click(visibleItem);
};

const assertMenuOptionVisibility = (assert, visibilityByOption, contextLabel, engineKey) => {
  for (const [option, isVisible] of Object.entries(visibilityByOption)) {
    if (isVisible) {
      assert.dom(GENERAL.menuItem(option)).exists(`${contextLabel} shows ${option} for ${engineKey}`);
    } else {
      assert
        .dom(GENERAL.menuItem(option))
        .doesNotExist(`${contextLabel} does not show ${option} for ${engineKey}`);
    }
  }
};

const clickVisibleConfirmButton = async () => {
  const visibleConfirmButton = findAll(GENERAL.confirmButton).find((el) => el.offsetParent !== null);
  if (!visibleConfirmButton) {
    return false;
  }
  await click(visibleConfirmButton);
  return true;
};

const expectedActionConfigureRoutes = (engineType) => {
  const { isConfigurable, configRoute } = engineDisplayData(engineType);
  if (!isConfigurable) {
    return ['vault.cluster.secrets.backend.configuration.general-settings'];
  }

  if (configRoute) {
    return [`vault.cluster.secrets.backend.${configRoute}`];
  }

  return [
    // if the engine is configured
    'vault.cluster.secrets.backend.configuration.plugin-settings',
    // if the engine is not configured
    'vault.cluster.secrets.backend.configuration.edit',
  ];
};

const expectedLandingRoute = ({ type, version = 1 }) => {
  const engineData = engineDisplayData(type);
  const isKvV1 = type === 'kv' && version === 1;

  if (engineData.isOnlyMountable) {
    return 'vault.cluster.secrets.backend.configuration.general-settings';
  }
  if (engineData.engineRoute && !isKvV1) {
    return `vault.cluster.secrets.backend.${engineData.engineRoute}`;
  }
  return 'vault.cluster.secrets.backend.list-root';
};

const expectedLandingConfigureRoutes = ({ type, version = 1 }) => {
  const engineData = engineDisplayData(type);
  const isKvV1 = type === 'kv' && version === 1;

  if (engineData.engineRoute && !isKvV1) {
    if (engineData.configRoute) {
      return [`vault.cluster.secrets.backend.${engineData.configRoute}`];
    }
  }

  if (engineData.isConfigurable) {
    return [
      // if the engine is configured
      'vault.cluster.secrets.backend.configuration.plugin-settings',
      // if the engine is not configured
      'vault.cluster.secrets.backend.configuration.edit',
    ];
  }

  return ['vault.cluster.secrets.backend.configuration.general-settings'];
};

const runEngineCase = async (assert, engine, uid, isEnterprise = false) => {
  const mountPath = `manage-${engine.key}-${uid}`;
  const actionConfigureRoutes =
    engine.expectedActionConfigureRoutesOverride || expectedActionConfigureRoutes(engine.type);
  const expectedManage = {
    showManageDropdown: engine.showManageDropdown ?? false,
    showGeneratePolicy: (engine.showGeneratePolicy ?? false) && isEnterprise,
    showConfigure: engine.showConfigure ?? true,
    showDelete: engine.showDelete ?? true,
  };

  // if engine path already exists, delete it before starting the test
  await runCmd(`delete sys/mounts/${mountPath}`);

  // mount the engine
  await mountEngine(engine, mountPath);

  // verify the engine shows in the list
  await filterEngineRowByPath(mountPath);
  assert.dom(GENERAL.tableRow()).exists(`row renders for ${engine.key}`);

  assert.dom(GENERAL.menuTrigger).exists(`Action menu is shown for ${engine.key}`);
  await click(GENERAL.menuTrigger);

  assertMenuOptionVisibility(
    assert,
    {
      Configure: expectedManage.showConfigure,
      Delete: expectedManage.showDelete,
    },
    'Action menu',
    engine.key
  );

  if (expectedManage.showConfigure) {
    // click configure and verify route
    await clickVisibleMenuItem('Configure');
    await settled();
    assert.true(
      actionConfigureRoutes.includes(currentRouteName()),
      `Action: Configure routes correctly for ${engine.key}`
    );
    await filterEngineRowByPath(mountPath);
  }

  if (expectedManage.showDelete) {
    // click delete and verify the engine is removed from the list
    await filterEngineRowByPath(mountPath);
    await click(GENERAL.menuTrigger);
    await clickVisibleMenuItem('Delete');
    const didConfirmActionDelete = await clickVisibleConfirmButton();
    assert.true(didConfirmActionDelete, `Action: Delete shows confirm button for ${engine.key}`);
    await settled();

    await filterEngineRowByPath(mountPath);
    assert.dom(GENERAL.tableRow()).doesNotExist(`Action: Delete removes ${engine.key} mount`);

    // remount the engine for manage dropdown testing
    await mountEngine(engine, mountPath);
    await filterEngineRowByPath(mountPath);
  }

  const isEnginePathClickable = engine.isEnginePathClickable ?? false;
  const backendLinkSelector = `a[href*="/vault/secrets-engines/${mountPath}"]`;

  if (!isEnginePathClickable) {
    // if the engine path is not expected to be clickable, verify it's not a link and skip the rest of the test
    assert.dom(backendLinkSelector).doesNotExist(`EnginePath is not a clickable link for ${engine.key}`);
    await runCmd(`delete sys/mounts/${mountPath}`);
    return;
  }

  assert.dom(backendLinkSelector).exists(`EnginePath is a clickable link for ${engine.key}`);
  await click(backendLinkSelector);

  const routeAfterPathClick = engine.expectedLandingRouteOverride || expectedLandingRoute(engine);
  assert.strictEqual(
    currentRouteName(),
    routeAfterPathClick,
    `Engine path click redirects to ${routeAfterPathClick} for ${engine.key}`
  );

  const shouldShowManageDropdown = expectedManage.showManageDropdown;

  if (!shouldShowManageDropdown) {
    // if manage dropdown is not expected to show on the landing page, verify it's not shown and skip the rest of the test
    assert
      .dom(GENERAL.dropdownToggle('Manage'))
      .doesNotExist(`Manage dropdown is not shown on landing page for ${engine.key}`);
    await runCmd(`delete sys/mounts/${mountPath}`);
    return;
  }

  assert
    .dom(GENERAL.dropdownToggle('Manage'))
    .exists(`Manage dropdown shows on landing page for ${engine.key}`);
  await click(GENERAL.dropdownToggle('Manage'));
  assertMenuOptionVisibility(
    assert,
    {
      'Generate policy': expectedManage.showGeneratePolicy,
      Configure: expectedManage.showConfigure,
      Delete: expectedManage.showDelete,
    },
    'Manage dropdown',
    engine.key
  );

  if (expectedManage.showConfigure) {
    // click configure and verify route
    await clickVisibleMenuItem('Configure');
    await settled();
    const allowedConfigureRoutes =
      engine.expectedLandingConfigureRoutesOverride || expectedLandingConfigureRoutes(engine);
    assert.true(
      allowedConfigureRoutes.includes(currentRouteName()),
      `Manage Configure routes correctly for ${engine.key}`
    );

    await filterEngineRowByPath(mountPath);
    await click(backendLinkSelector);
    await click(GENERAL.dropdownToggle('Manage'));
  }

  if (expectedManage.showDelete) {
    // click delete and verify the engine is removed from the list
    await clickVisibleMenuItem('Delete');
    const didConfirmManageDelete = await clickVisibleConfirmButton();
    if (!didConfirmManageDelete) {
      assert.true(
        engine.type === 'kubernetes',
        `Manage Delete missing confirm is only expected for kubernetes; got ${engine.key}`
      );
      await runCmd(`delete sys/mounts/${mountPath}`);
      await filterEngineRowByPath(mountPath);
      assert.dom(GENERAL.tableRow()).doesNotExist(`Manage Delete removes ${engine.key} mount`);
      return;
    }
    await settled();

    await filterEngineRowByPath(mountPath);
    assert.dom(GENERAL.tableRow()).doesNotExist(`Manage Delete removes ${engine.key} mount`);
    return; // if the delete action is confirmed, the engine should be removed and we can end the test here without needing to clean up again
  }
};

module('Acceptance | secrets-engines/manage-dropdown routing', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  for (const engine of SECRET_ENGINE_MANAGE_DROPDOWN_ROUTING_CASES) {
    const isEnterpriseOnly = !!engineDisplayData(engine.type).requiresEnterprise;
    const engineLabel = isEnterpriseOnly ? `${engine.key} (enterprise only)` : engine.key;

    test(`manage dropdown coverage | ${engineLabel}`, async function (assert) {
      const isEnterprise = this.owner.lookup('service:version').isEnterprise;
      await runEngineCase(assert, engine, this.uid, isEnterprise);
    });
  }
});

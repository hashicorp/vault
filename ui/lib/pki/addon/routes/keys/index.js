/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';
import { PkiListKeysListEnum } from '@hashicorp/vault-client-typescript';
import { paginate } from 'core/utils/paginate-list';

@withConfig()
export default class PkiKeysIndexRoute extends Route {
  @service secretMountPath;
  @service api;
  @service capabilities;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async fetchCapabilities(keyId) {
    const { pathFor } = this.capabilities;
    const backend = this.secretMountPath.currentPath;
    const pathMap = {
      import: pathFor('pkiKeysImport', { backend }),
      generate: pathFor('pkiKeysImport', { backend }),
      key: pathFor('pkiKey', { backend, keyId }),
    };
    const perms = await this.capabilities.fetch(Object.values(pathMap));

    return {
      canImportKey: perms[pathMap.import].canUpdate,
      canGenerateKey: perms[pathMap.generate].canUpdate,
      canRead: perms[pathMap.key].canRead,
      canEdit: perms[pathMap.key].canUpdate,
    };
  }

  async model(params) {
    const page = Number(params.page) || 1;
    const model = {
      hasConfig: this.pkiMountHasConfig,
      parentModel: this.modelFor('keys'),
    };

    try {
      const response = await this.api.secrets.pkiListKeys(
        this.secretMountPath.currentPath,
        PkiListKeysListEnum.TRUE
      );
      const keys = this.api.keyInfoToArray(response, 'key_id');
      const capabilities = await this.fetchCapabilities(keys[0].key_id);
      Object.assign(model, { ...capabilities, keys: paginate(keys, { page }) });
    } catch (e) {
      if (e.response.status === 404) {
        model.keys = [];
      } else {
        throw e;
      }
    }

    return model;
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: resolvedModel.parentModel.id },
      { label: 'Keys', route: 'keys.index', model: resolvedModel.parentModel.id },
    ];
    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}

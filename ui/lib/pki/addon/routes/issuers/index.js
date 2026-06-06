/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';
import { SecretsApiPkiListIssuersListEnum } from '@hashicorp/vault-client-typescript';
import { verifyCertificates, parseCertificate } from 'vault/utils/parse-pki-cert';

export default class PkiIssuersListRoute extends Route {
  @service secretMountPath;
  @service api;

  async getIssuerMetadata(issuer_id, listResponse) {
    try {
      const issuer = await this.api.secrets.pkiReadIssuer(issuer_id, this.secretMountPath.currentPath);
      const keyInfo = listResponse.key_info[issuer_id];
      const isRoot = await verifyCertificates(issuer.certificate, issuer.certificate);
      const parsedCertificate = parseCertificate(issuer.certificate);
      Object.assign(keyInfo, { ...keyInfo, ...issuer, isRoot, parsedCertificate });
    } catch (e) {
      return { ...listResponse.key_info[issuer_id], issuer_id };
    }
  }

  async model(params) {
    const page = Number(params.page) || 1;
    const parentModel = this.modelFor('issuers');

    try {
      const listResponse = await this.api.secrets.pkiListIssuers(
        this.secretMountPath.currentPath,
        SecretsApiPkiListIssuersListEnum.TRUE
      );
      // fetch full issuer data only if there are less than 10 issuers to avoid making too many requests
      if (listResponse.keys.length <= 10) {
        await Promise.all(
          listResponse.keys.map((issuer_id) => this.getIssuerMetadata(issuer_id, listResponse))
        );
      }
      const issuers = this.api.keyInfoToArray(listResponse);
      return {
        issuers: paginate(issuers, { page }),
        parentModel,
      };
    } catch (error) {
      if (error.response.status === 404) {
        return { parentModel };
      } else {
        throw error;
      }
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Issuers', route: 'issuers.index', model: currentPath },
    ];
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}

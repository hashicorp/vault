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

  // Returns a new enriched copy of the key_info entry for the given issuer_id.
  // Falls back to the bare list-stub on error so a single failing read doesn't break the page.
  async getIssuerMetadata(issuer_id, keyInfoEntry) {
    try {
      const issuer = await this.api.secrets.pkiReadIssuer(issuer_id, this.secretMountPath.currentPath);
      const isRoot = await verifyCertificates(issuer.certificate, issuer.certificate);
      const parsedCertificate = parseCertificate(issuer.certificate);
      return { ...keyInfoEntry, ...issuer, isRoot, parsedCertificate };
    } catch (e) {
      return { ...keyInfoEntry };
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

      // fetch full issuer data only if there are 10 or fewer issuers to avoid making too many requests
      let keyInfo = listResponse.key_info;
      if (listResponse.keys.length <= 10) {
        const enrichedEntries = await Promise.all(
          listResponse.keys.map((issuer_id) =>
            this.getIssuerMetadata(issuer_id, listResponse.key_info[issuer_id])
          )
        );
        keyInfo = Object.fromEntries(listResponse.keys.map((id, i) => [id, enrichedEntries[i]]));
      }

      const issuers = this.api.keyInfoToArray({ ...listResponse, key_info: keyInfo }, 'issuer_id');
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

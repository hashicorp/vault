/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { parseCertificate } from 'vault/utils/parse-pki-cert';

export default class PkiIssuerRotateRootRoute extends Route {
  @service secretMountPath;
  @service store;

  model() {
    const oldRoot = this.modelFor('issuers.issuer');
    const certData = parseCertificate(oldRoot.certificate);
    let parsingErrors;
    if (certData.parsing_errors && certData.parsing_errors.length > 0) {
      const errorMessage = certData.parsing_errors.map((e) => e.message).join(', ');
      parsingErrors = errorMessage;
    }
    return hash({
      oldRoot,
      certData,
      parsingErrors,
      backend: this.secretMountPath.currentPath,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: resolvedModel.oldRoot.backend },
      { label: 'Issuers', route: 'issuers.index', model: resolvedModel.oldRoot.backend },
      {
        label: resolvedModel.oldRoot.id,
        route: 'issuers.issuer.details',
        models: [resolvedModel.oldRoot.backend, resolvedModel.oldRoot.id],
      },
      { label: 'Rotate Root' },
    ];
  }
}

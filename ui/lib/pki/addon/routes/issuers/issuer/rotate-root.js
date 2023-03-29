/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiIssuerRoute from '../issuer';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import camelizeKeys from 'vault/utils/camelize-object-keys';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave('model.newRootModel')
export default class PkiIssuerRotateRootRoute extends PkiIssuerRoute {
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
    const newRootModel = this.store.createRecord('pki/action', {
      actionType: 'rotate-root',
      type: 'internal',
      ...camelizeKeys(certData), // copy old root settings over to new one
    });
    return hash({
      oldRoot,
      newRootModel,
      parsingErrors,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: resolvedModel.oldRoot.id, route: 'issuers.issuer.details' },
      { label: 'rotate root' },
    ];
  }
}

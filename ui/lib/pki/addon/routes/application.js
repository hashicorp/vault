/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiRoute extends Route {
  @service pathHelp;
  @service secretMountPath;

  beforeModel() {
    // We call pathHelp for all the models in this engine that use OpenAPI before any model hooks
    // so that the model attributes hydrate correctly. These only need to be called once to add
    // the openAPI attributes to the model prototype
    const mountPath = this.secretMountPath.currentPath;
    return hash({
      role: this.pathHelp.getNewModel('pki/role', mountPath),
      urls: this.pathHelp.getNewModel('pki/urls', mountPath),
      key: this.pathHelp.getNewModel('pki/key', mountPath),
      signCsr: this.pathHelp.getNewModel('pki/sign-intermediate', mountPath),
      certGenerate: this.pathHelp.getNewModel('pki/certificate/generate', mountPath),
      certSign: this.pathHelp.getNewModel('pki/certificate/sign', mountPath),
    });
  }
}

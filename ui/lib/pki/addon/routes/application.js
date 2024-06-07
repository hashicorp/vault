/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
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
      acme: this.pathHelp.getNewModel('pki/config/acme', mountPath),
      certGenerate: this.pathHelp.getNewModel('pki/certificate/generate', mountPath),
      certSign: this.pathHelp.getNewModel('pki/certificate/sign', mountPath),
      cluster: this.pathHelp.getNewModel('pki/config/cluster', mountPath),
      key: this.pathHelp.getNewModel('pki/key', mountPath),
      role: this.pathHelp.getNewModel('pki/role', mountPath),
      signCsr: this.pathHelp.getNewModel('pki/sign-intermediate', mountPath),
      tidy: this.pathHelp.getNewModel('pki/tidy', mountPath),
      urls: this.pathHelp.getNewModel('pki/config/urls', mountPath),
    });
  }
}

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
      acme: this.pathHelp.hydrateModel('pki/config/acme', mountPath),
      certGenerate: this.pathHelp.hydrateModel('pki/certificate/generate', mountPath),
      certSign: this.pathHelp.hydrateModel('pki/certificate/sign', mountPath),
      cluster: this.pathHelp.hydrateModel('pki/config/cluster', mountPath),
      key: this.pathHelp.hydrateModel('pki/key', mountPath),
      role: this.pathHelp.hydrateModel('pki/role', mountPath),
      signCsr: this.pathHelp.hydrateModel('pki/sign-intermediate', mountPath),
      tidy: this.pathHelp.hydrateModel('pki/tidy', mountPath),
      urls: this.pathHelp.hydrateModel('pki/config/urls', mountPath),
    });
  }
}

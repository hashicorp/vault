/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { duration } from 'core/helpers/format-duration';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';

/**
 * @module SecretsEngineMountConfig
 * SecretsEngineMountConfig component is used to display a "Show mount configuration" toggle section. It is generally used alongside the fetch-secret-engine-config decorator which displays the engine configuration above this component. Mount configuration is always available for display but is hidden by default behind a toggle.
 *
 * @example
 * <SecretsEngineMountConfig @secretsEngine={{this.model}} />
 *
 * @param {SecretsEngineResource} secretsEngine - the secrets engine resource containing the mount configuration details
 */

interface Args {
  model: SecretsEngineResource;
}
interface Field {
  label: string;
  value: string | boolean;
}

export default class SecretsEngineMountConfig extends Component<Args> {
  @tracked showConfig = false;
  @service declare readonly api: ApiService;

  submitForm = task(
    waitFor(async (event: Event) => {
      event?.preventDefault();
      try {
        // TODO we need to change how we send in the second parameter.
        await this.api.sys.mountsTuneConfigurationParameters(this.args.model.path, {
          description: 'blahlkj',
        });
      } catch (e) {
        debugger;
        // do something with the error
      }
    })
  );

  get fields(): Array<Field> {
    const { model } = this.args;
    return [
      { label: 'Secret engine type', value: model.engineType },
      { label: 'Path', value: model.path },
      { label: 'Accessor', value: model.accessor },
      { label: 'Local', value: model.local },
      { label: 'Seal wrap', value: model.sealWrap },
      { label: 'Default Lease TTL', value: duration([model.config.defaultLeaseTtl]) },
      { label: 'Max Lease TTL', value: duration([model.config.maxLeaseTtl]) },
      { label: 'Identity token key', value: model.config.identityTokenKey },
    ];
  }

  @action toggleSealWrap() {
    this.args.model.sealWrap = !this.args.model.sealWrap;
  }
}

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { duration } from 'core/helpers/format-duration';

import type SecretsEngineResource from 'vault/resources/secrets/engine';

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
  secretsEngine: SecretsEngineResource;
}
interface Field {
  label: string;
  value: string | boolean;
}

export default class SecretsEngineMountConfig extends Component<Args> {
  @tracked showConfig = false;

  get fields(): Array<Field> {
    const { secretsEngine } = this.args;
    return [
      { label: 'Secret engine type', value: secretsEngine.engineType },
      { label: 'Path', value: secretsEngine.path },
      { label: 'Accessor', value: secretsEngine.accessor },
      { label: 'Local', value: secretsEngine.local },
      { label: 'Seal wrap', value: secretsEngine.sealWrap },
      { label: 'Default Lease TTL', value: duration([secretsEngine.config.defaultLeaseTtl]) },
      { label: 'Max Lease TTL', value: duration([secretsEngine.config.maxLeaseTtl]) },
      { label: 'Identity token key', value: secretsEngine.config.identityTokenKey },
    ];
  }
}

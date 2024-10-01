/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { duration } from 'core/helpers/format-duration';

import type SecretEngineModel from 'vault/models/secret-engine';

/**
 * @module SecretsEngineMountConfig
 * SecretsEngineMountConfig component is used to display a "Show mount configuration" toggle section. It is generally used alongside the fetch-secret-engine-config decorator which displays the engine configuration above this component. Mount configuration is always available for display but is hidden by default behind a toggle.
 *
 * @example
 * <SecretsEngineMountConfig @model={{model}} />
 *
 * @param {Model} model- The secret engines model, generated via the secret-engine model and a belongsTo relationship connecting to the mount-config model.
 */

interface Args {
  model: SecretEngineModel;
}
interface Field {
  label: string;
  value: string | boolean;
}

export default class SecretsEngineMountConfig extends Component<Args> {
  @tracked showConfig = false;

  get fields(): Array<Field> {
    const { model } = this.args;
    return [
      { label: 'Secret Engine Type', value: model.engineType },
      { label: 'Path', value: model.path },
      { label: 'Accessor', value: model.accessor },
      { label: 'Local', value: model.local },
      { label: 'Seal Wrap', value: model.sealWrap },
      { label: 'Default Lease TTL', value: duration([model.config.defaultLeaseTtl]) },
      { label: 'Max Lease TTL', value: duration([model.config.maxLeaseTtl]) },
      { label: 'Identity Token Key', value: model.config.identityTokenKey },
    ];
  }
}

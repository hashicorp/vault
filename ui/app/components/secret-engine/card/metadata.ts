/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import SecretsEngineResource from 'vault/resources/secrets/engine';
interface Args {
  model: SecretsEngineResource;
}

export default class Metadata extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }

  updatePath() {
    // This method can be used to update the path of the secrets engine.
  }

  updateAccessor() {
    // This method can be used to update the accessor of the secrets engine.
  }
}

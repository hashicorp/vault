/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import SecretsEngineResource from 'vault/resources/secrets/engine';
interface Args {
  model: SecretsEngineResource;
}

export default class Version extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }
}

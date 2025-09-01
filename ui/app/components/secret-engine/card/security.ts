/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import SecretsEngineResource from 'vault/resources/secrets/engine';
interface Args {
  model: SecretsEngineResource;
}

export default class Security extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }

  @action toggleSealWrap() {
    this.args.model.seal_wrap = !this.args.model.seal_wrap;
  }

  @action toggleLocal() {
    this.args.model.local = !this.args.model.local;
  }
}

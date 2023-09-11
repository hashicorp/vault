/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import type PkiActionModel from 'vault/vault/models/pki/action';

interface Args {
  model: PkiActionModel;
}

export default class PagePkiIssuerGenerateIntermediateComponent extends Component<Args> {
  @tracked title = 'Generate intermediate CSR';
}

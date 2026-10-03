/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  breadcrumbs: Breadcrumb[];
}

export default class PagePkiIssuerImportComponent extends Component<Args> {
  @tracked title = 'Import a CA';
}

/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module PkiPageHeader
 * The `PkiPageHeader` is used to display pki page headers.
 *
 * @example ```js
 * <PkiPageHeader @backend="exampleBackend" />
 * ```
 */

interface Args {
  backend: { id: string };
}

export default class PkiPageHeader extends Component<Args> {
  get breadcrumbs() {
    return [
      {
        label: 'Secrets',
        route: 'secrets',
        linkExternal: true,
      },
      {
        label: this.args?.backend?.id,
      },
    ];
  }
}

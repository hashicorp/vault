/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

interface VaultReportingNamespacePickerSignature {
  Args: {
    namespaces: string[];
    onNamespaceChange: (namespace: string) => void;
  };
}

export default class VaultReportingNamespacePicker extends Component<VaultReportingNamespacePickerSignature> {
  @tracked selectedNamespace: string = this.args.namespaces[0] || '';
  @tracked search = '';

  get filteredNamespaces() {
    if (!this.search) return this.args.namespaces;

    return this.args.namespaces.filter((namespace) =>
      namespace.toLowerCase().includes(this.search.toLowerCase())
    );
  }

  handleNamespaceSelection = (namespace: string, close?: () => void) => {
    this.selectedNamespace = namespace;
    this.args.onNamespaceChange(namespace);
    close?.();
  };

  handleSearchInput = (event: Event) => {
    const input = event.target as HTMLInputElement;
    this.search = input.value;
  };
}

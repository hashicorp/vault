/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { SecurityPolicy } from './step-1';
import type NamespaceService from 'vault/services/namespace';
import {
  generateApiSnippet,
  generateCliSnippet,
  generateTerraformSnippet,
} from 'core/utils/code-generators/namespace-snippets';

interface Args {
  wizardState: {
    codeSnippet: null | string;
    creationMethod: CreationMethod;
    namespacePaths: string[];
    securityPolicyChoice: SecurityPolicy;
  };
  updateWizardState: (key: string, value: unknown) => void;
}

export enum CreationMethod {
  TERRAFORM = 'Terraform automation',
  APICLI = 'API/CLI',
  UI = 'Vault UI workflow',
}

interface CreationMethodChoice {
  icon: string;
  label: CreationMethod;
  description: string;
  isRecommended?: boolean;
}

export default class WizardNamespacesStep3 extends Component<Args> {
  @service declare readonly namespace: NamespaceService;
  @tracked creationMethodChoice: CreationMethod;
  @tracked selectedTabIdx = 0;

  methods = CreationMethod;
  policy = SecurityPolicy;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.creationMethodChoice = this.args.wizardState.creationMethod || CreationMethod.TERRAFORM;
  }

  creationMethodOptions: CreationMethodChoice[] = [
    {
      icon: 'terraform-color',
      label: CreationMethod.TERRAFORM,
      description:
        'Manage configurations by Infrastructure as Code. This creation method improves resilience and ensures common compliance requirements.',
      isRecommended: true,
    },
    {
      icon: 'terminal-screen',
      label: CreationMethod.APICLI,
      description:
        'Manage namespaces directly via the Vault CLI or REST API. Best for quick updates, custom scripting, or terminal-based workflows.',
    },
    {
      icon: 'sidebar',
      label: CreationMethod.UI,
      description:
        'Apply changes immediately. Note: Changes made in the UI will be overwritten by any future updates made via Infrastructure as Code (Terraform).',
    },
  ];

  get tfSnippet() {
    const { namespacePaths } = this.args.wizardState;
    return generateTerraformSnippet(namespacePaths, this.namespace.path);
  }

  get customTabs() {
    const { namespacePaths } = this.args.wizardState;
    return [
      {
        key: 'api',
        label: 'API',
        snippet: generateApiSnippet(namespacePaths, this.namespace.path),
      },
      {
        key: 'cli',
        label: 'CLI',
        snippet: generateCliSnippet(namespacePaths, this.namespace.path),
      },
    ];
  }

  @action
  onChange(choice: CreationMethodChoice) {
    this.creationMethodChoice = choice.label;
    this.args.updateWizardState('creationMethod', choice.label);
    // Update the code snippet whenever the creation method changes
    this.updateCodeSnippet();
  }

  @action
  onTabChange(idx: number) {
    this.selectedTabIdx = idx;

    // Update the code snippet whenever the tab changes
    this.updateCodeSnippet();
  }

  // Update the wizard state with the current code snippet
  @action
  updateCodeSnippet() {
    if (this.creationMethodChoice === CreationMethod.TERRAFORM) {
      this.args.updateWizardState('codeSnippet', this.tfSnippet);
    } else if (this.creationMethodChoice === CreationMethod.APICLI) {
      const snippet = this.customTabs[this.selectedTabIdx]?.snippet;
      this.args.updateWizardState('codeSnippet', snippet);
    }
  }
}

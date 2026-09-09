/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { SecurityPolicy } from './step-1';
import type NamespaceService from 'vault/services/namespace';
import type SnippetService from 'vault/services/snippet';
import type AnalyticsService from 'vault/services/analytics';
import { CreationMethod } from 'vault/utils/constants/snippet';
import {
  WIZARD_NAMESPACE_STEP3_SELECT_TERRAFORM,
  WIZARD_NAMESPACE_STEP3_SELECT_CLI_API,
  WIZARD_NAMESPACE_STEP3_SELECT_UI,
} from 'vault/utils/analytic-events';
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

interface CreationMethodChoice {
  icon: string;
  label: CreationMethod;
  description: string;
  isRecommended?: boolean;
}

const CREATION_METHOD_EVENTS: Record<CreationMethod, string> = {
  [CreationMethod.TERRAFORM]: WIZARD_NAMESPACE_STEP3_SELECT_TERRAFORM,
  [CreationMethod.APICLI]: WIZARD_NAMESPACE_STEP3_SELECT_CLI_API,
  [CreationMethod.UI]: WIZARD_NAMESPACE_STEP3_SELECT_UI,
};

export default class WizardNamespacesStep3 extends Component<Args> {
  @service declare readonly namespace: NamespaceService;
  @service declare readonly snippet: SnippetService;
  @service declare readonly analytics: AnalyticsService;

  methods = CreationMethod;
  policy = SecurityPolicy;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    window.scrollTo({ top: 0, behavior: 'smooth' });
    this.snippet.reset(this.args.wizardState.creationMethod || CreationMethod.TERRAFORM);
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

  get creationMethodChoice() {
    return this.snippet.creationMethodChoice;
  }

  get selectedTabIdx() {
    return this.snippet.selectedTabIdx;
  }

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
    this.snippet.setCreationMethod(choice.label, this.tfSnippet, this.customTabs);
    this.args.updateWizardState('creationMethod', choice.label);
    this.args.updateWizardState('codeSnippet', this.snippet.codeSnippet);
    this.analytics.trackEvent(CREATION_METHOD_EVENTS[choice.label], {
      namespace: 'namespace-wizard',
      action: 'selected',
      elementId: 'method-selector',
      channel: 'webpage',
      location: 'step-3',
    });
  }

  @action
  onTabChange(idx: number) {
    this.snippet.setSelectedTab(idx, this.tfSnippet, this.customTabs);
    this.args.updateWizardState('codeSnippet', this.snippet.codeSnippet);
  }

  @action
  updateCodeSnippet() {
    this.snippet.persistSnippet(this.tfSnippet, this.customTabs);
    this.args.updateWizardState('codeSnippet', this.snippet.codeSnippet);
  }
}

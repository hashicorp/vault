/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
import { cliTemplate } from 'core/utils/code-generators/cli';

import type NamespaceService from 'vault/services/namespace';
import type { CliTemplateArgs } from 'core/utils/code-generators/cli';
import type { TerraformResourceTemplateArgs } from 'core/utils/code-generators/terraform';

interface SnippetOption {
  key: string;
  label: string;
  language?: 'bash' | 'go' | 'hcl' | 'json' | 'log' | 'ruby' | 'shell-session' | 'yaml';
  snippet: string;
}

interface Args {
  customTabs?: SnippetOption[];
  tfvpArgs?: TerraformResourceTemplateArgs;
  cliArgs?: CliTemplateArgs;
}

export default class CodeGeneratorAutomationSnippets extends Component<Args> {
  @service declare readonly namespace: NamespaceService;

  get tabs() {
    return this.args.customTabs || this.snippetTabs;
  }

  get snippetTabs() {
    return [
      {
        key: 'terraform',
        label: 'Terraform Vault Provider',
        snippet: terraformResourceTemplate(this.terraformOptions),
        language: 'hcl',
      },
      {
        key: 'cli',
        label: 'CLI',
        snippet: cliTemplate(this.args.cliArgs),
        language: 'shell',
      },
    ];
  }

  get terraformOptions() {
    const { tfvpArgs } = this.args || {};
    // only add namespace if we're not in root (when namespace is '')
    if (tfvpArgs && !this.namespace.inRootNamespace) {
      const { resourceArgs } = tfvpArgs;
      return { ...tfvpArgs, resourceArgs: { namespace: `"${this.namespace.path}"`, ...resourceArgs } };
    }
    return tfvpArgs;
  }
}

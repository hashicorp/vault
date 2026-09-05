/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
import { cliTemplate } from 'core/utils/code-generators/cli';
import { apiTemplate } from 'core/utils/code-generators/api';

import type NamespaceService from 'vault/services/namespace';
import type { CliTemplateArgs } from 'core/utils/code-generators/cli';
import type { TerraformResourceTemplateArgs } from 'core/utils/code-generators/terraform';
import type { ApiTemplateArgs } from 'core/utils/code-generators/api';
import type { HTMLElementEvent } from 'vault/forms';

interface SnippetOption {
  key: string;
  label: string;
  language?: 'bash' | 'go' | 'hcl' | 'json' | 'log' | 'ruby' | 'shell-session' | 'yaml';
  snippet: string;
}

interface Args {
  customTabs?: SnippetOption[];
  tfSnippet?: string;
  tfvpArgs?: TerraformResourceTemplateArgs;
  cliArgs?: CliTemplateArgs;
  apiArgs?: ApiTemplateArgs;
  onTabChange?: (tabIdx: number) => void;
}

export default class CodeGeneratorAutomationSnippets extends Component<Args> {
  @service declare readonly namespace: NamespaceService;

  get tabs() {
    return this.args.customTabs || this.snippetTabs;
  }

  get snippetTabs() {
    const tabs = [];
    const { tfSnippet, tfvpArgs, cliArgs, apiArgs } = this.args;
    if (tfSnippet || tfvpArgs) {
      tabs.push({
        key: 'terraform',
        label: 'Terraform Vault Provider',
        snippet: this.terraformSnippet,
        language: 'hcl',
      });
    }
    if (cliArgs) {
      tabs.push({
        key: 'cli',
        label: 'CLI',
        snippet: cliTemplate(cliArgs),
        language: 'shell',
      });
    }
    if (apiArgs) {
      tabs.push({
        key: 'api',
        label: 'API',
        snippet: apiTemplate(this.apiOptions),
      });
    }
    return tabs;
  }

  get terraformSnippet() {
    const { tfSnippet, tfvpArgs } = this.args;
    if (tfSnippet !== undefined) {
      if (!this.namespace.inRootNamespace) {
        const namespaceLine = `  namespace = "${this.namespace.path}"`;
        return tfSnippet.replace(/^(resource "[^"]*" "[^"]*" \{)/, `$1\n${namespaceLine}`);
      }
      return tfSnippet;
    }
    // Legacy support for tfvpArgs - used in the event tfSnippet is not provided.
    if (tfvpArgs && !this.namespace.inRootNamespace) {
      // only add namespace if we're not in root (when namespace is '')
      const { resourceArgs } = tfvpArgs;
      return terraformResourceTemplate({
        ...tfvpArgs,
        resourceArgs: { namespace: `"${this.namespace.path}"`, ...resourceArgs },
      });
    }
    return terraformResourceTemplate(tfvpArgs);
  }

  get apiOptions() {
    const { apiArgs } = this.args || {};
    // only add namespace if we're not in root (when namespace is '')
    if (apiArgs && !this.namespace.inRootNamespace) {
      return { ...apiArgs, namespace: this.namespace.path };
    }
    return apiArgs;
  }

  @action
  handleTabChange(_event: HTMLElementEvent<HTMLInputElement>, tabIndex: number) {
    const { onTabChange } = this.args;
    if (onTabChange) {
      onTabChange(tabIndex);
    }
  }
}

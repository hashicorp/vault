/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from 'tracked-built-ins';
import { CreationMethod } from 'vault/utils/constants/snippet';

/**
 * @module ExternalPkiImplementationSelect
 * The `ExternalPkiImplementationSelect` is used to display the external overview route
 */

export enum SetupSteps {
  ACME_CONFIG = 'acme-config',
  ROLE_CONFIG = 'role-config',
}

interface Args {
  engineId: string;
  steps?: SetupSteps[];
  title: string;
}

interface StepConfig {
  number: number;
  title: string;
  config: { terraform?: object; cli?: object; api?: object };
}

export default class ExternalPkiImplementationSelect extends Component<Args> {
  @tracked selectedMethod = CreationMethod.TERRAFORM;

  get displaySteps(): StepConfig[] {
    const stepsToShow = this.args.steps || [SetupSteps.ACME_CONFIG, SetupSteps.ROLE_CONFIG];
    const steps: StepConfig[] = [];
    let stepNumber = 1;

    if (stepsToShow.includes(SetupSteps.ACME_CONFIG)) {
      steps.push({
        number: stepNumber++,
        title: 'Configure an ACME account',
        config: this.acmeConfig,
      });
    }

    if (stepsToShow.includes(SetupSteps.ROLE_CONFIG)) {
      steps.push({
        number: stepNumber++,
        title: 'Create a role',
        config: this.roleConfig,
      });
    }

    return steps;
  }

  methodOptions = [
    {
      icon: 'terraform-color',
      label: CreationMethod.TERRAFORM,
      description: 'Manage configurations by Infrastructure as Code.',
    },
    {
      icon: 'terminal-screen',
      label: CreationMethod.APICLI,
      description: 'Configure directly via the Vault CLI or REST API.',
    },
  ];

  get acmeConfig() {
    const mountPath = this.args.engineId;
    const payload = {
      directory_url: '<directory_url>',
      email_contacts: '[<email_contacts>]',
    };
    const terraform = {
      resource: 'vault_pki_external_ca_secret_backend_acme_account',
      resourceArgs: {
        mount: `"${mountPath}"`,
        name: '<name>',
        ...payload,
      },
    };
    const cli = {
      command: `write ${mountPath}/config/acme-account/<name> \\\n`,
      content: ' directory_url="<directory_url>" \\\n  email_contacts="<email_contacts>" \\\n',
    };
    const api = {
      url: `${mountPath}/config/acme-account/<name>`,
      payload,
    };

    return this.selectedMethod === CreationMethod.TERRAFORM ? { terraform } : { cli, api };
  }

  get roleConfig() {
    const mountPath = this.args.engineId;
    const payload = {
      allowed_domains: '[<allowed_domains>]',
      allowed_domain_options: '[<allowed_domain_options>]',
      allowed_challenge_types: '[<allowed_challenge_types>]',
    };
    const terraform = {
      resource: 'vault_pki_external_ca_secret_backend_role',
      resourceArgs: {
        mount: `"${mountPath}"`,
        name: '<name>',
        allowed_domains: '[<allowed_domains>]',
        allowed_domain_options: '[<allowed_domain_options>]',
        allowed_challenge_types: '[<allowed_challenge_types>]',
      },
    };
    const cli = {
      command: `write ${mountPath}/role/<name> \\\n`,
      content:
        ' allowed_domains="<allowed_domains>" \\\n  allowed_domain_options="<allowed_domain_options>" \\\n  allowed_challenge_types="<allowed_challenge_types>" \\\n',
    };
    const api = { url: `${mountPath}/role/<name>`, payload };

    return this.selectedMethod === CreationMethod.TERRAFORM ? { terraform } : { cli, api };
  }

  @action
  onChange(choice: { icon: string; label: CreationMethod; description: string }) {
    this.selectedMethod = choice.label;
  }
}

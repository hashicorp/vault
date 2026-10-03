/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

export enum SecurityPolicy {
  FLEXIBLE = 'flexible',
  STRICT = 'strict',
}

interface Args {
  wizardState: {
    securityPolicyChoice: SecurityPolicy | null;
  };
}

export default class WizardNamespacesStep1 extends Component<Args> {
  policy = SecurityPolicy;

  get cardInfo() {
    const { wizardState } = this.args;
    if (wizardState.securityPolicyChoice === SecurityPolicy.FLEXIBLE) {
      return {
        title: 'Single namespace',
        description:
          'Your organization should be comfortable with your current setup of one global namespace. You can always add more namespaces later.',
        bestFor: [
          'Small teams or orgs just getting started with Vault.',
          'Centralized platform teams managing all secrets.',
        ],
        avoidIf: [
          'You need strong isolation between teams or business units.',
          'You plan to scale to 100+ applications or secrets engines.',
          'You anticipate needing per-team Terraform workflows.',
        ],
      };
    }
    return {
      title: 'Multiple namespaces',
      description:
        'Create isolation for clear ownership and scalability for strictly separated teams or applications.',
      diagram: '~/multi-namespace.gif',
      bestFor: [
        'Heavily regulated organizations with strict boundary enforcement between tenants.',
        'Organizations already confident with Terraform and namespace scoping.',
      ],
      avoidIf: ["You're not absolutely sure you need hard isolation and nesting."],
    };
  }
}

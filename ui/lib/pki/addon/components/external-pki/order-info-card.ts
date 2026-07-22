/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached } from '@glimmer/tracking';

import type {
  PkiExternalCaReadLookupCertResponse,
  PkiExternalCaReadLookupOrderResponse,
  PkiExternalCaReadRoleOrderFetchCertResponse,
  PkiExternalCaReadRoleOrderStatusResponse,
} from '@hashicorp/vault-client-typescript';
import type { ApiParsedError } from 'vault/vault/api';

interface Challenge {
  challenge_status: string;
  challenge_type: string;
  expires: string;
  requires_manual_fulfillment: 'true' | 'false'; // yes this is a boolean as a string, also not currently rendered
}

interface ChallengeRow {
  challenge_status: string;
  challenge_type: string;
  expires: string;
  requires_manual_fulfillment: string;
}

interface IdentifierRow {
  identifier: string;
  challenge_status: string;
  challenge_type: string;
  isOpen: boolean;
  children: ChallengeRow[];
}

interface Args {
  order: {
    details: PkiExternalCaReadLookupOrderResponse | PkiExternalCaReadRoleOrderStatusResponse | undefined;
    error?: ApiParsedError;
  };
  certificate: {
    details: PkiExternalCaReadLookupCertResponse | PkiExternalCaReadRoleOrderFetchCertResponse | undefined;
    error?: ApiParsedError;
  };
  orderId: string;
  engineId: string;
}

export default class ExternalPkiOrderInfoCardComponent extends Component<Args> {
  tableColumns = [
    { key: 'identifier', label: 'Identifier', isExpandable: true },
    { key: 'challenge_status', label: 'Status' },
    { key: 'challenge_type', label: 'Type' },
    { key: 'expires', label: 'Expires' },
  ];

  get hasCertificate() {
    return !!this.args.certificate.details;
  }

  @cached
  get configDisplay() {
    const { certificate, order, orderId } = this.args;
    // If we have a certificate then the order has completed. Only show limited order info to prioritize cert details.
    if (certificate.details) {
      const { order_status, last_update, role_name } = order.details || {};
      return { order_status, last_update, role_name, ...certificate.details };
    }
    // Otherwise show all of the order information.
    return { order_id: orderId, ...order.details };
  }

  @cached
  get tableData(): IdentifierRow[] {
    const challenges = this.args.order.details?.challenges as Record<string, Challenge[]> | undefined;
    if (!challenges) return [];
    return Object.entries(challenges).map(([identifier, challenges]) => {
      // At least one challenge has to be valid for the identifier's authorization to be valid
      const validChallenges = challenges.filter((c) => c.challenge_status === 'valid');
      return {
        identifier,
        challenge_status: validChallenges.length ? 'valid' : 'pending',
        challenge_type: validChallenges.map((c) => c.challenge_type.toUpperCase()).join(', '),
        isOpen: true,
        children: challenges.map((challenge) => ({
          challenge_status: challenge.challenge_status,
          challenge_type: challenge.challenge_type.toUpperCase(),
          expires: challenge.expires,
          requires_manual_fulfillment: challenge.requires_manual_fulfillment,
        })),
      };
    });
  }
}

/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached } from '@glimmer/tracking';
import { PkiExternalCaReadLookupOrderResponse } from '@hashicorp/vault-client-typescript';
import { parseAPITimestamp } from 'core/utils/date-formatters';

interface Challenge {
  challenge_status: string;
  challenge_type: string;
  expires: string;
  requires_manual_fulfillment: 'true' | 'false'; // yes this is a boolean as as a string, also not currently rendered
}

interface ChallengeRow {
  challenge_status: string;
  challenge_type: string;
  expires: string;
  requires_manual_fulfillment: string;
}

interface IdentifierRow {
  isOpen: boolean;
  children: ChallengeRow[];
}

interface Args {
  order: PkiExternalCaReadLookupOrderResponse;
}

export default class ExternalPkiOrderInfoCardComponent extends Component<Args> {
  tableColumns = [
    {
      key: 'identifier',
      label: 'Identifier',
      isExpandable: true,
    },
    {
      key: 'challenge_status',
      label: 'Status',
    },
    {
      key: 'challenge_type',
      label: 'Type',
    },
    {
      key: 'expires',
      label: 'Expires',
    },
  ];

  @cached
  get tableData(): IdentifierRow[] {
    const { order } = this.args;
    if (!order?.challenges) return [];
    return Object.entries(order.challenges as Record<string, Challenge[]>).map(([identifier, challenges]) => {
      // At least one challenge has to be valid for the identifier's authorization to be valid
      const validChallenges = challenges.filter((c) => c.challenge_status === 'valid');
      return {
        identifier,
        challenge_status: validChallenges.length ? 'valid' : 'pending',
        challenge_type: validChallenges.map((c) => c.challenge_type.toUpperCase()).join(', '),
        isOpen: true,
        children: challenges.map((challenge: Challenge) => ({
          challenge_status: challenge.challenge_status,
          challenge_type: challenge.challenge_type.toUpperCase(),
          expires: challenge.expires,
          requires_manual_fulfillment: challenge.requires_manual_fulfillment,
        })),
      };
    });
  }

  // TEMPLATE HELPERS
  formatDate = (isoString: string) => parseAPITimestamp(isoString, "MM/dd/yyyy, HH:mm 'UTC'");
}

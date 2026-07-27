/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached } from '@glimmer/tracking';
import { OrdersOrderRouteModel } from 'pki/routes/external/orders/order';

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
  order: OrdersOrderRouteModel['order'];
  certificate: OrdersOrderRouteModel['certificate'];
  orderId: string;
  engineId: string;
}

export default class ExternalPkiOrderCertDetailsComponent extends Component<Args> {
  tableColumns = [
    { key: 'identifier', label: 'Identifier', isExpandable: true },
    { key: 'challenge_status', label: 'Status' },
    { key: 'challenge_type', label: 'Type' },
    { key: 'expires', label: 'Expires' },
  ];

  get orderCompleted() {
    // Only true if user has permission to read the order status AND it's completed
    return this.args.order?.details?.order_status === 'completed';
  }

  @cached
  get orderConfigDisplay() {
    const { order, orderId: order_id } = this.args;
    if (!order.details) return null;

    const { details } = order;
    // Rename last_update so it's clear the time corresponds to the order, not certificate
    const { order_status, last_update: last_order_update, role_name } = details;
    switch (order_status) {
      case 'completed':
        // Completed orders: cert display is priority, so only show summary order fields here.
        return { order_status, last_order_update, role_name };
      case 'expired':
      case 'error':
        // Expired or failed orders: cert is not fetchable, show relevant fields so the user understands why.
        return {
          order_id,
          order_status,
          last_order_update,
          creation_date: details.creation_date,
          expires: details.expires,
          role_name,
          last_error: details.last_error,
        };
      default:
        // All remaining statuses (pending, processing, etc.): show everything.
        return { order_id, ...details };
    }
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

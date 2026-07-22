/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export type OrderStatusName =
  | 'undefined'
  | 'new'
  | 'submitted'
  | 'awaiting-challenge-fulfillment'
  | 'vault-challenge-fulfillment'
  | 'vault-challenge-propagating'
  | 'notify-acme-server-challenges-completed'
  | 'processing-challenge'
  | 'fetching-certificate'
  | 'completed'
  | 'revoked'
  | 'expired'
  | 'error';

// Challenge statuses originate from RFC 8555:
// https://datatracker.ietf.org/doc/html/rfc8555#section-7.1.6
export type ChallengeStatusName = 'valid' | 'pending' | 'processing' | 'invalid';

export default function mapStatusBadge(status: OrderStatusName | ChallengeStatusName) {
  switch (status) {
    case 'new':
    case 'submitted':
    case 'awaiting-challenge-fulfillment':
    case 'vault-challenge-fulfillment':
    case 'vault-challenge-propagating':
    case 'notify-acme-server-challenges-completed':
    case 'processing-challenge':
    case 'fetching-certificate':
    case 'pending':
      return { text: 'Pending', color: 'warning' };
    case 'processing':
      return { text: 'Processing', color: 'warning' };
    case 'error':
      return { text: 'Failed', color: 'critical' };
    case 'invalid':
      return { text: 'Invalid', color: 'critical' };
    case 'expired':
      return { text: 'Expired', color: 'neutral' };
    case 'revoked':
      return { text: 'Revoked', color: 'critical' };
    case 'completed':
      return { text: 'Issued', color: 'success' };
    case 'valid':
      return { text: 'Valid', color: 'success' };
    default:
      return { text: 'Unknown', color: 'neutral' };
  }
}

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

export default function mapOrderStatus(orderStatus: OrderStatusName) {
  switch (orderStatus) {
    case 'new':
    case 'submitted':
    case 'awaiting-challenge-fulfillment':
    case 'vault-challenge-fulfillment':
    case 'vault-challenge-propagating':
    case 'notify-acme-server-challenges-completed':
    case 'processing-challenge':
    case 'fetching-certificate':
      return { text: 'Pending', color: 'warning' };
    case 'error':
      return { text: 'Failed', color: 'critical' };
    case 'expired':
      return { text: 'Expired', color: 'neutral' };
    case 'revoked':
      return { text: 'Revoked', color: 'critical' };
    case 'completed':
      return { text: 'Issued', color: 'success' };
    default:
      return { text: 'Unknown', color: 'neutral' };
  }
}

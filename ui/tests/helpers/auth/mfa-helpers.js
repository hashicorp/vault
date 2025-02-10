/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const setupTotpMfaResponse = (mountPath) => ({
  auth: {
    mfa_requirement: {
      mfa_request_id: '0edf0945-da02-1300-9a0a-cb052cd94eb4',
      mfa_constraints: {
        [mountPath]: {
          any: [
            {
              type: 'totp',
              id: '7028db82-7de3-01d7-26b5-84b147c80966',
              uses_passcode: true,
            },
          ],
        },
      },
    },
    num_uses: 0,
  },
});

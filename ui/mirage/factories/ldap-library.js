/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  name: (i) => `library-${i}`,
  service_account_names: () => ['fizz@example.com', 'buzz@example.com'],
  ttl: '10h',
  max_ttl: '20h',
  disable_check_in_enforcement: false,
});

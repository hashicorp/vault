/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';
import { v4 as uuidv4 } from 'uuid';

export default Factory.extend({
  id: (i) => `${i + 1}`,
  display_name: (i) => `agent-${i + 1}`,
  entity_id: () => uuidv4(),
  description: (i) => `Description for agent ${i + 1}`,
  owner: (i) => `owner-${i + 1}@example.com`,
  ceiling_policy: () => ['default', 'agent-policy'],
  no_default_ceiling_policy: () => Math.random() < 0.5,
  creation_time: () => new Date().toISOString(),
  last_updated_time: () => new Date().toISOString(),
});

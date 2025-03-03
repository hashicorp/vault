/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MessagesRoute extends Route {
  @service pagination;

  model() {
    return [
      {
        id: '1',
        pluginName: 'AWS',
        author: 'cool',
        description: 'hello hello hello',
        externalUrl: 'www.hello.com',
        pluginVersion: '1.0.0',
        pluginType: 'SECRET',
        tags: 'community',
        official: { author: 'hashicorp', tags: 'official' },
      },
    ];
  }
}

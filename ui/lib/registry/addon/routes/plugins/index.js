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
        pluginName: 'Plugin 1',
        author: 'Some awesome author',
        description:
          "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
        externalUrl: 'www.hello.com',
        pluginVersion: '1.0.0',
        pluginType: 'SECRET',
        tags: 'community',
        official: { author: 'hashicorp', tags: 'official' },
        publishDate: new Date(),
      },
    ];
  }
}

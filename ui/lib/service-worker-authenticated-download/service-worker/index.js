/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { createUrlRegEx, urlMatchesAnyPattern } from 'ember-service-worker/service-worker/url-utils';

var patterns = ['/v1/sys/storage/raft/snapshot'];
var REGEXES = patterns.map(createUrlRegEx);

function sendMessage(message) {
  return self.clients.matchAll({ includeUncontrolled: true, type: 'window' }).then(function (results) {
    var client = results[0];
    return new Promise(function (resolve, reject) {
      var messageChannel = new MessageChannel();
      messageChannel.port2.onmessage = function (event) {
        if (event.data.error) {
          reject(event.data.error);
        } else {
          resolve(event.data.token);
        }
      };

      client.postMessage(message, [messageChannel.port1]);
    });
  });
}

function authenticateRequest(request) {
  // copy the reaquest headers so we can mutate them
  const headers = new Headers(request.headers);

  // get and set vault token so the request is authenticated
  return sendMessage({ action: 'getToken' }).then(function (token) {
    headers.set('X-Vault-Token', token);

    // continue the fetch with the new request
    // that has the auth header
    return fetch(
      new Request(request.url, {
        method: request.method,
        headers,
      })
    );
  });
}

self.addEventListener('fetch', function (fetchEvent) {
  const request = fetchEvent.request;

  if (urlMatchesAnyPattern(request.url, REGEXES) && request.method === 'GET') {
    return fetchEvent.respondWith(authenticateRequest(request));
  } else {
    return fetchEvent.respondWith(fetch(request));
  }
});

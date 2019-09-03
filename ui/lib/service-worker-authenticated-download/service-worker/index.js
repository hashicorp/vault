import { createUrlRegEx, urlMatchesAnyPattern } from 'ember-service-worker/service-worker/url-utils';

const patterns = ['/v1/sys/storage/raft/snapshot'];
const REGEXES = patterns.map(createUrlRegEx);

async function sendMessage(message) {
  const [client] = await self.clients.matchAll({ includeUncontrolled: true, type: 'window' });
  return new Promise(function(resolve, reject) {
    var messageChannel = new MessageChannel();
    messageChannel.port2.onmessage = function(event) {
      if (event.data.error) {
        reject(event.data.error);
      } else {
        resolve(event.data.token);
      }
    };

    client.postMessage(message, [messageChannel.port1]);
  });
}

const authenticateRequest = async request => {
  // copy the reaquest headers so we can mutate them
  let headers = new Headers(request.headers);

  // get and set vault token so the request is authenticated
  let token = await sendMessage({ action: 'getToken' });
  headers.set('X-Vault-Token', token);

  // continue the fetch with the new request
  // that has the auth header
  return fetch(
    new Request(request.url, {
      method: request.method,
      headers,
    })
  );
};

self.addEventListener('fetch', fetchEvent => {
  const request = fetchEvent.request;

  if (urlMatchesAnyPattern(request.url, REGEXES) && request.method === 'GET') {
    return fetchEvent.respondWith(authenticateRequest(request));
  } else {
    return fetchEvent.respondWith(fetch(request));
  }
});

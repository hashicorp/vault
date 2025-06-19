import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString  } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    discardResponseBodies: true,
    scenarios: {
%{for idx, host in hosts~}
        ${idx}: {
            executor: 'ramping-arrival-rate',
            exec: 'kv',
            timeUnit: '1s',
            preAllocatedVUs: 3000,
            maxVUs: 6000,
            stages: [
              { duration: '10s', target: 1000 },
              { duration: '60s', target: 1000 },
              { duration: '10s', target: 4000 },
              { duration: '120s', target: 4000 },
              { duration: '10s', target: 1000 },
              { duration: '120s', target: 1000 },
              { duration: '10s', target: 0 },
            ],
            env: {
                VAULT_ADDR: 'http://${host.private_ip}:8200',
            },
        },
%{endfor~}
    },
};

function linearJitterBackoff(attemptNumber, minWait, maxWait) {
  let rand = Math.random();
  let jitter = rand * (maxWait - minWait);
  let jitterMin = int(minWait + jitter);
  return jitterMin * attemptNumber;
}

function retryRequest(url, data, params, retries) {
  let res = http.post(url, data, params);
  let attemptNum = 1;
  while (res.status >= 400 && attemptNum <= retries) {
    sleep(linearJitterBackoff(attemptNum, 1, 5));
    attemptNum++;
    res = http.post(url, data, params);
  }
  return res;
}

export function setup() {
  let data = {
    "type": "kv",
    "config": {
      "default_lease_ttl": "0s",
      "max_lease_ttl": "0s",
      "force_no_cache": false,
    },
    "options": {
      "version": "1"
    },
    "local": false,
    "seal_wrap": false,
  }
  http.post('http://${leader_addr}:8200/v1/sys/mounts/kv', JSON.stringify(data), {
    headers: { 'Content-Type': 'application/json', 'X-Vault-Token': '${vault_token}' },
  });
}

export function kv() {
  const key = randomString(8);
  const url = `$${__ENV.VAULT_ADDR}/v1/kv/data/` + key;

  let data = {"data": {"foo": "bar"}};
  let params = {'headers': {'Content-Type': 'application/json', 'X-Vault-Token': '${vault_token}'}};
  params['tags'] = {'name': 'create-secret'};
  let res = retryRequest(url, JSON.stringify(data), params);
  check(res, { 'put was success': (r) => r.status < 400 });
}

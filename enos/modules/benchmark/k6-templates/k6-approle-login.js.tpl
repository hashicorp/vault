import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
import { randomString  } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    scenarios: {
%{for idx, host in hosts~}
        ${idx}: {
            executor: 'ramping-arrival-rate',
            exec: 'approle_login',
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
  // Unique auth method name for each instance
  const auth_name = "approle-" + randomString(8);

  // Mount approle auth method
  let mount_authmethod_data = {
    "type": "approle",
  }
  const auth_mount_url = 'http://${leader_addr}:8200/v1/sys/auth/' + auth_name;

  let params = {'headers': { 'Content-Type': 'application/json', 'X-Vault-Token': '${vault_token}' }};
  let res = http.post(auth_mount_url, JSON.stringify(mount_authmethod_data), params);
  if (res.status >= 400) {
    console.log("Failed to mount approle auth method")
    console.log(res);
    abort("Failed to mount approle auth method: ", auth_name);
  }

  // Create approle roles
  const auth_api_url = 'http://${leader_addr}:8200/v1/auth/' + auth_name;
  let create_approle_data = {
    "policies": "default",
    "secret_id_ttl": "0m",
    "token_ttl": "1m",
  }

  for (let i = 0; i < 10; i++) {
    const role_name = "approle" + i;
    let create_role_url =  auth_api_url + '/role/' + role_name;
    res = http.post(create_role_url, JSON.stringify(create_approle_data), params);
    if (res.status >= 400) {
      console.log("Failed to create approle role");
      console.log(res);
      return;
    }

    let role_id_payload = {
      "role_name": role_name,
      "role_id": role_name + "-role",
    }
    res = http.post(create_role_url + '/role-id', JSON.stringify(role_id_payload), params);
    if (res.status >= 400) {
      console.log("Failed to create role-id");
      console.log(res);
      return;
    }

    let secret_id_payload = {
      "role_name": role_name,
      "secret_id": role_name + "-secret",
    }
    res = http.post(create_role_url + '/custom-secret-id', JSON.stringify(secret_id_payload), params);
    if (res.status >= 400) {
      console.log("Failed to create secret-id");
      console.log(res);
      return;
    }

    let login_data = {
      "role_id": role_name + "-role",
      "secret_id": role_name + "-secret",
    }
    const login_url = 'http://${leader_addr}:8200/v1/auth/' + auth_name + '/login';
    res = http.post(login_url, JSON.stringify(login_data), params);
    if (res.status >= 400) {
      console.log("Failed to login to " + auth_name);
      console.log(res);
      return;
    }
  }

  return auth_name;
}

export function approle_login(data) {
  const auth_name = data;
  const role_id = Math.floor(Math.random() * 9);
  let login_data = {
    "role_id": "approle" + role_id + "-role",
    "secret_id": "approle" + role_id + "-secret",
  }

  const login_url = `$${__ENV.VAULT_ADDR}/v1/auth/` + auth_name + '/login';
  let params = {
    'headers': {
      'Content-Type': 'application/json',
      'X-Vault-Token': '${vault_token}'
    },
    timeout: '10s'
  };
  let res = retryRequest(login_url, JSON.stringify(login_data), params);
  if (res.status >= 400 && res.status != 503) {
    console.log("Failed to login");
    console.log(res);
  }
  check(res, { 'login success': (r) => r.status < 400 });
}

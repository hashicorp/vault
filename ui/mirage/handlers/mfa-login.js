/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import Ember from 'ember';

// initial auth response cache -- lookup by mfa_request_id key
const authResponses = {};
// mfa requirement cache -- lookup by mfa_request_id key
const mfaRequirement = {};

// may be imported in tests when the validation request needs to be intercepted to make assertions prior to returning a response
// in that case it may be helpful to still use this validation logic to ensure to payload is as expected
export const validationHandler = (schema, req) => {
  try {
    const { mfa_request_id, mfa_payload } = JSON.parse(req.requestBody);
    const mfaRequest = mfaRequirement[mfa_request_id];

    if (!mfaRequest) {
      return new Response(404, {}, { errors: ['MFA Request ID not found'] });
    }
    // validate request body
    for (const constraintId in mfa_payload) {
      // ensure ids were passed in map
      const method = mfaRequest.methods.find(({ id }) => id === constraintId);
      if (!method) {
        return new Response(400, {}, { errors: [`Invalid MFA constraint id ${constraintId} passed in map`] });
      }
      // test non-totp validation by rejecting all pingid requests
      if (method.type === 'pingid') {
        return new Response(403, {}, { errors: ['PingId MFA validation failed'] });
      }
      // validate totp passcode
      const passcode = mfa_payload[constraintId][0];
      if (method.uses_passcode) {
        const expectedPasscode = method.type === 'duo' ? 'passcode=test' : 'test';
        if (passcode !== expectedPasscode) {
          const error =
            {
              used: 'code already used; new code is available in 30 seconds',
              limit:
                'maximum TOTP validation attempts 4 exceeded the allowed attempts 3. Please try again in 15 seconds',
            }[passcode] || 'failed to validate';
          console.log(error); // eslint-disable-line
          return new Response(403, {}, { errors: [error] });
        }
      } else if (passcode) {
        // for okta and duo, reject if a passcode was provided
        return new Response(400, {}, { errors: ['Passcode should only be provided for TOTP MFA type'] });
      }
    }
    return authResponses[mfa_request_id];
  } catch (error) {
    console.log(error); // eslint-disable-line
    return new Response(500, {}, { errors: ['Mirage Handler Error: /sys/mfa/validate'] });
  }
};

export default function (server) {
  // generate different constraint scenarios and return mfa_requirement object
  const generateMfaRequirement = (req, res) => {
    const { user } = req.params;
    // uses_passcode automatically set to true in factory for totp type
    const m = (type, uses_passcode = false) => server.create('mfa-method', { type, uses_passcode });
    let mfa_constraints = {};
    let methods = []; // flat array of methods for easy lookup during validation

    function generator() {
      const methods = [];
      const constraintObj = [...arguments].reduce((obj, methodArray, index) => {
        obj[`test_${index}`] = { any: methodArray };
        methods.push(...methodArray);
        return obj;
      }, {});
      return [constraintObj, methods];
    }

    if (user === 'mfa-a') {
      [mfa_constraints, methods] = generator([m('totp')]); // 1 constraint 1 passcode
    } else if (user === 'mfa-b') {
      [mfa_constraints, methods] = generator([m('okta')]); // 1 constraint 1 non-passcode
    } else if (user === 'mfa-c') {
      [mfa_constraints, methods] = generator([m('totp'), m('duo', true)]); // 1 constraint 2 passcodes
    } else if (user === 'mfa-d') {
      [mfa_constraints, methods] = generator([m('okta'), m('duo')]); // 1 constraint 2 non-passcode
    } else if (user === 'mfa-e') {
      [mfa_constraints, methods] = generator([m('okta'), m('totp')]); // 1 constraint 1 passcode 1 non-passcode
    } else if (user === 'mfa-f') {
      [mfa_constraints, methods] = generator([m('totp')], [m('duo', true)]); // 2 constraints 1 passcode for each
    } else if (user === 'mfa-g') {
      [mfa_constraints, methods] = generator([m('okta')], [m('duo')]); // 2 constraints 1 non-passcode for each
    } else if (user === 'mfa-h') {
      [mfa_constraints, methods] = generator([m('totp')], [m('okta')]); // 2 constraints 1 passcode 1 non-passcode
    } else if (user === 'mfa-i') {
      [mfa_constraints, methods] = generator([m('okta'), m('totp')], [m('totp')]); // 2 constraints 1 passcode/1 non-passcode 1 non-passcode
    } else if (user === 'mfa-j') {
      [mfa_constraints, methods] = generator([m('pingid')]); // use to test push failures
    } else if (user === 'mfa-k') {
      [mfa_constraints, methods] = generator([m('duo', true)]); // test duo passcode and prepending passcode= to user input
    }
    const mfa_request_id = crypto.randomUUID();
    const mfa_requirement = {
      mfa_request_id,
      mfa_constraints,
    };
    // cache mfa requests to test different validation scenarios
    mfaRequirement[mfa_request_id] = { methods };
    // cache auth response to be returned later by sys/mfa/validate
    authResponses[mfa_request_id] = { ...res };
    return mfa_requirement;
  };
  // passthrough original request, cache response and return mfa stub
  const passthroughLogin = async (schema, req) => {
    // test totp not configured scenario
    if (req.params.user === 'totp-na') {
      return new Response(400, {}, { errors: ['TOTP mfa required but not configured'] });
    }
    const mock = req.params.user ? req.params.user.includes('mfa') : null;
    // bypass mfa for users that do not match type
    if (!mock) {
      req.passthrough();
    } else if (Ember.testing) {
      // use root token in test environment
      const res = await fetch('/v1/auth/token/lookup-self', { headers: { 'X-Vault-Token': 'root' } });
      if (res.status < 300) {
        const json = res.json();
        if (Ember.testing) {
          json.auth = {
            ...json.data,
            policies: [],
            metadata: { username: 'foobar' },
          };
          json.data = null;
        }
        return { auth: { mfa_requirement: generateMfaRequirement(req, json) } };
      }
      return new Response(500, {}, { errors: ['Mirage error fetching root token in testing'] });
    } else {
      const xhr = req.passthrough();
      xhr.onreadystatechange = () => {
        if (xhr.readyState === 4 && xhr.status < 300) {
          // XMLHttpRequest response prop only has a getter -- redefine as writable and set value
          Object.defineProperty(xhr, 'response', {
            writable: true,
            value: JSON.stringify({
              auth: { mfa_requirement: generateMfaRequirement(req, JSON.parse(xhr.responseText)) },
            }),
          });
        }
      };
    }
  };
  server.post('/auth/:method/login/:user', passthroughLogin);

  server.post('/sys/mfa/validate', validationHandler);
}

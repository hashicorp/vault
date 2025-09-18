/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* 
* HOW TO USE THIS HANDLER
- Create the mfa users below in Vault (e.g. "mfa-a") - there is a config.sh script in the vault-tools repo to expedite this
- do NOT configure MFA, this mirage handler stubs all of that
- For TOTP the passcode for successful login is "test"
- PingID is used to test push failure states
*/

import { Response } from 'miragejs';
import Ember from 'ember';

export const QR_CODE_URL =
  'otpauth://totp/vault-not-self-enroll:daf8420c-0b6b-34e6-ff38-ee1ed093bea9?algorithm=SHA1\u0026digits=6\u0026issuer=vault-not-self-enroll\u0026period=30\u0026secret=JGPHY3TZBIUCHWYN7ZO3LHISKQIAJZGL';

// initial auth response cache -- lookup by mfa_request_id key
const authResponses = {};
// mfa requirement cache -- lookup by mfa_request_id key
const mfaRequirement = {};
// mfa constraint cache -- lookup by mfa_request_id key
const mfaConstraints = {};

export const buildEnforcementError = (method, constraintName) => {
  // Although this block is just for pingid, adding failure message for posterity and testing other method error states
  const failure =
    method.type === 'totp'
      ? 'failed to validate TOTP passcode'
      : `${method.type} authentication failed: "Login request denied."`;
  const name = constraintName || 'My Secure Enforcement';
  const msg = `failed to satisfy enforcement ${name}. error: 2 errors occurred:\n\t* ${failure}\n\t* login MFA validation failed for methodID: [${method.id}]\n\n`;
  return new Response(403, {}, { errors: [msg] });
};

// may be imported in tests when the validation request needs to be intercepted to make assertions prior to returning a response
// in that case it may be helpful to still use this validation logic to ensure to payload is as expected
export const validationHandler = (schema, req) => {
  try {
    const { mfa_request_id, mfa_payload } = JSON.parse(req.requestBody);
    const mfaRequest = mfaRequirement[mfa_request_id];
    const constraintNameLookup = mfaConstraints[mfa_request_id];

    if (!mfaRequest) {
      return new Response(404, {}, { errors: ['MFA Request ID not found'] });
    }
    // validate request body
    for (const methodId in mfa_payload) {
      // ensure ids were passed in map
      const method = mfaRequest.methods.find(({ id }) => id === methodId);
      if (!method) {
        return new Response(400, {}, { errors: [`Invalid MFA method id ${methodId} passed in map`] });
      }
      // test non-totp validation by rejecting all pingid requests
      if (method.type === 'pingid') {
        return buildEnforcementError(method, constraintNameLookup[method.id]);
      }
      // validate totp passcode
      const passcode = mfa_payload[methodId][0];
      if (method.uses_passcode) {
        const expectedPasscode = method.type === 'duo' ? 'passcode=test' : 'test';
        if (passcode !== expectedPasscode) {
          const error =
            {
              used: 'code already used; new code is available in 30 seconds',
              limit:
                'maximum TOTP validation attempts 4 exceeded the allowed attempts 3. Please try again in 15 seconds',
            }[passcode] || null;
          console.log(error); // eslint-disable-line
          if (error) {
            return new Response(403, {}, { errors: [error] });
          }
          return buildEnforcementError(method, constraintNameLookup[method.id]);
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
    const m = (type, { uses_passcode = false, self_enrollment_enabled = false } = {}) =>
      server.create('mfa-method', { type, uses_passcode, self_enrollment_enabled });
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
      [mfa_constraints, methods] = generator([m('totp'), m('duo', { uses_passcode: true })]); // 1 constraint 2 passcodes
    } else if (user === 'mfa-d') {
      [mfa_constraints, methods] = generator([m('okta'), m('duo')]); // 1 constraint 2 non-passcode
    } else if (user === 'mfa-e') {
      [mfa_constraints, methods] = generator([m('okta'), m('totp')]); // 1 constraint 1 passcode 1 non-passcode
    } else if (user === 'mfa-f') {
      [mfa_constraints, methods] = generator([m('totp')], [m('duo', { uses_passcode: true })]); // 2 constraints 1 passcode for each
    } else if (user === 'mfa-g') {
      [mfa_constraints, methods] = generator([m('okta')], [m('duo')]); // 2 constraints 1 non-passcode for each
    } else if (user === 'mfa-h') {
      [mfa_constraints, methods] = generator([m('totp')], [m('okta')]); // 2 constraints 1 passcode 1 non-passcode
    } else if (user === 'mfa-i') {
      [mfa_constraints, methods] = generator([m('okta'), m('totp')], [m('duo', { uses_passcode: true })]); // 2 constraints 1 non-passcode or 1 non-passcode and 1 passcode
    } else if (user === 'mfa-j') {
      [mfa_constraints, methods] = generator([m('pingid')]); // use to test push failures
    } else if (user === 'mfa-k') {
      [mfa_constraints, methods] = generator([m('duo', true)]); // test duo passcode and prepending passcode= to user input

      // * SELF-ENROLLMENT USERS BELOW
      // users match counterpart config scenario above
      // e.g. "mfa-a" is the same as "mfa-a-self", but with self-enroll enabled
    } else if (user === 'mfa-a-self') {
      // 1 constraint 1 passcode
      [mfa_constraints, methods] = generator([m('totp', { self_enrollment_enabled: true })]);
    } else if (user === 'mfa-c-self') {
      // 1 constraint 2 passcodes
      [mfa_constraints, methods] = generator([
        m('totp', { self_enrollment_enabled: true }),
        m('duo', { uses_passcode: true }),
      ]);
    } else if (user === 'mfa-f-self') {
      // 2 constraints 1 passcode for each
      [mfa_constraints, methods] = generator(
        [m('totp', { self_enrollment_enabled: true })],
        [m('duo', { uses_passcode: true })]
      );
    } else if (user === 'mfa-h-self') {
      // 2 constraints 1 passcode 1 non-passcode
      [mfa_constraints, methods] = generator([m('totp', { self_enrollment_enabled: true })], [m('okta')]);
    } else if (user === 'mfa-i-self') {
      // 2 constraints 1 non-passcode or 1 non-passcode and 1 passcode
      [mfa_constraints, methods] = generator(
        [m('okta'), m('totp', { self_enrollment_enabled: true })],
        [m('duo', { uses_passcode: true })]
      );
    } else if (user === 'mfa-z-self') {
      // We've discussed that this scenario likely won't be allowed in the real-world
      // by having the API restrict self-enrollment so it's only be possible when only ONE
      // constraint has self_enrollment_enabled.
      // 3 constraints, two have 2 methods (and each includes a method with self_enrollment_enabled)
      [mfa_constraints, methods] = generator(
        [m('totp', { self_enrollment_enabled: true }), m('pingid')],
        [m('totp', { self_enrollment_enabled: true }), m('okta')],
        [m('duo')]
      );
    }

    const mfa_request_id = crypto.randomUUID();
    const mfa_requirement = {
      mfa_request_id,
      mfa_constraints,
    };
    // cache mfa requests to test different validation scenarios
    mfaRequirement[mfa_request_id] = { methods };
    // cache login enforcement names
    for (const [constraintName, constraint] of Object.entries(mfa_constraints)) {
      // Create lookup by method ID
      constraint.any.forEach((method) => {
        if (!mfaConstraints[mfa_request_id]) {
          mfaConstraints[mfa_request_id] = {};
        }
        mfaConstraints[mfa_request_id][method.id] = constraintName;
      });
    }

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

  server.post('/identity/mfa/method/totp/self-enroll', async () => {
    // For this endpoint to return a legitimate QR code, the user has to actually exist in Vault.
    // Since we're using mirage to stub MFA users and methods this just returns a dummy QR code for testing.
    return {
      data: {
        barcode:
          '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',
        url: QR_CODE_URL,
      },
    };
  });
}

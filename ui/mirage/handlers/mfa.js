import faker from 'faker';
import { Response } from 'miragejs';

export default function (server) {
  // initial auth response cache -- lookup by mfa_request_id key
  const authResponses = {};
  // mfa enforcement cache -- lookup by mfa_request_id key
  const mfaEnforcement = {};
  // passthrough original request, cache response and return mfa stub
  const passthroughLogin = (schema, req) => {
    const xhr = req.passthrough();
    xhr.onreadystatechange = () => {
      if (xhr.readyState === 4 && xhr.status < 300) {
        const type = ['duo', 'okta', 'pingid', 'totp'].includes(req.params.user) ? req.params.user : null;
        // bypass mfa for users that do not match type
        if (type) {
          const res = JSON.parse(xhr.responseText);
          const mfa_request_id = faker.datatype.uuid();
          // cache auth response to be returned later by sys/mfa/validate
          authResponses[mfa_request_id] = { ...res };
          // unsure of final response shape when mfa is enabled
          // it looks like the new object will be added under the auth key so return only that for now
          const mfa_enforcement = {
            mfa_request_id,
            mfa_constraints: [server.create('mfa-method', { type })],
          };
          // cache mfa requests to test different validation scenarios
          mfaEnforcement[mfa_request_id] = mfa_enforcement;
          // XMLHttpRequest response prop only has a getter -- redefine as writable and set value
          Object.defineProperty(xhr, 'response', {
            writable: true,
            value: JSON.stringify({ auth: { mfa_enforcement } }),
          });
        }
      }
    };
  };
  server.post('/auth/:method/login/:user', passthroughLogin);

  // unsure if the token method will utilize mfa
  // server.get('/auth/token/lookup-self', passthroughLogin);

  server.post(
    '/sys/mfa/validate',
    (schema, req) => {
      try {
        const { mfa_request_id, mfa_payload } = JSON.parse(req.requestBody);
        const mfaRequest = mfaEnforcement[mfa_request_id];

        if (!mfaRequest) {
          return new Response(404, {}, { errors: ['MFA Request ID not found'] });
        }
        // validate request body
        for (let constraintId in mfa_payload) {
          // ensure ids were passed in map
          const mfaConstraint = mfaRequest.mfa_constraints.find(({ id }) => id === constraintId);
          if (!mfaConstraint) {
            return new Response(
              400,
              {},
              { errors: [`Invalid MFA constraint id ${constraintId} passed in map`] }
            );
          }
          // test non-totp validation by rejecting all pingid requests
          if (mfaConstraint.type === 'pingid') {
            return new Response(403, {}, { errors: ['PingId MFA validation failed'] });
          }
          // validate totp passcode
          const passcode = mfa_payload[constraintId];
          if (mfaConstraint.type === 'totp') {
            if (passcode !== 'test') {
              const error = !passcode ? 'TOTP passcode not provided' : 'Incorrect TOTP passcode provided';
              return new Response(403, {}, { errors: [error] });
            }
          } else if (passcode) {
            // for okta and duo, reject if a passcode was provided
            return new Response(400, {}, { errors: ['Passcode should only be provided for TOTP MFA type'] });
          }
        }
        return authResponses[mfa_request_id];
      } catch (error) {
        console.log(error);
        return new Response(500, {}, { errors: ['Mirage Handler Error: /sys/mfa/validate'] });
      }
    },
    { timing: 3000 }
  );
}

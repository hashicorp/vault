import decodeConfigFromJWT from 'vault/utils/decode-config-from-jwt';
import { module, test } from 'qunit';

module('Unit | Util | decode config from jwt', function() {
  const PADDING_STRIPPED_TOKEN =
    'eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiaHR0cDovLzE5Mi4xNjguNTAuMTUwOjgyMDAiLCJleHAiOjE1MTczNjkwNzUsImlhdCI6MTUxNzM2NzI3NSwianRpIjoiN2IxZDZkZGUtZmViZC00ZGU1LTc0MWUtZDU2ZTg0ZTNjZDk2IiwidHlwZSI6IndyYXBwaW5nIn0.MIGIAkIB6s2zbohbxLimwhM6cg16OISK2DgoTgy1vHbTjPT8uG4hsrJndZp5COB8dX-djWjx78ZFMk-3a6Ij51su_By9xsoCQgFXV8y3DzH_YzYvdL9x38dMSWaVHpR_lpoKWsQnMvAukSchJp1FfHZQ8JcSkPu5IAVZdfwlG5esJ_ZOMxA3KIQFnA';
  const NO_PADDING_TOKEN =
    'eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiaHR0cDovLzEyNy4wLjAuMTo4MjAwIiwiZXhwIjoxNTE3NDM0NDA2LCJpYXQiOjE1MTc0MzI2MDYsImp0aSI6IjBiYmI1ZWMyLWM0ODgtMzRjYi0wMzY5LTkxZmJiMjVkZTFiYSIsInR5cGUiOiJ3cmFwcGluZyJ9.MIGHAkIBAGzB5EW6PolAi2rYOzZNvfJnR902WxprtRqnSF2E2I2ye9XLGX--L7npSBjBhnd27ocQ4ZO9VhfDIFqMzu1TNiwCQT52O6xAoz9ElRrq76PjkEHO4ns5_ZgjSKXuKaqdGysHYSlry8KEjWLGQECvZWg9LQeIf35jwqeQUfyJUfmwl5r_';
  const INVALID_JSON_TOKEN = `foo.${btoa({ addr: 'http://127.0.0.1' })}.bar`;

  test('it decodes token with no padding', function(assert) {
    const config = decodeConfigFromJWT(NO_PADDING_TOKEN);

    assert.ok(!!config, 'config was decoded');
    assert.ok(!!config.addr, 'config.addr is present');
  });

  test('it decodes token with stripped padding', function(assert) {
    const config = decodeConfigFromJWT(PADDING_STRIPPED_TOKEN);

    assert.ok(!!config, 'config was decoded');
    assert.ok(!!config.addr, 'config.addr is present');
  });

  test('it returns nothing if the config is invalid JSON', function(assert) {
    const config = decodeConfigFromJWT(INVALID_JSON_TOKEN);

    assert.notOk(config, 'config is not present');
  });
});

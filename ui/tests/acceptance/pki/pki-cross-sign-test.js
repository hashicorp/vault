/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, click, fillIn, find } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-issuer-cross-sign';
import { verifyCertificates } from 'vault/utils/parse-pki-cert';
module('Acceptance | pki/pki cross sign', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    this.parentMountPath = `parent-mount-${uuidv4()}`;
    this.oldParentIssuerName = 'old-parent-issuer'; // old parent issuer we're transferring from
    this.parentIssuerName = 'new-parent-issuer'; // issuer where cross-signing action will begin
    this.intMountPath = `intermediate-mount-${uuidv4()}`; // first input box in cross-signing page
    this.intIssuerName = 'my-intermediate-issuer'; // second input box in cross-signing page
    this.newlySignedIssuer = 'my-newly-signed-int'; // third input
    await enablePage.enable('pki', this.parentMountPath);
    await enablePage.enable('pki', this.intMountPath);

    await runCommands([
      `write "${this.parentMountPath}/root/generate/internal" common_name="Long-Lived Root X1" ttl=8960h issuer_name="${this.oldParentIssuerName}"`,
      `write "${this.parentMountPath}/root/generate/internal" common_name="Long-Lived Root X2" ttl=8960h issuer_name="${this.parentIssuerName}"`,
      `write "${this.parentMountPath}/config/issuers" default="${this.parentIssuerName}"`,
    ]);
  });

  hooks.afterEach(async function () {
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.intMountPath}`]);
    await runCommands([`delete sys/mounts/${this.parentMountPath}`]);
    await logout.visit();
  });

  test('it cross-signs an issuer', async function (assert) {
    const oldParent = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUMasDTM5kAFfObliTvzO1pNPNSpowDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgxMB4XDTIzMDUwMTIxMjEy\nNFoXDTIzMDYwMjIxMjE1NFowHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgx\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA20kuor8GsTO5yl3vok7K\n9XWgZIMUMRjNTfx5NbXQ0A7RsYaAJC4cyd+f/0bzzzbmHWv4t7M77KYTHzT4tYpK\nwr5v2jl2qbS2RKtmZfg7EE2R64oJ8A9BEPaFvAAPbvuWS9MMjr1BGQFdWMXGzMfL\nrhwPUi75LnZ7rmS8+SMBNZjZCzQAQtFkVOVPOoCq7KUJt++nOb8jEdshXGu2Ojom\nbDpB4out+/GLtEyTeK1XPwRY+iKEt1jTsy3GdlzpR93ZJ5YYWeyYNqzJDP0jhmPR\nd1FgFuvsLtX5FS19MHPV7x74j8sgnlKXUjxufZZ04Wtk9AV/xL8rXKLxon+OVsgH\n2QIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQUi021H8/CiaTX5rbr8pQB697EzdkwHwYDVR0jBBgwFoAUi021H8/CiaTX\n5rbr8pQB697EzdkwDQYJKoZIhvcNAQELBQADggEBAC+bkHDlXB7tp1tNNowXVzO/\ncSUtgmMr9BVDj3ITCQui2/lntgmHqHaHJ2k8RAzWNWih2lrzVMhRAlFQjHniQ5/I\nfjq/OvY5r0IqI80Urdnf4sbBJx7gPLJlgM52yZJRdXjc7E0mbxfGuhUNWDgy+MOJ\n9Qc9ug1go5UyQX+ehftoQi86X0wxenv4TMl93cy0UF7DQBYY0vtjdyo2PbcQGeh5\nuDmPdXWxz3c/LHPcbY0ysNtEdADK9iqVTxvq2bIwFamYjFadSJRAFuHTPkEqdTHE\nybeq46NgaivdNjHDsKVSJChiE8/4ZJtnLaHNvWM2P4A/+UyburjTX05SATC+EtY=\n-----END CERTIFICATE-----\n`;
    const newParent = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUXzLxW05qCjEM8+8dvx6dCk1WcGEwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgyMB4XDTIzMDUwMTIxMjEy\nNFoXDTIzMDYwMjIxMjE1NFowHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgy\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0hPlH+ImGjoHuKOt7ZZN\nsLk2F5EVE/BpshO4OzmIcRSABJgoBa8M0gcy5u6iC4zAFdXybd396Mi0LnQ1oytO\nfXQBNp3bXMhLq/LiF69fQ47e0+z5N2GaizWeQCTYlpDhU6woQ3ZCM7K6N8wzKGS9\n4KDZYehO8H8pTBVxOdTMikEpCj8GPLABZG1d1snF1P+vlaO34dr7WKpnGorDWsTw\nOTn/LzBDJBFnAZywtHHhd6SVrIPvR/uONBsI3CErlW63Hia8HlEVH9PwDPUzGEX7\nqekktXMDj+cKFqv+etfCn/4zORhJPQgkj+myE2riR8xOpvS7G1GKw/+o4T1IH2WE\nxwIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQUak5d37QOIAQq77iKPRMGotCm5RIwHwYDVR0jBBgwFoAUak5d37QOIAQq\n77iKPRMGotCm5RIwDQYJKoZIhvcNAQELBQADggEBAEhQQ6OKQGdeKYl6Yteim+3B\nEk3XLX+Va0mA3lu4mz64ctZHo+gMjzgeD5yJUKlJQRm/tLGi3k5B8mr8M/2WiBYS\nkNzCHBYyg71DxKTd5iimcpn1BRjh4yuMJyeQbhpwzp8qRY0Yt7R3yfWTvWIC7wBV\nkHrCFQ4/QM67P+rkP1tgG68fbvzfWFs/HGqwqb0bXgq5YS2H+eV9xHiMdqTieVdG\neelgAR/rOZStukMHsaQqBXJYG862ZDh/ZBLso6G9QcSLT0VzozOKI4cXtAbyg+HZ\nGmMLdlvgwqg5qRjdPW52ODRhDzEEcgdNBCeH8DvM2jk7EyELvVhhhYE4hOkcbwQ=\n-----END CERTIFICATE-----\n`;
    const oldInt = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUb0r8eP0ix4aineR1cjwQKs7LyzswDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgxMB4XDTIzMDUwMTIxMjEy\nNVoXDTIzMDYwMjIxMjE1NVowHTEbMBkGA1UEAxMSU2hvcnQtTGl2ZWQgSW50IFIx\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoS0AgZa+fiD4lscP44Nj\nLg5wU6alr2kibta48oqVhNaXcWZ4dLoHArMWdjhSRcDPJVqqT2IPK7fyAkq/re8I\nAH0rzM8BNOQPXK7XcMqdHn0cIYL1evVjCdq8i6VFw3/VRFVIJRpNRneSvPlViIw9\n+u5dJEVx5T5VTuadTDrBoXMib4+3yvbly9rKgb5lXvm+EWX6SaFY6tDvGLemANcK\niAW9AElmOq3YN1Byt4zCuknUppNyx8+2/eu6r0MIxnMJS+ZiMeWsUEzfYtsYVAdH\nx/7P/ynQmWaKCi0rDV6Rl1gQt6cjBbshsQ1aydjEKhpAl0f50dY/zOLccgOKCGKJ\n6QIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQUO5GvpYDyzbBUXJN3PZPF31xSux0wHwYDVR0jBBgwFoAUi021H8/CiaTX\n5rbr8pQB697EzdkwDQYJKoZIhvcNAQELBQADggEBACcSVeI86od/FyDYbe8kHPgl\nukBocM+1dycNGnAX7UxfMxP8Q8if3k1AECTqTwDH6BfbS6+SxNg/6GNIbIH/OQAO\ns/hrPIK7+VWeKDga9VJUAcUBAElHJqs7Y1ye75DLTVn3fjG8xL/MnSrU/1V/XUBZ\nLTvq+Hm2D28Ej+Sj6vqsEpUXCZM7dLcrtTx82g7fOoFnA+bzXxIOEggBCMVlZxAv\nS1ZTsdxsmhJMxYHovddpq4CuOUuKcDyqmIkfm8k3sPLxdSejTr6odfDKumUrX11n\n1F0ChQnJemAglP096LMs75uJIYC+R3AiWOkh6zGede5pz6U4o4KMFyHfHuOuShE=\n-----END CERTIFICATE-----\n`;
    const newInt = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUMC5/gh7BP3+NGl8w9+NNEXcujkYwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgyMB4XDTIzMDUwMTIxMjcw\nMFoXDTIzMDYwMjIxMjczMFowHTEbMBkGA1UEAxMSU2hvcnQtTGl2ZWQgSW50IFIx\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoS0AgZa+fiD4lscP44Nj\nLg5wU6alr2kibta48oqVhNaXcWZ4dLoHArMWdjhSRcDPJVqqT2IPK7fyAkq/re8I\nAH0rzM8BNOQPXK7XcMqdHn0cIYL1evVjCdq8i6VFw3/VRFVIJRpNRneSvPlViIw9\n+u5dJEVx5T5VTuadTDrBoXMib4+3yvbly9rKgb5lXvm+EWX6SaFY6tDvGLemANcK\niAW9AElmOq3YN1Byt4zCuknUppNyx8+2/eu6r0MIxnMJS+ZiMeWsUEzfYtsYVAdH\nx/7P/ynQmWaKCi0rDV6Rl1gQt6cjBbshsQ1aydjEKhpAl0f50dY/zOLccgOKCGKJ\n6QIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQUO5GvpYDyzbBUXJN3PZPF31xSux0wHwYDVR0jBBgwFoAUak5d37QOIAQq\n77iKPRMGotCm5RIwDQYJKoZIhvcNAQELBQADggEBAMHNCLpkk16OZRwMbaIOQzwI\nJk7ONcrpTTBrkvFuaCaZKyr+KLSj1jSGc5AtTqMr3nT9vIhnG8VnTDnn9nvn5OYx\nH1wX6+iaypF+4EZFSHOmkufwX1hxfVHf7Q70Tg1IQW3+BqjBj1aCP5QwLvf8lxj8\nDxBJiuBmYRdjJcUUCUztZJw9ccSUVBFMbn2qkUAKf4rIp8m7yYNdYC5Alu5ObDnA\nQI/MxBcJksyfUu8t+pgow983gtm9joT+0TcRdk+q/0fLqbzj+MunxC7Edzs6ySL7\noe2Nk1w2YxeSYhgmfTdAmMvOPcKAFnd8LJwe4VtCNqRDzjJ6rDORKh3Pvh7EHAk=\n-----END CERTIFICATE-----\n`;
    const leaf = `-----BEGIN CERTIFICATE-----\nMIIDSDCCAjCgAwIBAgIUPmoMZlVBR+y/E5/NsPPMjCe81RgwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSU2hvcnQtTGl2ZWQgSW50IFIxMB4XDTIzMDUwMTIxMjgw\nMVoXDTIzMDUwMzIxMjgzMVowFDESMBAGA1UEAxMJdGVzdC1sZWFmMIIBIjANBgkq\nhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt9VzrSlVEeivf7s8ZOFB6HrfM/GKT59n\najqPjKumrj6yypo5l+rTV36v/8nFJcDLSLszNPBqi907O7E+Z58LxzURZCLn8CEn\nerOhedlmQtqiD3GCO3bc0tnuUdXOnPAMMp9m3LFxyfS2VPGgM1+r3mBfjUaYuW2H\nijx87LLPPVUXZtwVce0+kll10K13vg6PrFiyfj6yJl/DYY/8ZTRmi5aG8WLyOePA\npWX/oWmzXKRgKvRI60vOHts5iIplsYkGEU5zNZB306SayV63Y66ujoWaQ8tjPePO\nA3av9Yf2hoE9D000WE6+STmK9EiZIXMEJOUkHwt6zECJ/iZvQ6cjGQIDAQABo4GI\nMIGFMA4GA1UdDwEB/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH\nAwIwHQYDVR0OBBYEFPsaZp3scYFdIo22W/A0YO7EULciMB8GA1UdIwQYMBaAFDuR\nr6WA8s2wVFyTdz2Txd9cUrsdMBQGA1UdEQQNMAuCCXRlc3QtbGVhZjANBgkqhkiG\n9w0BAQsFAAOCAQEAXSMkEEYArPi248Xa+q2EbAy5F/GvID0ND5byD5j2bt+BTrCv\nx9Zdo2dNbkuuLw5T7beucD9etyvCZ9J3KPkVkJmbKXO51YvX6ZBFslsKBoeFl2z3\nIJ+A3AkwTNSKRi11pPlCtvEp++9RovME4PWVF9E5I5VmqWytSJTURrDxmTnn+BJZ\nh4QaXWDNArLOno5fJaiSpgwen0R24UhsfxE5xVv46K4UDljm0StiswG2rvQXRHFM\nmimCQLbzEAs6VQrCisZqzc78AqjD51X3xk7+DLQFygiPUHWEzH9maZ1q4EW7vFbZ\nCIWVg9ji3w2Tccf2VVHMB7WP768jZRn+M1RXVw==\n-----END CERTIFICATE-----\n`;
    await verifyCertificates(oldParent, oldInt, leaf);
    await verifyCertificates(newParent, newInt, leaf);
    // console.log(isValidAgainstOldChain, isValidAgainstNewChain);

    // configure mounts to make them cross-signable
    await visit(`/vault/secrets/${this.intMountPath}/pki/configuration/create`);
    await click(SELECTORS.optionByKey('generate-csr'));
    await fillIn(SELECTORS.typeField, 'internal');
    await fillIn(SELECTORS.inputByName('commonName'), 'Short-Lived Int R1');
    await click('[data-test-save]');
    const csr = find('[data-test-value-div="CSR"] [data-test-copy-button]').getAttribute(
      'data-clipboard-text'
    );
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.oldParentIssuerName}/sign`);
    await fillIn(SELECTORS.inputByName('csr'), csr);
    await fillIn(SELECTORS.inputByName('format'), 'pem_bundle');
    await click('[data-test-pki-sign-intermediate-save]');
    const pemBundle = find('[data-test-value-div="CA Chain"] [data-test-copy-button]')
      .getAttribute('data-clipboard-text')
      .replace(/,/, '\n');
    await visit(`vault/secrets/${this.intMountPath}/pki/configuration/create`);
    await click(SELECTORS.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', pemBundle);
    await click(SELECTORS.importSubmit);
    await visit(`vault/secrets/${this.intMountPath}/pki/issuers`);
    await click('[data-test-is-default]');
    // name default issuer of intermediate
    const oldIntIssuerId = find(SELECTORS.rowValue('Issuer ID')).innerText;
    const oldIntCert = find('[data-test-value-div="Certificate"] [data-test-copy-button]').getAttribute(
      'data-clipboard-text'
    );
    await click('[data-test-pki-issuer-configure]');
    await fillIn(SELECTORS.inputByName('issuerName'), this.intIssuerName);
    await click('[data-test-save]');

    // perform cross-sign
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.parentIssuerName}/cross-sign`);
    await fillIn(SELECTORS.input('intermediateMount'), this.intMountPath);
    await fillIn(SELECTORS.input('intermediateIssuer'), this.intIssuerName);
    await fillIn(SELECTORS.input('newCrossSignedIssuer'), this.newlySignedIssuer);
    await click(SELECTORS.submitButton);
    assert
      .dom(`${SELECTORS.signedIssuerCol('intermediateMount')} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.intMountPath}/pki/overview`);
    assert
      .dom(`${SELECTORS.signedIssuerCol('intermediateIssuer')} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.intMountPath}/pki/issuers/${oldIntIssuerId}/details`);

    // get certificate data of newly signed issuer
    await click(`${SELECTORS.signedIssuerCol('newCrossSignedIssuer')} a`);
    const newIntCert = find('[data-test-value-div="Certificate"] [data-test-copy-button]').getAttribute(
      'data-clipboard-text'
    );

    // get certificate data of parent issuers
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.parentIssuerName}/details`);
    const newParentCert = find('[data-test-value-div="Certificate"] [data-test-copy-button]').getAttribute(
      'data-clipboard-text'
    );
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.oldParentIssuerName}/details`);
    const oldParentCert = find('[data-test-value-div="Certificate"] [data-test-copy-button]').getAttribute(
      'data-clipboard-text'
    );

    // create a role to issue a leaf certificate
    const myRole = 'some-role';
    await runCommands([
      `write ${this.intMountPath}/roles/${myRole} \
    issuer_ref=${this.newlySignedIssuer}\
    allow_any_name=true \
    max_ttl="720h"`,
    ]);
    await visit(`vault/secrets/${this.intMountPath}/pki/roles/${myRole}/generate`);
    await fillIn(SELECTORS.inputByName('commonName'), 'my-leaf');
    await fillIn('[data-test-ttl-value="TTL"]', '3600');
    await click('[data-test-pki-generate-button]');
    const myLeafCert = find('[data-test-value-div="Certificate"] [data-test-copy-button]').getAttribute(
      'data-clipboard-text'
    );

    assert.true(await verifyCertificates(oldParentCert, oldIntCert, myLeafCert));
    assert.true(await verifyCertificates(newParentCert, newIntCert, myLeafCert));
  });
});

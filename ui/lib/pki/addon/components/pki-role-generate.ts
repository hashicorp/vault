import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { parseCertificate } from 'vault/helpers/parse-pki-cert';
import Router from '@ember/routing/router';
import Store from '@ember-data/store';
import PkiRoleAdapter from 'vault/adapters/pki/role';
import { HTMLElementEvent } from 'forms';

interface Args {
  role: string;
  backend: string;
  onSuccess: CallableFunction;
}
interface CertResponse {
  ca_chain: string[];
  certificate: string;
  expiration: number;
  issuing_ca: string;
  private_key: string;
  private_key_type: string;
  serial_number: string;
}
interface Certificate {
  contents: string;
  commonName: string;
  issueDate: Date;
  serialNumber: string;
  notValidBefore: Date;
  notValidAfter: Date;
}
export default class PkiRoleGenerate extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @tracked commonName = '';

  @tracked certificate: Certificate | null = null;

  // TODO: only for testing, do not commit
  // constructor(owner: unknown, args: Args) {
  //   super(owner, args);
  //   const res: CertResponse = {
  //     ca_chain: [
  //       '-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gAwIBAgIUTPb/quDfym2Js9NeauQrsFldHIQwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjIxMjA3MTUwODE4WhcNMzIx\nMjA0MTUwODQ4WjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBANUZ/+NWdRzjKd02RtX6rXbPi/Xf+hRL1J0gZ+78\ngHTcBi/AiSwfKxuhzh024/lhq91zC3CWMVy4eZirReoinVhbdCf759kT78DwEHv/\nl71w0rN9RBL41PsAfAJSXz3CQxOdR35OOmVa1gbOmW9z9Nz6CZ8d129Hp2DM7Cka\nd0pbKNfZ/c3nTVNXUz/lsShrUq1WYbsKe0SMwSbUq9UJiW6ZckRV9/0CNxt3ol+h\nAFWPIxJwnJud053ngz1ULKRNcZpiQaFVaqKTt56nvzZWukx3cFkVI88l5YQ754VJ\nM2FMoB8h0gRZUXg4X+zbQ7NsjvVhe/Jzqreu25KHkdgP2yECAwEAAaN7MHkwDgYD\nVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFCxYI4iV1hR4\n7vLqttxt2Dyu/jOgMB8GA1UdIwQYMBaAFCxYI4iV1hR47vLqttxt2Dyu/jOgMBYG\nA1UdEQQPMA2CC2V4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQBUh2mTl+6G\naG3LQ4YS76K1wcBjfSjwRp/o+J8lFl+HfRoqm/DBIVOUEGP2D/Q9cm07HHkYDzo3\noXofrNc5LWDQi1L8sIa7uSygNMn+4+k2h6YojZ1ji/RiwQ0VRWw40vlNHFF+grEs\nuzAETZpUAvL+WQTfp+q1nFBaIuwxR7xpUM4G/YX9ab41JU+S3D7vfXjhz5NziAJv\nrfPdI1n1lEJVOGwM51Xu4klVBiGsfyTLHrC4yys6AD99APezgjUmUEzWxJrOV9sG\n8l550kHO6WkflqcAK3RDYwB7hFFnoQs1wkoazYQB99QlGvWBDtDmkaD+VLTALQ8D\ne4w0wr0MoxLN\n-----END CERTIFICATE-----',
  //     ],
  //     certificate:
  //       '-----BEGIN CERTIFICATE-----\nMIIDSzCCAjOgAwIBAgIUZ4jT5jT3Dot395eyj57DDg2uOoUwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjIxMjA3MTU1MzU4WhcNMjMw\nMTA4MTU1NDI3WjAZMRcwFQYDVQQDEw5leGFtcGxlLmNvbW1vbjCCASIwDQYJKoZI\nhvcNAQEBBQADggEPADCCAQoCggEBAMcLADtxqer6+2c+V/Od+20PUTKthlEaHaRw\n4GlIb7p30VtcH8ANxg4Y2a2AD8Dbyr3Mav0c4X/zkD4CuJnlOuH6jUKn4obMGGTq\nVM7of0jfdQ4VYAqQUUfthpoh2asUrc/RGbP2X8nXW+7htMUC9/BsKBzeWHM3hZF8\nVsi705/Vx5+G1du+spdiUJ3VUd+kLqOeDwA4E57baJu+81MYusM+T7ED/S3W3VT/\n965Bu+GpxNyl4dzrsKpHU+8/jn5FgGqrdFFGgnQYYBTkJxAvD6AIYDayhkE+eYco\nOoEwmG+9ptzuJ87I2ztSIRbrelXXKn9dyNH8S4DWOvcWRV57bMcCAwEAAaOBjTCB\nijAOBgNVHQ8BAf8EBAMCA6gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC\nMB0GA1UdDgQWBBRlTSr5wVoQvOtb/mBVwuIoNJSUJTAfBgNVHSMEGDAWgBQsWCOI\nldYUeO7y6rbcbdg8rv4zoDAZBgNVHREEEjAQgg5leGFtcGxlLmNvbW1vbjANBgkq\nhkiG9w0BAQsFAAOCAQEAMuzjR0noFVIKjNbBYlYTilR7l1KdFWm1cUv0xq3CKSbf\nUD1l2ooHQuqKZWpP1hTFf0kEsIgGEecNd13WIwxmgMKRu2p9CXtpFJoCKc5005dW\n2j2kI27wqAh4fZ1Gf5yvA3We/dd+3S4CEmrBq6nZ8T6bP5/Q3otNk+y2XZ7vn7fC\nmU5XERQ5dwWDhHgJvwCnCKjINtafbed+RGJlaIgfXLsCZWLxZqLpi2hITJgeX8qv\nUoMfQkSrQ8lh3PjN8Qon/lbcjP2bMVnu56tDGvA10ZlfoaB9Raoqz7NP19rzaPP8\nsJGQv+d3ASVPxteLijUE7P05p6EYKo7mh9u9FUpenQ==\n-----END CERTIFICATE-----',
  //     expiration: 1673193267,
  //     issuing_ca:
  //       '-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gAwIBAgIUTPb/quDfym2Js9NeauQrsFldHIQwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjIxMjA3MTUwODE4WhcNMzIx\nMjA0MTUwODQ4WjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBANUZ/+NWdRzjKd02RtX6rXbPi/Xf+hRL1J0gZ+78\ngHTcBi/AiSwfKxuhzh024/lhq91zC3CWMVy4eZirReoinVhbdCf759kT78DwEHv/\nl71w0rN9RBL41PsAfAJSXz3CQxOdR35OOmVa1gbOmW9z9Nz6CZ8d129Hp2DM7Cka\nd0pbKNfZ/c3nTVNXUz/lsShrUq1WYbsKe0SMwSbUq9UJiW6ZckRV9/0CNxt3ol+h\nAFWPIxJwnJud053ngz1ULKRNcZpiQaFVaqKTt56nvzZWukx3cFkVI88l5YQ754VJ\nM2FMoB8h0gRZUXg4X+zbQ7NsjvVhe/Jzqreu25KHkdgP2yECAwEAAaN7MHkwDgYD\nVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFCxYI4iV1hR4\n7vLqttxt2Dyu/jOgMB8GA1UdIwQYMBaAFCxYI4iV1hR47vLqttxt2Dyu/jOgMBYG\nA1UdEQQPMA2CC2V4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQBUh2mTl+6G\naG3LQ4YS76K1wcBjfSjwRp/o+J8lFl+HfRoqm/DBIVOUEGP2D/Q9cm07HHkYDzo3\noXofrNc5LWDQi1L8sIa7uSygNMn+4+k2h6YojZ1ji/RiwQ0VRWw40vlNHFF+grEs\nuzAETZpUAvL+WQTfp+q1nFBaIuwxR7xpUM4G/YX9ab41JU+S3D7vfXjhz5NziAJv\nrfPdI1n1lEJVOGwM51Xu4klVBiGsfyTLHrC4yys6AD99APezgjUmUEzWxJrOV9sG\n8l550kHO6WkflqcAK3RDYwB7hFFnoQs1wkoazYQB99QlGvWBDtDmkaD+VLTALQ8D\ne4w0wr0MoxLN\n-----END CERTIFICATE-----',
  //     private_key:
  //       '-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAxwsAO3Gp6vr7Zz5X8537bQ9RMq2GURodpHDgaUhvunfRW1wf\nwA3GDhjZrYAPwNvKvcxq/Rzhf/OQPgK4meU64fqNQqfihswYZOpUzuh/SN91DhVg\nCpBRR+2GmiHZqxStz9EZs/Zfyddb7uG0xQL38GwoHN5YczeFkXxWyLvTn9XHn4bV\n276yl2JQndVR36Quo54PADgTnttom77zUxi6wz5PsQP9LdbdVP/3rkG74anE3KXh\n3OuwqkdT7z+OfkWAaqt0UUaCdBhgFOQnEC8PoAhgNrKGQT55hyg6gTCYb72m3O4n\nzsjbO1IhFut6Vdcqf13I0fxLgNY69xZFXntsxwIDAQABAoIBAQCZwCWtxV875CPO\n3JmT1bUhyXMvclsOyt2a6JZwvUORBnzx1XexIvKacRe0rfd9QkqZ0g3S9zw5WitR\nu0hdmHSjmqcDliuldIJjcZ+RNIceA36oIlrNziz7Ir+W0A8T2CVrIlp4aRgVEgYx\nwKeih2h+jw1tP1OTrI/Akgt3g581Fs8u99aGst0R7d/L200emQuODoWqbEdzlGJK\nwiCQYWJKExUTl2ReYqkkbNOaPp/7Oi193ZXvO/FsULwcdfxsb1b9tpBpcATk19Dm\ngsHM0B/EQMn97BO0evV0TpfxBLf0lvwDh4fye1cMfZQ2S+wiNbCBd1pHV2co0SlI\nxog4pJQBAoGBAMslE+3xSK9tgEwmAzEAH5WrTu7soXg6/UsWn3xw/w7anZiKw2AS\niIl/4znbxmKwfmZZ7BwUz9spkRQfD0Tbf/F10l4DVn+SaKlc6Wi1KlMAHUHXnjto\ntJKXuQhj5wt+eBDjruB3zP8JHJdIQRC84zhdZGUc8uObCJhqq+43wqepAoGBAPrU\ntYxVNwuN1STdJS2GpbgqV3Cy5AdTOJ/KMqnsxUTFZ5ja8VzA+X7XREQY8YGto47F\n/KqH7/oxu4v2xa8VxD8c6TKffi3SAaWS/ChuWKh0taVWEUpbigOAFNaGpjpypxRi\nFsP4yaXI6+gMUY+NxzBRBjHfgZCs4Zvxy/56UXbvAoGAFDKXjKzUwTxt6SROZOzS\nNxtVOcQlOcMDtBeHu+OwOFXcHXKOglrVYHZqrTIAw4cwyhReuVwIXo3/crSz2/DH\nA8bnJ5nFW+G+rjgirNp2XtJAFm/Nt7JtYbXcG81zB12HqoY4uPCwXRsW0KPKdFOT\nc+M1PChoreCYNi9E8OZyYCECgYA1tXA+YUzNE4ytPREl42v+uEpK3nNHQkGgrXoV\nupYu+JoLN+5wrv19dHiwoCquWtDn1Gsa1MrE5vtCqA+CQwXngbhJV697/jjODGAk\nBCTFxV/TzE8dfeZag4Vyvhg/8abnDW3UfqQm9JPW9zRLqc8aAG90JGio4uvYKXlF\nv0lMTwKBgQC7cgeKfYz8cZcrDn0lDdJ0zb+cgxG5h6LtF4e5Th8NUtKHrA2FbsYI\nS4BlsHBhwR40llfU3JuGowFuez6T47PKyc9n/tPh1GOxL69MRH2nhzHIl5pAe4RC\nxIQaO8h0QMbbODuZH6if0VKsX5fNAoW9WBCLt6W2ENuKCApFZKApNA==\n-----END RSA PRIVATE KEY-----',
  //     private_key_type: 'rsa',
  //     serial_number: '67:88:d3:e6:34:f7:0e:8b:77:f7:97:b2:8f:9e:c3:0e:0d:ae:3a:85',
  //   };
  //   this.displayCertificate(res);
  // }

  displayCertificate(response: CertResponse) {
    const parsed = parseCertificate(response.certificate);
    this.certificate = {
      contents: response.certificate,
      commonName: parsed.common_name || '',
      issueDate: new Date(parsed.issue_date),
      serialNumber: response.serial_number,
      notValidAfter: new Date(parsed.expiry_date),
      notValidBefore: new Date(parsed.issue_date),
    };
  }

  @task
  *save(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    const { role, backend, onSuccess } = this.args;
    const adapter = this.store.adapterFor('pki/role') as PkiRoleAdapter;
    try {
      const response: CertResponse = yield adapter.generateCertificate(backend, role, {
        common_name: this.commonName,
      });
      this.displayCertificate(response);
      onSuccess();
    } catch (e) {
      // TODO: show error
      console.log('ERROR', e);
    }
  }

  @task
  *revoke() {
    const { backend } = this.args;
    const adapter = this.store.adapterFor('pki/role') as PkiRoleAdapter;
    const payload = this.certificate?.serialNumber
      ? { serial_number: this.certificate.serialNumber }
      : { certificate: this.certificate?.contents };
    try {
      yield adapter.revokeCertificate(backend, payload);
      this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
    } catch (e) {
      // TODO: show error
      console.log('ERROR', e);
    }
  }

  @action download() {
    // TODO
  }

  @action cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
  }
}

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { waitFor } from '@ember/test-waiters';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { addToArray } from 'vault/helpers/add-to-array';
/**
 * @module PkiIssuerCrossSign
 * PkiIssuerCrossSign components render from a parent issuer's details page to cross-sign an intermediate issuer (from a different mount).
 * The component reads an existing intermediate issuer, cross-signs it with a parent issuer and imports the new
 * issuer into the existing intermediate mount using three inputs from the user:
 * intermediateMount (the mount path where the issuer to be cross signed lives)
 * intermediateIssuer (the name of the intermediate issuer, located in the above mount)
 * newCrossSignedIssuer (the name of the to-be-cross-signed, new issuer)
 *
 * The requests involved and how those inputs are used:
 * 1. Read an existing intermediate issuer
 *    -> GET /:intermediateMount/issuer/:intermediateIssuer
 * 2. Create a new CSR based on this existing issuer ID
 *    -> POST /:intermediateMount/intermediate/generate/existing
 * 3. Sign it with the new parent issuer, minting a new certificate.
 *    -> POST /this.args.parentIssuer.backend/issuer/this.args.parentIssuer.issuerRef/sign-intermediate
 * 4. Import it back into the existing mount
 *    -> POST /:intermediateMount/issuers/import/bundle
 * 5. Read the imported issuer
 *    -> GET /:intermediateMount/issuer/:issuer_id
 * 6. Update this issuer with the newCrossSignedIssuer
 *    -> POST /:intermediateMount/issuer/:issuer_id
 *
 * @example
 * ```js
 * <PkiIssuerCrossSign @parentIssuer={{this.model}} />
 * ```
 * @param {object} parentIssuer - the model of the issuing certificate that will sign the issuer to-be cross-signed
 */

export default class PkiIssuerCrossSign extends Component {
  @service store;
  @tracked formData = [];
  @tracked signedIssuers = [];
  @tracked intermediateIssuers = {};
  @tracked validationErrors = [];

  inputFields = [
    {
      label: 'Mount path',
      key: 'intermediateMount',
      placeholder: 'Mount path',
      helpText: 'The mount in which your new certificate can be found.',
    },
    {
      label: "Issuer's current name",
      key: 'intermediateIssuer',
      placeholder: 'Current issuer name',
      helpText: 'The API name of the previous intermediate which was cross-signed.',
    },
    {
      label: 'New issuer name',
      key: 'newCrossSignedIssuer',
      placeholder: 'Enter a new issuer name',
      helpText: `This is your new issuer’s name in the API.`,
    },
  ];

  get statusCount() {
    const error = this.signedIssuers.filter((issuer) => issuer.hasError).length;
    const success = this.signedIssuers.length - error;
    return `${success} successful, ${error} ${error === 1 ? 'error' : 'errors'}`;
  }

  @task
  @waitFor
  *submit(e) {
    e.preventDefault();
    this.signedIssuers = [];
    this.validationErrors = [];

    // Validate name input for new issuer does not already exist in mount
    for (let row = 0; row < this.formData.length; row++) {
      const { intermediateMount, newCrossSignedIssuer } = this.formData[row];
      const issuers = yield this.store
        .query('pki/issuer', { backend: intermediateMount })
        .then((resp) => resp.map(({ issuerName, issuerId }) => ({ issuerName, issuerId })))
        .catch(() => []);

      // for cross-signing error handling we want to record the list of issuers before the process starts
      this.intermediateIssuers[intermediateMount] = issuers;
      this.validationErrors = addToArray(this.validationErrors, {
        newCrossSignedIssuer: this.nameValidation(newCrossSignedIssuer, issuers),
      });
    }
    if (this.validationErrors.any((row) => !row.newCrossSignedIssuer.isValid)) return;

    // iterate through submitted data and cross-sign each certificate
    for (let row = 0; row < this.formData.length; row++) {
      const { intermediateMount, intermediateIssuer, newCrossSignedIssuer } = this.formData[row];
      try {
        // returns data from existing and newly cross-signed issuers
        // { intermediateIssuer: existingIssuer, newCrossSignedIssuer: crossSignedIssuer, intermediateMount: intMount }
        const data = yield this.crossSignIntermediate(
          intermediateMount,
          intermediateIssuer,
          newCrossSignedIssuer
        );
        this.signedIssuers = addToArray(this.signedIssuers, { ...data, hasError: false });
      } catch (error) {
        this.signedIssuers = addToArray(this.signedIssuers, {
          ...this.formData[row],
          hasError: errorMessage(error),
          hasUnsupportedParams: error.cause ? error.cause.map((e) => e.message).join(', ') : null,
        });
      }
    }
  }

  @action
  async crossSignIntermediate(intMount, intName, newCrossSignedIssuer) {
    // 1. Fetch issuer we want to sign
    // What/Recovery: any failure is early enough that you can bail safely/normally.
    const existingIssuer = await this.store.queryRecord('pki/issuer', {
      backend: intMount,
      id: intName,
    });

    // Return if user is attempting to self-sign issuer
    if (existingIssuer.issuerId === this.args.parentIssuer.issuerId) {
      throw new Error('Cross-signing a root issuer with itself must be performed manually using the CLI.');
    }

    // Translate certificate values to API parameters to pass along: CSR -> Signed CSR -> Cross-Signed issuer
    // some of these values do not apply to a CSR, but pass anyway. If there is any issue parsing the certificate,
    // (ex. the certificate contains unsupported values) direct user to manually cross-sign via CLI
    const certData = parseCertificate(existingIssuer.certificate);
    if (certData.parsing_errors.length > 0) {
      throw new Error('Certificate must be manually cross-signed using the CLI.', {
        cause: certData.parsing_errors,
      });
    }

    // 2. Create the new CSR
    // What/Recovery: any failure is early enough that you can bail safely/normally.
    const newCsr = await this.store
      .createRecord('pki/action', {
        keyRef: existingIssuer.keyId,
        commonName: existingIssuer.commonName,
        type: 'existing',
        ...certData,
      })
      .save({
        adapterOptions: { actionType: 'generate-csr', mount: intMount, useIssuer: false },
      })
      .then(({ csr }) => csr);

    // 3. Sign newCSR with correct parent to create cross-signed cert, "issuing"
    // an intermediate certificate.
    // What/Recovery: any failure is early enough that you can bail safely/normally.
    const signedCaChain = await this.store
      .createRecord('pki/action', {
        csr: newCsr,
        commonName: existingIssuer.commonName,
        ...certData,
      })
      .save({
        adapterOptions: {
          actionType: 'sign-intermediate',
          mount: this.args.parentIssuer.backend,
          issuerRef: this.args.parentIssuer.issuerRef,
        },
      })
      .then(({ caChain }) => caChain.join('\n'));

    // 4. Import the newly cross-signed cert to become an issuer
    // What/Recovery:
    //   1. Permission issue -> give the cert (`signedCaChain`) to the user,
    //      let them import & name. (Issue you have is that you already issued
    //      it (step 3) and so "undo" would mean revoking the cert, which
    //      you might not have permissions to do either).
    //
    //   2. CRL rebuilding fails ("the CRL" in error message). Server returns
    //      an error, we wanted the CRL rebuilt -- but the issuer was still
    //      imported anyways. Only way to detect would be to do a list issuers
    //      before and after. Recovery would be on the operator in this case;
    //      reproduce the error and let them deal with it.
    //
    // End result: user should solve this issue, but we shouldn't undo anything
    // either.
    //
    //    -> For 1 though, make sure to give the `signedCaChain` in the
    //       error message for them.
    //    -> For 2, you could list before and after to find the id of the
    //       new issuer(s) so they can name them and fix any issues with
    //       them.
    //
    // If its not a permissions error _and_ you did two lists, not finding
    // a new issuer...
    //
    //    -> Unknown error. Could give them `signedCaChain` and serial of
    //       the newly issued intermediate CA, so that they can do recovery
    //       as they'd like.
    const issuerId = await this.store
      .createRecord('pki/action', { pemBundle: signedCaChain })
      .save({ adapterOptions: { actionType: 'import', mount: intMount, useIssuer: true } })
      .then((importedIssuer) => {
        return Object.keys(importedIssuer.mapping).find(
          // matching key is the issuer_id
          (key) => importedIssuer.mapping[key] === existingIssuer.keyId
        );
      })
      .catch((e) => {
        console.debug('CA_CHAIN \n', signedCaChain); // eslint-disable-line
        throw new Error(`${errorMessage(e)} See console for signed ca_chain data.`);
      });

    // 5. Fetch issuer imported above by issuer_id, name and save
    // Recovery: cosmetic issue; can let the user deal with it. Usually
    // fails because the name is in use.
    // Pre-fix: list all issuers, check the desired name isn't either
    // an existing issuer_id or an issuer_name.
    const crossSignedIssuer = await this.store.queryRecord('pki/issuer', { backend: intMount, id: issuerId });
    crossSignedIssuer.issuerName = newCrossSignedIssuer;
    await crossSignedIssuer.save({ adapterOptions: { mount: intMount } });

    // 6. Return the data to our caller.
    return {
      intermediateIssuer: existingIssuer,
      newCrossSignedIssuer: crossSignedIssuer,
      intermediateMount: intMount,
    };
  }

  @action
  reset() {
    this.signedIssuers = [];
    this.validationErrors = [];
    this.formData = [];
  }

  nameValidation(nameInput, existing) {
    if (existing.any((i) => i.issuerName === nameInput || i.issuerId === nameInput))
      return {
        errors: [`Issuer reference '${nameInput}' already exists in this mount.`],
        isValid: false,
      };
    return { errors: [], isValid: true };
  }
}

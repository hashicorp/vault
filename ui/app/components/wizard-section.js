import templateOnly from '@ember/component/template-only';

/**
 * @module WizardSection
 * WizardSection components are instruction areas for the wizard.
 *
 * @example
 * ```js
 * <WizardSection
 *  @headerText="Enable secrets Engine"
 *  @docText="Docs: Secret Engine"
 *  @docPath: "docs/secret/index.html"
 *  @instructions="select and engine"
 *  @class="has-bottom-margin-l"/>
 * ```
 * @param {string} [headerText] - Title text.
 * @param {string} [docText] - Text for docs link.
 * @param {string} [docPath] - Link for the docs.
 * @param {string} [instructions] - what the user is to do in this step. Under the title What to do.
 * @param {string} [class] - class to add for section.
 */

export default templateOnly();

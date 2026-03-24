# Forms

- [Guidelines](#guidelines)

## Guidelines

- Render `FlashMessage` on success
- Handling errors/validation messages:
  - Render API errors using a `<MessageError>` or `Hds::Alert` at the top of forms
  - Display validation error messages `onsubmit` (not `onchange` for inputs)
  - Render an `<AlertInline>` [beside](../lib/pki/addon/components/pki-role-generate.hbs) form buttons, especially if the error banner is hidden from view (long forms). Message options:
    - The `invalidFormMessage` from a model's `validate()` method that includes an error count
    - Generic message for API errors or forms without model validations: 'There was an error submitting this form.'
  - Add `has-error-border` class to invalid inputs

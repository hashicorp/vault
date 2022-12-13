// accepts an error and returns error.errors joined with a comma, error.message or a fallback message
export default function (error, fallbackMessage = 'An error occurred, please try again') {
  if (error instanceof Error && error?.errors) {
    return error.errors.join(', ');
  }
  return error?.message || fallbackMessage;
}

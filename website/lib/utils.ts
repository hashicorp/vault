export const isInternalLink = (link: string): boolean => {
  if (
    link.startsWith('/') ||
    link.startsWith('#') ||
    link.startsWith('https://vaultproject.io') ||
    link.startsWith('https://www.vaultproject.io')
  ) {
    return true
  }
  return false
}

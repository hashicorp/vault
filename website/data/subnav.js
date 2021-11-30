export default [
  { text: 'Overview', url: '/' },
  {
    text: 'Use Cases',
    submenu: [
      { text: 'Secrets Management', url: '/use-cases/secrets-management' },
      { text: 'Data Encryption', url: '/use-cases/data-encryption' },
      {
        text: 'Identity-based Access',
        url: '/use-cases/identity-based-access',
      },
    ],
  },
  {
    text: 'Enterprise',
    url: 'https://www.hashicorp.com/products/vault/enterprise',
  },
  'divider',
  { text: 'Tutorials', url: 'https://learn.hashicorp.com/vault' },
  { text: 'Docs', url: '/docs' },
  { text: 'API', url: '/api-docs' },
  { text: 'Community', url: '/community' },
]

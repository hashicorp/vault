# PKI filenames

> delete this file when PKI redesign is complete

Now that all things PKI live in a `/pki` folder we are removing `pki` from the filename in any newly created files for the pki redesign (ex. 'role.js' instead of 'pki-role.js'). Aside from `cert.js` all of the old pki files are prepended with `pki-`

Old files:

├── models/
│   ├── pki/
│   │   ├── cert.js
│   │   ├── pki-config.js
│   │   ├── pki-role.js
│   ├── pki-ca-certificate-sign.js
│   ├── pki-ca-certificate.js
│   ├── pki-certificate-sign.js
├── serializers/
│   ├── pki/
│   │   ├── cert.js
│   │   ├── pki-config.js
│   │   ├── pki-role.js
├── adapters/
│   ├── pki/
│   │   ├── cert.js
│   │   ├── pki-config.js
│   │   ├── pki-role.js
│   ├── pki-ca-certificate-sign.js
│   ├── pki-ca-certificate.js
│   ├── pki-certificate-sign.js
│   ├── pki.js

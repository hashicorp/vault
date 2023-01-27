To generate a declaration file run `yarn tsc <javascript file to declare>  --declaration --allowJs --emitDeclarationOnly --outDir <type file location>`

For example, the following command generates a declaration file called base.d.ts for the pki certificate base.js model:

`yarn tsc ./app/models/pki/certificate/base.js  --declaration --allowJs --emitDeclarationOnly --outDir types/vault/models/pki/certificate`

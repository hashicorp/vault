import babylonParser from './jscodeshift-babylon-parser';

// use babylon parser with decorators-legacy plugin
export const parser = babylonParser;

// checks for access of specified service on this
// injects service if not present and imports inject as service if needed
// example usage - npx jscodeshift -t ./scripts/codemods/inject-service.js ./app/**/*.js --service=store
// --service arg is required with name of service
// pass -d for dry run (no files transformed)

export default function transformer({ source }, api, { service }) {
  if (!service) {
    throw new Error('Missing service arg. Pass --service=store for example to script');
  }
  const j = api.jscodeshift;
  const filterForService = (path) => {
    return j(path.value).find(j.MemberExpression, {
      object: {
        type: 'ThisExpression',
      },
      property: {
        name: service,
      },
    }).length;
  };
  const recastOptions = {
    reuseWhitespace: false,
    wrapColumn: 110,
    quote: 'single',
    trailingComma: true,
  };
  let didInjectService = false;

  // find class bodies and filter down to ones that access service
  const classesAccessingService = j(source).find(j.ClassBody).filter(filterForService);

  if (classesAccessingService.length) {
    // filter down to class bodies where service is not injected
    const missingService = classesAccessingService.filter((path) => {
      return !j(path.value)
        .find(j.ClassProperty, {
          key: {
            name: service,
          },
        })
        .filter((path) => {
          // ensure service property belongs to service decorator
          return path.value.decorators.find((path) => path.expression.name === 'service');
        }).length;
    });

    if (missingService.length) {
      // inject service
      const serviceInjection = j.classProperty(j.identifier(`@service ${service}`), null);
      // adding a decorator this way will force injection down to a new line and then add a new line
      // leaving in just in case it's needed
      // serviceInjection.decorators = [j.decorator(j.identifier('service'))];

      source = missingService
        .forEach((path) => {
          path.value.body.unshift(serviceInjection);
        })
        .toSource();

      didInjectService = true;
    }
  }

  // find .extend object expressions and filter down to ones that access this[service]
  const objectsAccessingService = j(source)
    .find(j.CallExpression, {
      callee: {
        type: 'MemberExpression',
        property: {
          name: 'extend',
        },
      },
    })
    .filter(filterForService)
    .find(j.ObjectExpression)
    .filter((path) => {
      // filter object expressions that belong to .extend
      // otherwise service will also be injected in actions: { } block of component for example
      const callee = path.parent.value.callee;
      return callee && callee.property?.name === 'extend';
    });

  if (objectsAccessingService.length) {
    // filter down to objects where service is not injected
    const missingService = objectsAccessingService.filter((path) => {
      return !j(path.value).find(j.ObjectProperty, {
        key: {
          name: service,
        },
        value: {
          callee: {
            name: 'service',
          },
        },
      }).length;
    });

    if (missingService.length) {
      // inject service
      const serviceInjection = j.objectProperty(
        j.identifier(service),
        j.callExpression(j.identifier('service'), [])
      );

      source = missingService
        .forEach((path) => {
          path.value.properties.unshift(serviceInjection);
        })
        .toSource(recastOptions);

      didInjectService = true;
    }
  }

  // if service was injected here check if inject has been imported
  if (didInjectService) {
    const needsImport = !j(source).find(j.ImportSpecifier, {
      imported: {
        name: 'inject',
      },
    }).length;

    if (needsImport) {
      const injectionImport = j.importDeclaration(
        [j.importSpecifier(j.identifier('inject'), j.identifier('service'))],
        j.literal('@ember/service')
      );

      const imports = j(source).find(j.ImportDeclaration);
      source = imports
        .at(imports.length - 1)
        .insertAfter(injectionImport)
        .toSource(recastOptions);
    }
  }

  return source;
}

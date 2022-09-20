export const statuses = [
  'connected',
  'disconnected since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: UNKNOWN',
  'disconnected since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: some other error other than unknown',
  'connecting since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: dial tcp [::1]:28083: connect: connection refused',
  'connecting since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: principal does not have permission to register as provider: rpc error: code = PermissionDenied desc =',
  'connecting since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: failed to get access token: oauth2: cannot fetch token: 401 Unauthorized.  Response: {"error":"access_denied","error_description":"Unauthorized"}',
  'connecting since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: connection error we are unaware of',
  // the following were identified as dev only errors -- leaving in case they need to be handled
  // 'connecting since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: failed to get access token: Post "https://aauth.idp.hcp.dev/oauth2/token": x509: “*.hcp.dev” certificate name does not match input',
  // 'connecting since 2022-09-13 14:45:40.666697 -0700 PDT m=+21.065498483; error: UNKNOWN',
];
let index = null;

export default function (server) {
  const handleResponse = (req, props) => {
    const xhr = req.passthrough();
    xhr.onreadystatechange = () => {
      if (xhr.readyState === 4 && xhr.status < 300) {
        // XMLHttpRequest response prop only has a getter -- redefine as writable and set value
        Object.defineProperty(xhr, 'response', {
          writable: true,
          value: JSON.stringify({
            ...JSON.parse(xhr.responseText),
            ...props,
          }),
        });
      }
    };
  };

  server.get('sys/seal-status', (schema, req) => {
    // return next status from statuses array
    if (index === null || index === statuses.length - 1) {
      index = 0;
    } else {
      index++;
    }
    return handleResponse(req, { hcp_link_status: statuses[index] });
  });
  // enterprise only feature initially
  server.get('sys/health', (schema, req) => handleResponse(req, { version: '1.12.0-dev1+ent' }));
}

// Package api is the official Golang API client for HashiCorp Vault. Vault's
// own internal core consumes this API.
//
// To get started, create a new API client:
//
//     client, err := api.NewClient(api.DefaultClient())
//     if err != nil {
//         log.Fatal(err)
//     }
//
// Chaining
//
// The library is designed to work in a "chained" fashion. For example, basic
// CRUD operations are available via the Logical() chain:
//
//     // Read a secret
//     secret, err := client.Logical().Read("secret/foo")
//
//     // Write a secret
//     secret, err := client.Logical().Write("secret/foo", map[string]interface{}{
//         "hello": "world"
//     })
//
//     // List secrets
//     secrets, err := client.Logical().List("secret/")
//
//     // Delete a secret
//     err := client.Logical().Delete("secret/foo")
//
// Similarly, the system backend is available via the Sys() chain:
//
//     // Check seal status
//     status, err := client.Sys().SealStatus()
//
// Context Support
//
// Most methods in this library support Go's standard context pattern by passing
// a context as the first argument to a method. All functions which support
// context are suffixed with "WithContext" like "MyFunctionWithContext".
//
//     // Read a secret with context
//     ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
//     defer cancel()
//     secret, err := client.Logical().ReadWithContext(ctx, "secret/foo")
//
//     // Check seal status with context
//     ctx, cancel := context.WithCancel(context.Background)
//     defer cancel()
//     status, err := client.Sys().SealStatusWithContext(ctx)
//
package api

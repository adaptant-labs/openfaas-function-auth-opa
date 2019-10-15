# openfaas-function-auth-opa

This repository provides an example of [Open Policy Agent](https://www.openpolicyagent.org/)\-backed authentication in OpenFaaS Serverless functions.

## Quick Start

To try it out, you will need to have an OPA server in your OpenFaaS stack. A version implementing this by default can
be found [here](https://github.com/adaptant-labs/faas/tree/opa-integration). Once this is up and running, fetch the
[golang-http-gomod](https://github.com/adaptant-labs/openfaas-golang-http-gomod-template) template and deploy as normal:

```
$ faas-cli template pull https://github.com/adaptant-labs/openfaas-golang-http-gomod-template.git
$ faas-cli up --skip-push
```
## Example Policy

A simple example rego policy is provided in order to get started. This policy
prohibits access by default, allowing access to the named function only for a
specified user:

```
package openfaas.authz

default allow = false

allow {
  input.function == "opa-auth"
  input.user == "alice"
}
```

## Function Invocation

Invocation of the function is prohibited by default by the example policy:

```
$ curl -X POST http://127.0.0.1:8080/function/opa-auth
Unauthorized.
```

Retrying the request with the permitted named user succeeds:

```
$ curl -H 'Authorization: alice' -X POST http://127.0.0.1:8080/function/opa-auth
Authorization OK
```

## Licensing

Released under the terms of the [MIT](LICENSE) license.

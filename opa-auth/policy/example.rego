package openfaas.authz

default allow = false

allow {
  input.function == "opa-auth"
  input.user == "alice"
}

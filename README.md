# nikolay's hashicups

progress: 5 / 11 steps

messing around with [terraform custom provider tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider)

## notes

- there is experimental support for [generating from OAS](https://developer.hashicorp.com/terraform/plugin/code-generation#openapi-provider-spec-generator). our API has long way to go before it fits their generator's expectation
- converting api client <-> provider's local model is booooooring already
- provider's assume CRUD is available and returns the resource for all resources
- fails might leave the service in weird state? needs to be tested like crazy against the spec
- hurrrrr

## Identity Schemas

Schemas for the HCC `x-rh-identity` header.


## Getting Started

There are multiple `x-rh-identity` header shapes based on the type of authentication
being used. This project is meant to be a source of truth for what the myriad
identities look like, for testing, enforcement and building clients.


## Usage

To run schema validation against the various types of identities, we have fixtures
for each. These fixtures are static, but should be recreated whenever a change to
the identity header is proposed or made.

Schemas can also be used to construct valid objects/structs for consuming the `x-rh-identity`.

### Validation
To run schema validation, simply run:
```shell
$ go run main.go
```

You will see output similar to:

```
*** Running schema validation against sample 3scale identities ***
Validating basic.json...
Error(s) found in validation:
- identity.internal.cross_access: Invalid type. Expected: number, given: boolean

*** Running schema validation against sample turnpike identities ***
Validating turnpike_saml.json...
Validating turnpike_x509.json...
```

### Schema consumption
You can use the schema to help validate and enforce valid identity objects to
work with in your project. Here are a few examples, though you can use whatever
you'd like:

Go:
```shell
# https://github.com/a-h/generate
$ go get -u github.com/a-h/generate/...
$ schema-generate ./turnpike/schema.json
```

Python:
```shell
# https://github.com/bcwaldon/warlock
$ pip install warlock
$ python3
>>> import warlock, json
>>> file = open("turnpike/schema.json")
>>> schema = json.loads(file.read())
>>> Identity = warlock.model_factory(schema)
>>> turnpike_identity = Identity(key="val")
```

Java:
```
https://github.com/joelittlejohn/jsonschema2pojo
```

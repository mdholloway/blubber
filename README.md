# Blubber

**Very experimental proof of concept.**

Blubber is a highly opinionated abstraction for container build configurations
and a command-line compiler which currently supports outputting multi-stage
Dockerfiles. It aims to provide a handful of declarative constructs that
accomplish build configuration in a more secure and determinate way than
running ad-hoc commands.

## Example configuration

```json
{
  "base": "debian:jessie",
  "apt": { "packages": ["libjpeg", "libyaml"] },
  "npm": { "install": true },
  "run": { "in": "/srv/service", "as": "runuser", "uid": 666, "gid": 666 },
  "variants": {
    "development": {
      "apt": { "packages": ["libjpeg-dev", "libyaml-dev"] }
    },
    "test": {
      "includes": ["development"],
      "apt": { "packages": ["chromium"] },
      "copiestree": true,
      "entrypoint": ["npm", "test"]
    },
    "production": {
      "base": "debian:jessie-slim",
      "npm": { "env": "production" },
      "artifacts": [{ "from": "test", "source": "/srv/service", "destination": "." }],
      "entrypoint": ["npm", "start"]
    }
  }
}
```

## Variants

Blubber supports a concept of composeable configuration variants for defining
slightly different container images while still maintaining a sufficient
degree of parity between them. For example, images for development and testing
may require some development and debugging packages which you wouldn't want in
production lest they contain vulnerabilities and somehow end up linked or
included in the application runtime.

Properties declared at the top level are shared among all variants unless
redefined, and one variant can include the properties of others. Some
properties, like `apt:packages` are combined when inherited or included.

In the example configuration, the `test` variant when expanded effectively
becomes:

```json
{
  "base": "debian:jessie",
  "apt": { "packages": ["libjpeg", "libyaml", "libjpeg-dev", "libyaml-dev", "chromium"] },
  "npm": { "install": true },
  "run": { "in": "/srv/service", "as": "runuser", "uid": 666, "gid": 666 },
  "copiestree": true,
  "entrypoint": ["npm", "test"]
}
```

## Artifacts

When trying to ensure optimally sized Docker images for production, there's a
common pattern that has emerged which is essentially to use one image for
building an application and copying the resulting build artifacts to another
much more optimized image, using the latter for production.

The Docker community has responded to this need by implementing
[multi-stage builds](https://github.com/moby/moby/pull/32063) and Blubber
makes use of this with its `artifacts` configuration property.

In the example configuration, the `production` variant declares artifacts to
be copied over from the result of building the `test` image.

## Usage

Running the `blubber` command will be produce `Dockerfile` output for the
given variant.

    blubber config.json variant

You can see the result of the example configuration by cloning this repo and
running (assuming you have go and your GOPATH set up properly):

    go build
    ./blubber blubber blubber.example.json development
    ./blubber blubber blubber.example.json test
    ./blubber blubber blubber.example.json production

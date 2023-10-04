# Contributing to vtProtobuf

## Workflow

For all contributors, we recommend the standard [GitHub flow](https://guides.github.com/introduction/flow/)
based on [forking and pull requests](https://guides.github.com/activities/forking/).

For significant changes, please [create an issue](https://github.com/planetscale/vtprotobuf/issues)
to let everyone know what you're planning to work on, and to track progress and design decisions.

## Development

### Protobuf version upgrade

1. Bump protobuf version in [./protobuf.sh](./protobuf.sh)) (PROTOBUF_VERSION variable).
1. Run `./protobuf.sh` to download and build protobuf.
1. Run `make genall` to regenerate proto files with a new compiler version, including well-known types.

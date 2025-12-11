[![REUSE status](https://api.reuse.software/badge/github.com/openmcp-project/service-provider-template)](https://api.reuse.software/info/github.com/openmcp-project/service-provider-template)

# service-provider-template

## About this project

A template for building @openmcp-project Service Providers

## Requirements and Setup

1. Create a new repository based on this template.
2. Execute the template to create a new `ServiceProvider`.
3. Test your `ServiceProvider`.

The template includes a basic code generation command that lets you create a `ServiceProvider` for your Go module, API kind and group.
You can also choose to add sample code to get a fully functional `ServiceProvider`.

For a complete usage overview with the default settings, run:

```shell
task gen-tmpl-help
```

Then execute the template, for example:

```shell
task gen-tmpl -- -module github.com/yourorg/yourrepo -kind YourKind -group yourgroup
```

Test your `ServiceProvider`:

```shell
task test
```

Note that `task gen-tmpl` is a one-time operation. To start over, delete your repository and recreate it from scratch.

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/openmcp-project/service-provider-template/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure

If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/openmcp-project/service-provider-template/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2025 SAP SE or an SAP affiliate company and service-provider-template contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/openmcp-project/service-provider-template).

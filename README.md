`pocli` is the command line interface for Veraison policy management service.
Policies are used to provide deployment-specific rules for amending and/or
augmenting attestation results generated for an attestation scheme (such as
PSA_IOT or CCA). Please see Veraison services [policy documentation] for more
details of how policies are used in Veraison.

[policy documentation]: https://github.com/veraison/services/tree/main/policy#readme


## Policy management overview

A policy is always associated with an [attestation scheme]. Veraison can keep
track of multiple policies for a single scheme, however at most one of them
can be "active" at a time. The active policy is the one that is actually
utilised by Veraison verification service.

When adding a new policy, `pocli` will automatically activate it (unless told
otherwise via a flag), however it is also possible to activate any previously
added policy (e.g. to roll back a bad update).

It is also possible to deactivate all policies for a scheme, ensuring that
attestation results are generated solely based on the scheme.

It is not possible to delete policies. This is to ensure that a policy that may
have been used to evaluate an attestation result is always available for
traceability.

For the same reason, it is not possible to update the rules on an existing
policy. Instead, a new policy containing updated rules must be added.

`pocli` also allows retrieving a specific policy by its UUID, retrieving the
active policy, or listing all policies for a scheme.

[attestation scheme]: https://github.com/veraison/services/tree/main/scheme


## Policy description

`pocli` represents policies as JSON objects. For example:

```json
{
    "uuid": "97609859-2567-11ee-a344-0242c0a82005",
    "ctime": "2023-07-18T12:35:34.661333946Z",
    "name": "default",
    "type": "opa",
    "rules": "package policy\n\nexecutables = APPROVED_RT\n",
    "active": false
}
```

A policy objects contains the following fields:

- **uuid** is the unique identifier of the policy. This is automatically
  generated when a policy is created.
- **ctime** is the creation timestamp.
- **name** is a descriptive label associated with the policy. This can be
  specified when creating a policy. If not specified, the default name
  `"default"` is used. Names do not need to be unique -- multiple policies can
  have the same name. This can be used to group policies that are different
  version of the same logical policy.
- **type** indicates the policy engine that will be used to evaluate the
  policy, and therefore the syntax of the policy rules. Currently, the only
  supported type is `"opa"`.
- **rules** are the policy rules that will be applied to the attestation
  result. The format of the rules depends on the policy type. See the Veraison
  services [OPA policy documentation] for details on how to write "opa"
  policy rules. You can also seen an example in
  `misc/example-PSA_IOT-policy.rego`.
- **active** is a boolean value indicating whether this policy is the current
  active policy for the scheme. At most one policy is active, however it is
  also possible for a scheme to have no active policies.

[OPA policy documentation]: https://github.com/veraison/services/blob/main/policy/README.opa.md


## Configuration

Configuration for `pocli` is specified using YAML markup language. `pocli`
reads configuration from `config.yaml` in the current working directory and/or
the user configuration directory (`$XDG_CONFIG_HOME/pocli/config.yaml` on
Linux). An alternate file may also be specified on the command line with
`-c`/`--config` option.

The following configuration options are supported:

- **host**: the host name of the Veraison management service. This can also be
  specified on the command line using `-H`/`--host` option.
- **port**: the port on which the Veraison management service is listening.
  This can also be specified on the command line using `-p`/`--port` option.
- **auth**: authentication method used by the Veraison management service. This
  can be either `passthrough`/`none`, `basic`, or `oauth2`. If not specified
  this defaults to `passthrough`. This can also be specified on the command
  line using `-a`/`--auth` flag.
- **username**: username for authenticating with the remote service. This is
  only used if `auth` is set to `basic` or `oauth2`. This can also be specified
  on the command line with `-U`/`--username`.
- **password**: password for authenticating with the remote service. This is
  only used if `auth` is set to `basic` or `oauth2`. This can also be specified
  on the command line with `-P`/`--password`.
- **client_id**: OAuth2 client ID. This is used only if `auth` is set to `oauth2`.
  This can also be specified on the command line using `-C`/`--client-id`.
- **client_secret**: OAuth2 client secret. This is used only if `auth` is set to
  `oauth2`. This can also be specified on the command line using
  `-S`/`--client-secret`.
- **token_url**: OAuth2 token endpoint URL. This is used only if `auth` is set
  to `oauth2`. This can also be specified on the command line using
  `-T`/`--token-url`.

See `misc/example-config.yaml` for an example configuration file.


## Usage examples

All policy management operations are performed with respect to a specific
attestation scheme. Thus most `pocli` commands take the scheme name as their
first argument. To see a list of valid scheme names for the service, you can
execute

    ./pocli schemes

A new policy can be created by passing the name of the attestation scheme to
which it is to be added and the path to the file containing the policy rules to
the `create` subcommand:

    ./pocli create PSA_IOT misc/example-PSA_IOT-policy.rego

This will also automatically activate the new policy (which can be suppressed
with `-d`/`--dont-activate` option). Optionally, the policy can be given a name
with `-n`/`--name`. If successful, `pocli` will output a JSON dump of the newly
created policy.

You can obtain the currently active policy using the `get` subcommand:

    ./pocli get PSA_IOT

This will write a JSON dump of the currently active policy to stdout. If there
is no active policy for the scheme, `pocli` will output an error to that
effect. `-o`/`--output` option can be used to write the JSON to a file instead
of stdout. Alternatively, `-w`/`--write-rules` can be used to write the rules
to the file specified with this option. That rules file can then be used with
`create` subcommand to create a new version of a policy.

The `get` subcommand can also be used to obtain a specific policy associated
with a scheme, rather than the active one specifying its UUID:

    ./pocli get PSA_IOT 94659974-2580-11ee-8fab-0242c0a84005

You can use `list` subcommand to get a list of all policies that have been
created for a scheme:

    ./pocli list PSA_IOT

As with the `get` subcommand, `-o`/`--output` option may be used to write the
resulting JSON array to a file.

A previously created policy may be activated with `activate` subcommand:

    ./pocli activate PSA_IOT 94659974-2580-11ee-8fab-0242c0a84005

As only one policy per scheme may be active at a time, this will also
deactivate the previously active policy, if there is one.

You can also deactivate all policies associated with a scheme, so that no
policy would be applied when the evidence for that scheme is appraised by the
verification service:

    ./pocli deactivate PSA_IOT

### Note on TLS

`-s`/`--tls` flag can be used to enable TLS, in which case system CA certs will
be used to validate the certificate sent by the server. It is possible to
disable server certificate validation with `-i`/`--insecure` flag (note that if
this flag is used, `-s` flag is implied and does not need to be explicitly specified.
specified). Alternatively, if the CA cert for the server is available but is
not installed in the system, it may be specified using `-E`/`--ca-cert` flag.

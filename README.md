# 1Password extension for Redact

This is a separate extension for [redact](https://github.com/julian7/redact), to store redact keys in 1Password.

It implements the standard extension interface:

- list
- get
- put

Parameters:

- key: `op://<vault>/<item>/[<section>/]<field>`

Usage:

- create a service account for the 1Password account, and grant read/write access to the vault
- set `OP_SERVICE_ACCOUNT_TOKEN` variable to the service account token
- configure redact to use the extension, and set key to point to a field in a 1Password vault item
- use `redact` extensions to retrieve and store secrets

Hint: if there's no section, and field is set to "Note", it uses the note message instead of a password field.

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
- set `OP_ACCOUNT_UUID` or `OP_SERVICE_ACCOUNT_TOKEN` variable to the service account token
- configure redact to use the extension, and set key to point to a field in a 1Password vault item
- use `redact` extensions to retrieve and store secrets

Hint: if there's no section, and field is set to "Note", it uses the note message instead of a password field.

## Authenticate with account UUID

There is a possibility to authenticate with the locally running 1Password app, just take a few steps in preparation:

First, allow integration with 1Password desktop app:

1. Run 1Password app
2. Open the drop-down menu at the top-left corner, and select "Developer"
3. Check "Integrate with other apps" under "Integrate with the 1Password SDKs"
4. Once you're at it, enable "Integrate with 1Password CLI" if not yet checked.

Second, get account UUID:

Run `op account list --format json`, and find your account's `account_uuid` field. Set it as an environment variable, like:

```
export OP_ACCOUNT_UUID=ABCDEFGH
```

## Authenticate with Service account token

Please refer to the official documentation: https://developer.1password.com/docs/service-accounts/

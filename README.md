# apidiff

APISIX declarative config diff and validation tool. It compares local YAML/JSON config with live APISIX Admin API data, outputs a plan-style diff, and performs semantic validation (reachability checks and plugin rules).

## Features

- Load local APISIX declarative config from YAML/JSON
- Fetch live config from APISIX Admin API
- Plan-style diff output with field-level changes
- Semantic validation
  - Upstream node reachability checks
  - Plugin rules (conflicts, required fields, require-one-of, deny fields, regex) with configurable rule sets

## Install

```
go build -o apidiff ./cmd/apidiff
```

## CLI Usage

### Plan

```
apidiff plan -f ./apisix.yaml --admin-url http://127.0.0.1:9180 --token <X-API-KEY>
```

Optional flags:
- `--skip-reachability` Skip upstream node reachability checks
- `--rules <file>` Use a custom plugin rules file (YAML/JSON)
- `--color` Enable/disable colored output (default: true)

### Validate

```
apidiff validate -f ./apisix.yaml
```

Optional flags:
- `--skip-reachability` Skip upstream node reachability checks
- `--rules <file>` Use a custom plugin rules file (YAML/JSON)

### Version

```
apidiff version
```

## Exit Codes

- `0` No diff and validation passed
- `1` Diff detected (plan output still printed)
- `2` Validation failed

## Example Output

```
Plan: 2 to add, 1 to change, 1 to delete.

+ route.demo_foo
~ upstream.u_1
  Nodes:
    ["10.0.0.1:8080"]: 1 -> 2
  Timeout:
    Connect: "1s" -> "2s"
- plugin_config.p1
```

## Plugin Rules Configuration

Rules can be loaded from a YAML/JSON file using `--rules`. If not provided, built-in defaults are used.

### Rules Schema

- `conflicts`: Plugins that cannot appear together
- `requires`: Fields that must be present when a plugin is enabled
- `require_one_of`: At least one field must be present when a plugin is enabled
- `deny_fields`: Fields that must not appear when a plugin is enabled
- `regex`: Regex constraints for plugin fields

### Example (rules.yaml)

```yaml
conflicts:
  - name: limit-req conflicts with limit-count
    scope: [route, service, plugin_config]
    plugins: [limit-req, limit-count]

requires:
  - name: key-auth requires key
    scope: [consumer]
    plugin: key-auth
    fields: [key]

require_one_of:
  - name: jwt-auth requires one of key/secret
    scope: [consumer]
    plugin: jwt-auth
    fields: [key, secret]

deny_fields:
  - name: limit-req forbids allow_degradation
    scope: [route, service]
    plugin: limit-req
    fields: [allow_degradation]

regex:
  - name: key-auth key format
    scope: [consumer]
    plugin: key-auth
    field: key
    pattern: "^[A-Za-z0-9_-]{8,64}$"
```

### Supported Scopes

- `route`
- `service`
- `plugin_config`
- `consumer`

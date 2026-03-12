# apidiff

APISIX 声明式配置差异对比与语义校验工具。对比本地 YAML/JSON 配置与 APISIX Admin API 的实时配置，输出 plan 风格 diff，并执行语义校验（上游可达性与插件规则）。

## 功能

- 读取本地 APISIX 声明式配置（YAML/JSON）
- 调用 APISIX Admin API 拉取实时配置
- Plan 风格 diff 输出（字段级变更）
- 语义校验
  - Upstream 节点可达性检查
  - 插件规则（冲突/必填字段），支持自定义规则集

## 安装

```
go build -o apidiff ./cmd/apidiff
```

## CLI 用法

### Plan

```
apidiff plan -f ./apisix.yaml --admin-url http://127.0.0.1:9180 --token <X-API-KEY>
```

可选参数：
- `--skip-reachability` 跳过 upstream 可达性检查
- `--rules <file>` 使用自定义插件规则文件（YAML/JSON）

### Validate

```
apidiff validate -f ./apisix.yaml
```

可选参数：
- `--skip-reachability` 跳过 upstream 可达性检查
- `--rules <file>` 使用自定义插件规则文件（YAML/JSON）

### Version

```
apidiff version
```

## 退出码

- `0` 无差异且校验通过
- `1` 存在差异（仍会输出 plan）
- `2` 校验失败

## 输出示例

```
Plan: 2 to add, 1 to change, 1 to delete.

+ route.demo_foo
~ upstream.u_1
  .Nodes["10.0.0.1:8080"]: 1 -> 2
  .Timeout.Connect: "1s" -> "2s"
- plugin_config.p1
```

## 插件规则配置

通过 `--rules` 指定 YAML/JSON 规则文件，不传则使用内置默认规则。

### 规则结构

- `conflicts`：互斥规则，插件不能同时出现
- `requires`：必填规则，插件启用时必须包含字段

### 示例 (rules.yaml)

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
```

### 支持的 scope

- `route`
- `service`
- `plugin_config`
- `consumer`

# Convert struct to map[string]any

## How to use

### Install this tool
```shell
 go install github.com/bytedance/nn/convToMap@latest
```
### Prepare your struct like example.go
[Check example](./unit_test/example.go)

### Generate code
```shell
 convToMap example.go
```

## Tips

1. Currently, code generation for structs that contain structs from different packages is not yet supported

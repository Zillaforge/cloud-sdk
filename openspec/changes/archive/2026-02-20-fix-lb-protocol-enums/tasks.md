## 1. Listener Protocol Enum

- [x] 1.1 Define `ListenerProtocol` typed string type in `models/vps/loadbalancers/listener.go`
- [x] 1.2 Add constants: `ListenerProtocolHTTP`, `ListenerProtocolHTTPS`, `ListenerProtocolTerminatedHTTPS`, `ListenerProtocolTCP`, `ListenerProtocolUDP`
- [x] 1.3 Add `Valid()` method to `ListenerProtocol`
- [x] 1.4 Change `Listener.Protocol` field type from `string` to `ListenerProtocol`
- [x] 1.5 Change `ListenerCreateRequest.Protocol` field type from `string` to `ListenerProtocol`
- [x] 1.6 Update `ListenerCreateRequest.Validate()` to use `ListenerProtocol.Valid()` and return "protocol must be one of" error for invalid values

## 2. Pool Protocol Enum

- [x] 2.1 Define `PoolProtocol` typed string type in `models/vps/loadbalancers/pool.go`
- [x] 2.2 Add constants: `PoolProtocolHTTP`, `PoolProtocolHTTPS`, `PoolProtocolTCP`, `PoolProtocolUDP`
- [x] 2.3 Add `Valid()` method to `PoolProtocol`
- [x] 2.4 Change `Pool.Protocol` field type from `string` to `PoolProtocol`
- [x] 2.5 Change `PoolCreateRequest.Protocol` field type from `string` to `PoolProtocol`
- [x] 2.6 Update `PoolCreateRequest.Validate()` to use `PoolProtocol.Valid()` and return "protocol must be one of" error for invalid values

## 3. Pool Method Enum (Consistent Usage)

- [x] 3.1 Change `Pool.Method` field type from `string` to `PoolMethod`
- [x] 3.2 Change `PoolCreateRequest.Method` field type from `string` to `PoolMethod`
- [x] 3.3 Change `PoolUpdateRequest.Method` field type from `string` to `PoolMethod`

## 4. Test Updates

- [x] 4.1 Update `models/vps/loadbalancers/listener_test.go` to use `ListenerProtocol` constants instead of string literals
- [x] 4.2 Update `models/vps/loadbalancers/pool_test.go` to use `PoolProtocol` and `PoolMethod` constants instead of string literals
- [x] 4.3 Update `modules/vps/loadbalancers/listeners_test.go` to use `ListenerProtocol` constants
- [x] 4.4 Update `modules/vps/loadbalancers/pools_test.go` to use `PoolProtocol` and `PoolMethod` constants

## 5. Verification

- [x] 5.1 Run `make build` and confirm zero compile errors
- [x] 5.2 Run `make test` and confirm all tests pass

package main

import "sync"

var keyStringValueMap = map[string]string{}
var ksvRWMutex = sync.RWMutex{}

var keyStringMapValueMap = map[string]map[string]string{}
var ksmvRWMutex = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hashSet,
	"HGET":    hashGet,
	"HGETALL": hashGetAll,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	ksvRWMutex.RLock()
	value, ok := keyStringValueMap[key]
	ksvRWMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	ksvRWMutex.Lock()
	keyStringValueMap[key] = value
	ksvRWMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hashGet(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	key1 := args[0].bulk
	key2 := args[1].bulk

	ksmvRWMutex.RLock()
	value, ok := keyStringMapValueMap[key1][key2]
	ksmvRWMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hashSet(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	key1 := args[0].bulk
	key2 := args[1].bulk
	value := args[2].bulk

	ksmvRWMutex.Lock()

	if _, ok := keyStringMapValueMap[key1]; !ok {
		keyStringMapValueMap[key1] = map[string]string{}
	}

	keyStringMapValueMap[key1][key2] = value

	ksmvRWMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hashGetAll(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	key1 := args[0].bulk

	ksmvRWMutex.RLock()
	valueMap, ok := keyStringMapValueMap[key1]
	ksmvRWMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	valueArray := []Value{}

	for key2, value := range valueMap {
		valueArray = append(valueArray, Value{typ: "bulk", bulk: key2}, Value{typ: "bulk", bulk: value})
	}

	return Value{typ: "array", array: valueArray}
}

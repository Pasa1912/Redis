package main

import (
	"strconv"
	"sync"
	"time"
)

type Node struct {
	value string
	ttl   time.Time
}

var keyStringValueMap = map[string]Node{}
var ksvRWMutex = sync.RWMutex{}

var keyStringMapValueMap = map[string]map[string]Node{}
var ksmvRWMutex = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hashSet,
	"HGET":    hashGet,
	"HGETALL": hashGetAll,
	"EXPIRE":  expire,
	"TTL":     getTimeToLive,
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
	node, ok := keyStringValueMap[key]
	ksvRWMutex.RUnlock()

	if !ok || (!node.ttl.IsZero() && time.Now().After(node.ttl)) {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: node.value}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	val := args[1].bulk

	ksvRWMutex.Lock()
	keyStringValueMap[key] = Node{value: val}
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
	node, ok := keyStringMapValueMap[key1][key2]
	ksmvRWMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: node.value}
}

func hashSet(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	key1 := args[0].bulk
	key2 := args[1].bulk
	val := args[2].bulk

	ksmvRWMutex.Lock()

	if _, ok := keyStringMapValueMap[key1]; !ok {
		keyStringMapValueMap[key1] = map[string]Node{}
	}

	keyStringMapValueMap[key1][key2] = Node{value: val}

	ksmvRWMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hashGetAll(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	key1 := args[0].bulk

	ksmvRWMutex.RLock()
	nodeMap, ok := keyStringMapValueMap[key1]
	ksmvRWMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	valueArray := []Value{}

	for key2, node := range nodeMap {
		valueArray = append(valueArray, Value{typ: "bulk", bulk: key2}, Value{typ: "bulk", bulk: node.value})
	}

	return Value{typ: "array", array: valueArray}
}

func expire(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'expire' command"}
	}

	key := args[0].bulk
	seconds := args[1].num

	ksvRWMutex.Lock()
	ttl := time.Now().Add(time.Duration(seconds) * time.Second)
	keyStringValueMap[key] = Node{value: keyStringValueMap[key].value, ttl: ttl}
	ksvRWMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func getTimeToLive(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'ttl' command"}
	}

	key := args[0].bulk

	ksvRWMutex.RLock()
	node, ok := keyStringValueMap[key]
	ksvRWMutex.RUnlock()

	if !ok || node.ttl.IsZero() {
		return Value{typ: "null"}
	}

	timeToLive := int(time.Since(node.ttl))

	if timeToLive < 0 {
		return Value{typ: "null"}
	}

	return Value{typ: "string", str: strconv.Itoa(timeToLive)}
}

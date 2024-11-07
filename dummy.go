package jitlua

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// #cgo CFLAGS: -I./lib/LuaJIT-2.1/src
// #cgo windows LDFLAGS: -L./lib/LuaJIT-2.1/src -lluajit-5.1.dll -lm
// #cgo linux LDFLAGS:  -L./lib/LuaJIT-2.1/src -lluajit -ldl -lm
// #include "glua.h"
import "C"

type dummy struct {
	key []byte
	val interface{}
}

var (
	dummyCache map[uintptr]map[uintptr]*dummy
	// dummyCache map[unsafe.Pointer]map[unsafe.Pointer]*dummy
	dummyRW sync.RWMutex
)

func init() {
	dummyCache = make(map[uintptr]map[uintptr]*dummy)
	// dummyCache = make(map[unsafe.Pointer]map[unsafe.Pointer]*dummy)
}

// lua dummy method
func pushDummy(vm *C.struct_lua_State, obj interface{}) unsafe.Pointer {
	vmKey := generateLuaStateId(vm)

	val := reflect.ValueOf(obj)
	var (
		realObj interface{}
		dummyId uintptr
	)

	switch val.Kind() {
	case reflect.Pointer:
		{
			realObj = val.Elem().Interface()
		}
	default:
		{
			realObj = obj
		}
	}

	dObj := &dummy{
		key: []byte(fmt.Sprintf("%p", &realObj)),
		val: obj,
	}

	// dummyId = uintptr(unsafe.Pointer(&(dObj.key[0])))
	pDumyId := unsafe.Pointer(&(dObj.key[0]))
	dummyId = uintptr(pDumyId)

	dummyRW.Lock()
	defer dummyRW.Unlock()

	target, ok := dummyCache[vmKey]
	if !ok {
		target = make(map[uintptr]*dummy)
		dummyCache[vmKey] = target
	}
	target[dummyId] = dObj

	// 返回 unsafe.Pointer 类型的 dummyId
	// return unsafe.Pointer(dummyId)
	return pDumyId
}

func findDummy(vm *C.struct_lua_State, ptr unsafe.Pointer) (interface{}, error) {
	vmKey := generateLuaStateId(vm)
	dummyId := uintptr(ptr) // 确保类型安全转换
	// dummyId := unsafe.Pointer(uintptr(ptr))  // 修改此行，确保类型正确转换

	dummyRW.RLock()
	defer dummyRW.RUnlock() // 添加 defer 解锁，确保在函数退出时解锁，增强代码稳定性

	target, ok := dummyCache[vmKey]
	if !ok {
		return nil, fmt.Errorf("%s", "Invalid VMKey")
	}
	dObj, ok := target[dummyId]
	if !ok {
		return nil, fmt.Errorf("%s", "Invalid DummyId")
	}
	return dObj.val, nil
}

func cleanDummy(vm *C.struct_lua_State) {
	vmKey := generateLuaStateId(vm)

	dummyRW.Lock()
	defer dummyRW.Unlock() // 添加 defer 解锁，确保在函数退出时解锁，增强代码稳定性
	delete(dummyCache, vmKey)
}

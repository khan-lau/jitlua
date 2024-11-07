# JitLua

## 介绍
基于LuaJit的golang包装，用于Lua脚本的执行和调用. 

本项目基于RyouZhang的[go-lua](https://github.com/RyouZhang/go-lua), 暂时只在其基础上进行了少量的修改

特点:
1. 执行速度快, 约比gopher-lua 快5-10倍左右
2. 使用cgo, buid时需要依赖[LuaJit 2.1](https://github.com/LuaJIT/LuaJIT)

## 编译
1. 下载LuaJit源码, 在其目录下执行`make -C ./`
2. 编译 jitLua demo, 在jitLua目录下执行`make`

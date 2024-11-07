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
3. 执行`cd dist && ./jitlua_example && cd -`

备注:
> 1. `jitlua_example` 执行需要依赖 `luajit.so` 或 `luajit.dll` 库, 请将 `libluajit-5.1.so` 复制到 `dist`目录下
> 2. `jitlua_example` 执行需要依赖 `script.lua` 文件, 请将 `script.lua` 从`cmd`目录复制到 `dist`目录下
> 3.  `make check`脚本中依赖`golangci-lint` 静态检查工具, 如要进行静态检查请先安装 `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
> 4. makefile 支持`windows msys2`环境, 支持`windows powershell`环境, 支持`windows cmd`环境, 支持`linux`环境, 其他环境未进行测试
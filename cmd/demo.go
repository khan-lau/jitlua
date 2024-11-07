package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	glua "github.com/khan-lau/jitlua"
)

func test_sum(ctx context.Context, args ...interface{}) (interface{}, error) {
	sum := 0
	for _, arg := range args {
		sum = sum + int(arg.(int64))
	}
	if sum%2 == 0 {
		return sum, nil
	} else {
		return nil, fmt.Errorf("bad sum")
	}
}

func json_decode(ctx context.Context, args ...interface{}) (interface{}, error) {
	raw := args[0].(string)

	var res map[string]interface{}
	err := json.Unmarshal([]byte(raw), &res)
	return []any{res, args[1]}, err
}

func get_header_field(ctx context.Context, args ...interface{}) (interface{}, error) {
	req := args[0].(*http.Request)
	key := args[1].(string)
	return req.Header.Get(key), nil
}

type A interface {
	Name() string
}

type AA struct {
}

func (a *AA) Name() string {
	return "hello world"
}

func main() {

	_ = glua.RegisterExternMethod("json_decode", json_decode)
	_ = glua.RegisterExternMethod("test_sum", test_sum)
	_ = glua.RegisterExternMethod("get_header_field", get_header_field)

	s := time.Now()
	res, err := glua.NewAction().WithScript(`
	function fib(n)
		if n == 0 then
			return 0
		elseif n == 1 then
			return 1
		end
		return fib(n-1) + fib(n-2)
	end
	`).WithEntrypoint("fib").AddParam(35).Execute(context.Background())
	// fmt.Println(time.Now().Sub(s))
	fmt.Println("time since:", time.Since(s))
	fmt.Println("inner fib result:", res, " error:", err)

	s = time.Now()
	res, err = glua.NewAction().WithScriptPath("script.lua").WithEntrypoint("fib").AddParam(35).Execute(context.Background())
	fmt.Println("time since:", time.Since(s))
	fmt.Println("fib result:", res, " error:", err)

	s = time.Now()
	res, err = glua.NewAction().WithScriptPath("script.lua").WithEntrypoint("fibt").AddParam(35).Execute(context.Background())
	fmt.Println("time since:", time.Since(s))
	fmt.Println("fibt result:", res, " error:", err)

	s = time.Now()
	res, err = glua.NewAction().WithScriptPath("script.lua").WithEntrypoint("test_args").AddParam([]interface{}{69, 56}).Execute(context.Background())
	fmt.Println("time since:", time.Since(s))
	fmt.Println("test_args golang test_sum result:", res, " error:", err)

	fmt.Println("======")
	s = time.Now()
	res, err = glua.NewAction().WithScriptPath("script.lua").WithEntrypoint("async_json_encode").Execute(context.Background())
	fmt.Println("time since:", time.Since(s))
	fmt.Println("async_json_encode result:", res, " error:", err)
	fmt.Println("======")

	s = time.Now()
	res, err = glua.NewAction().WithScriptPath("script.lua").WithEntrypoint("test_pull_table").AddParam(69).Execute(context.Background())
	fmt.Println("time since:", time.Since(s))
	fmt.Println("test_pull_table result:", res, " error:", err)

	req, _ := http.NewRequest("GET", "https://www.bing.com", nil)
	req.Header.Add("test", "3121232")

	//nolint:all
	var a A
	a = &AA{}

	res, err = glua.NewAction().WithScriptPath("script.lua").WithEntrypoint("test").AddParam(req, a).Execute(context.Background())
	fmt.Println("time since:", time.Since(s))
	fmt.Println("test golang get_header_field result:", res, " error:", err)
}

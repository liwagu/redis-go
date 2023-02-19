package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitee.com/wedone/redis_course/example"
)

func main() {
	defer example.RedisClient.Close()

	argsProg := os.Args
	var argsWithoutProg []string
	if len(argsProg) > 0 {
		argsWithoutProg = os.Args[1:]
		// fmt.Printf("\n==%v==\n", argsWithoutProg[0])
		fmt.Printf("输入参数:\n%s\n----------\n", strings.Join(argsWithoutProg, "\n"))
	}
	ctx := context.Background()
	runExample := argsWithoutProg[0]
	exampleParams := argsWithoutProg[1:]
	// funcValue := reflect.ValueOf(runExample)
	// paramList := make([]reflect.Value, 0, len(exampleParams)+1)
	// paramList = append(paramList, reflect.New(reflect.TypeOf(ctx)))
	// // paramList[0] = reflect.ValueOf(ctx)
	// for _, param := range exampleParams {
	// 	paramList = append(paramList, reflect.ValueOf(param))
	// }
	// // 反射调用函数
	// retList := funcValue.Call(paramList)
	// // 获取第一个返回值, 取整数值
	// fmt.Println(retList[0].Int())
	switch runExample {
	case "Ex01":
		example.Ex01(ctx, exampleParams)
	case "Ex02":
		example.Ex02(ctx)
	case "Ex03":
		example.Ex03(ctx)
	case "Ex04":
		example.Ex04(ctx)
	case "Ex05":
		example.Ex05(ctx, exampleParams)
	case "Ex06":
		example.Ex06(ctx, exampleParams)
	case "Ex06_2":
		fmt.Printf("%v\n", exampleParams)
		example.Ex06_2(ctx, exampleParams)
	case "Ex07":
		example.Ex07(ctx)
	}
}

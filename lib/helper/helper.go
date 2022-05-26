package helper

import (
	"os"
	"runtime"
	"strings"
)

//当前运行的方法名
func GetRunMethodName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	full_name := f.Name()
	names := strings.Split(full_name, ".")
	return names[len(names)-1]
}

func WriteFile(file_path string, content *string) (int, error) {
	f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	// 关闭文件
	defer f.Close()
	// 字符串写入
	return f.WriteString(*content)
}

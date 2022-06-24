package pattern_utils

// 保护方式允许一个函数
func ProtectRun(entry func(), handle func(err interface{})) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err != nil {
			handle(err)
		}
	}()
	entry()
}
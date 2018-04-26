package python 

func OptimiseCode(name string) string {
	if name == "text_itoa" {
		return `	runtime.Stack.pop()
	runtime.Lists.append(bytearray(str(runtime.Stack.pop()), "utf8"))
	return
`}
	
	return ""
}

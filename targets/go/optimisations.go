package golang

func OptimiseCode(name string) string {
	if name == "text_itoa" {
		return `
	var b = runtime.Stack[len(runtime.Stack)-1]
	var value = runtime.Stack[len(runtime.Stack)-2]
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]

	var str string
	if value.Bits() == nil {
		str = strconv.FormatInt(value.Small, int(b.Small)) 
	} else {
		str = value.Text(int(b.Small))
	}

	var result List

	for _, char := range []byte(str) {
		result.Put(Int{Small:int64(char)})
	}

	runtime.Lists = append(runtime.Lists, &result)
	return
`}

	return ""
}

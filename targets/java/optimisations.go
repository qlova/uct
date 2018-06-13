package java
func OptimiseCode(name string) string {
	if name == "text_itoa" {
		return `
	Int base = runtime.Stack[runtime.StackPointer];
	Int value = runtime.Stack[runtime.StackPointer-1];
	runtime.StackPointer -= 2;
	
	String str;
	if (value.Big == null) {
		str = String.valueOf(value.Small);
	} else {
		str = String.valueOf(value.Big);
	}
	
	List result = new List();

	for (int i = 0; i < str.length(); i++) {
		result.put(new Int(str.charAt(i)));
	}

	runtime.pushList(result);
	return;
`}

	return ""
}

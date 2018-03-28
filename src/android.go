package uct

import "flag"

//This is the Android compiler for uct.
var Android bool

var AndroidAssembly = Assemblable{}

func init() {
	flag.BoolVar(&Android, "android", false, "Target Android")
	
	for k, v := range JavaAssembly {
		AndroidAssembly[k] = v
	}
	
	for _, word := range JavaReserved {
		AndroidAssembly[word] = Reserved()
	}
	
	AndroidAssembly["ANDROID"] = AndroidAssembly["JAVA"]
	delete(AndroidAssembly, "JAVA")
	
	AndroidAssembly["HEADER"] = Instruction{
		Data:   " ",
		Indent: 1,
		Args:   1,
	}

	RegisterAssembler(AndroidAssembly , &Android, "android", "//")
}

package uct

import "flag"

//This is the QML compiler for uct.
var Qml bool

var QmlAssembly = Assemblable{}

func init() {
	flag.BoolVar(&Qml, "qml", false, "Target QML")
	
	for k, v := range JavascriptAssembly {
		QmlAssembly[k] = v
	}
	
	QmlAssembly["QML"] = QmlAssembly["JAVASCRIPT"]
	delete(QmlAssembly, "JAVASCRIPT")

	RegisterAssembler(QmlAssembly , &Qml, "qml", "//")
}

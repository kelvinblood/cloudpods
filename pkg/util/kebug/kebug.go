package kebug

import (
	"runtime"

	"yunion.io/x/log"
)

func Info(myTag string) {
	// pc,file,line,_ := runtime.Caller(2)
	pc0, file0, line0, _ := runtime.Caller(3)
	f0 := runtime.FuncForPC(pc0)
	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	pc2, file2, line2, _ := runtime.Caller(1)
	f2 := runtime.FuncForPC(pc2)
	log.Infoln("============"+myTag, file0, line0, f0.Name()+"()\n\t\t\t\t"+"->"+file, line, f.Name()+"()\n\t\t\t\t"+"->"+file2, line2, f2.Name()+"()")
}

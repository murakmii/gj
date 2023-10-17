package main

import (
	"flag"
	"fmt"
	"github.com/murakmii/gj"
	"github.com/murakmii/gj/vm"
	_ "github.com/murakmii/gj/vm/native"
	"os"
	"strings"
	"unicode/utf16"
)

var (
	configPath string
	mainClass  string
	print      bool
)

func init() {
	flag.StringVar(&configPath, "config", "", "path of configuration file")
	flag.StringVar(&mainClass, "main", "", "main class name")
	flag.BoolVar(&print, "print", false, "print disassembled class file")
}

func main() {
	flag.Parse()
	if len(configPath) == 0 || len(mainClass) == 0 {
		flag.Usage()
		return
	}
	mainClass = strings.ReplaceAll(mainClass, ".", "/") + ".class"

	config, err := readConfig()
	if err != nil {
		fmt.Printf("failed to read config: %s", err)
		return
	}

	classPaths, err := gj.InitClassPaths(config.ClassPath)
	if err != nil {
		fmt.Printf("failed to init class path: %s", err)
		return
	}
	defer func() {
		for _, classPath := range classPaths {
			classPath.Close()
		}
	}()

	if print {
		execPrint(classPaths)
	} else {
		execVM(config)
	}
}

func readConfig() (*gj.Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return gj.ReadConfig(f)
}

func execPrint(classPaths []gj.ClassPath) {
	for _, classPath := range classPaths {
		classFile, err := classPath.SearchClass(mainClass)
		if err != nil {
			fmt.Printf("failed to search class: %s", err)
			return
		}

		if classFile == nil {
			continue
		}

		fmt.Printf(classFile.String())
		return
	}

	fmt.Println("class not found")
}

func execVM(config *gj.Config) {
	vmInstance, err := vm.InitVM(config)
	if err != nil {
		panic(err)
	}

	className := "java/lang/String"
	javaLangString, state, err := vmInstance.FindInitializedClass(&className, vm.NewThread(vmInstance))
	if err != nil {
		panic(err)
	}
	if state == vm.FailedInitialization {
		panic("string class initialization failed")
	}

	className = "java/lang/Class"
	_, state, err = vmInstance.FindInitializedClass(&className, vm.NewThread(vmInstance))
	if err != nil {
		panic(err)
	}
	if state == vm.FailedInitialization {
		panic("class class initialization failed")
	}

	fmt.Println("succeeded string class initialization")

	thread := vm.NewThread(vmInstance)

	strA := "Hello, ワ"
	strB := "ールド!"

	jsA, _ := vmInstance.JavaString(thread, &strA)
	jsB, _ := vmInstance.JavaString(thread, &strB)

	class, method := javaLangString.ResolveMethod("concat", "(Ljava/lang/String;)Ljava/lang/String;")
	frame := vm.NewFrame(class, method).SetLocals([]interface{}{jsA, jsB})

	if _, err = thread.Execute(frame); err != nil {
		panic(err)
	}

	fClass := "java/lang/String"
	fName := "value"
	charArray := thread.GetResult().(*vm.Instance).GetField(&fClass, &fName).(*vm.Array)

	u16 := make([]uint16, charArray.Length())
	for i := 0; i < charArray.Length(); i++ {
		u16[i] = uint16(charArray.Get(i).(int))
	}

	fmt.Printf("**** concat = %s ****\n", string(utf16.Decode(u16)))
}

package main

import (
	"flag"
	"fmt"
	"github.com/murakmii/gj"
	"os"
	"strings"
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
		// TODO: execute jvm
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

		fmt.Printf("Constant pool size: %d\n", classFile.ConstantPool().Size())

		break
	}

	fmt.Println("class not found")
}

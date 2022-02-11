package cmd

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type ArgsNumber int

const (
	VERSION                       = "1.1.0"
	loginArgsNumber    ArgsNumber = 1
	copyFileArgsNumber ArgsNumber = 2
)

var (
	ver *bool
	cp  *bool
)

func init() {
	// 参数定义
	ver = flag.Bool("v", false, "Show Version")
	cp = flag.Bool("cp", false, "tryssh -cp [Src File Path] [Dest File Path]\n"+
		"Eg: \n\tdownload: tryssh -cp 192.168.1.1:/root/test.txt ./\n"+
		"\tupload: tryssh -cp ./test.txt 192.168.1.1:/root/\n")

	flag.Usage = func() {
		fmt.Println("Usage: \n\ttryssh [IpAddress]\tTry to login [IpAddress]" +
			"\n\ttryssh -cp [Src File Path] [Dest File Path]\tTry to upload/download file\n" +
			"Options:")
		flag.PrintDefaults()
	}
}

// FlagsParse 数据结构体，命令行参数信息
type FlagsParse struct {
	Version  bool
	CopyFile bool
	Args     []string
}

// NewFlagsParse 通过参数校验后创建FlagsParse对象
func NewFlagsParse() FlagsParse {
	flag.Parse()
	args := flag.Args()

	if *ver {
		printVersionAndExit()
	}
	argsNumberCheck(args)

	return FlagsParse{
		Version:  *ver,
		CopyFile: *cp,
		Args:     args,
	}
}

// printVersionAndExit 打印版本号后退出
func printVersionAndExit() {
	log.Infof("Tryssh Version: %s", VERSION)
	os.Exit(0)
}

// argsNumberCheck 参数个数检查
func argsNumberCheck(args []string) {
	argsLen := ArgsNumber(len(args))
	if *cp {
		if argsLen == copyFileArgsNumber {
			return
		} else {
			log.Errorf("Wrong number of parameters. There can only be %d, got %d\n\n",
				copyFileArgsNumber, argsLen)
			flag.Usage()
			os.Exit(1)
		}
	} else {
		if argsLen == loginArgsNumber {
			return
		} else {
			log.Errorf("Wrong number of parameters. There can only be %d, got %d\n\n",
				loginArgsNumber, argsLen)
			flag.Usage()
			os.Exit(1)
		}
	}
}

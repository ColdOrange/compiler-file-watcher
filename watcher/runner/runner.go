package runner

type Runner interface {
	Run() error
}

var (
	_ Runner = &ProtocolRunner{}
	_ Runner = &OssDescRunner{}
	_ Runner = &ExcelRunner{}
)

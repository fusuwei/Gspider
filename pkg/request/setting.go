package request

type Setting struct {
	Host           string
	Timeout        int
	Verify         bool // 是否验证证书
	AllowRedirects bool // 禁止跳转
	Header         map[string]string
}

func (s *Setting) AutoHeaders() {

}

package request

type Setting struct {
	Timeout        int
	Verify         bool // 是否验证证书
	AllowRedirects bool // 禁止跳转
}

func (s *Setting) AutoHeaders() {

}

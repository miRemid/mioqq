package http

const (
	// Info Logger信息
	Info = iota
	// Error 错误信息
	Error
	// Warn 警告信息
	Warn
)

// SendLog 日志处理
func (server *Server) SendLog(status int, format string, a ...interface{}) {
	if !server.Logger {
		return
	}
	switch status {
	case Info:
		server.logger.Printf("[I] "+format, a...)
		break
	case Error:
		server.logger.Printf("[E] "+format, a...)
		break
	case Warn:
		server.logger.Printf("[W] "+format, a...)
		break
	default:
		break
	}
}

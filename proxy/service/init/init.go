package init

import (
	"github.com/imsilence/kaleidoscope/proxy/service"
	"github.com/imsilence/kaleidoscope/proxy/service/provider/ssh"
)

func init() {
	service.DefaultManager.Register(new(ssh.Ssh))
}

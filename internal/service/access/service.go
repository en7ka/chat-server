package access

import (
	repinf "github.com/en7ka/chat-server/internal/repository/repointerface"
	servinf "github.com/en7ka/chat-server/internal/service/servinterface"
)

type ServiceAcc struct {
	accessRepository servinf.Access
}

func NewServiceAcc(accessRepository repinf.Access) ServiceAcc {
	return ServiceAcc{accessRepository: accessRepository}
}

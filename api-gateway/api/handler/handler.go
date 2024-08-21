package handler
import(
	"gitlab.com/bahodirova/api-gateway/client"
)
type HandlerStruct struct {
	Clients client.Clients
}

func NewHandlerStruct() *HandlerStruct {
	return &HandlerStruct{
		Clients: *client.NewClients(),
	}
}
package admin

import (
	"github.com/beatrice950201/araneid/extend/service"
	"github.com/gorilla/websocket"
	"net/http"
)

type Admin struct{ Main }

// @router /admin/socket [get]
func (c *Admin) Socket() {
	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok == false {
		service.SocketInstanceGet().Join(c.UserInfo.Id, ws)
		defer service.SocketInstanceGet().Leave(c.UserInfo.Id)
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				c.StopRun()
			}
		}
	} else {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
	}
}

// @router / [get]
func (c *Admin) Index() {
	c.Data["disk"] = c.adminService.DiskDashboard()
	c.Data["processing"] = c.adminService.DashboardProcessing()
}

package index

import "github.com/beatrice950201/araneid/controllers"

type Index struct{ Main }

// @router / [get]
func (c *Index) Index() {
	c.Succeed(&controllers.ResultJson{
		Message: "success!!!",
		Data:    c.DomainCache,
	})
}

package global

import "autoLogin/models"

var Config *models.Config

var Status struct {
	Output   bool
	Guardian bool
}

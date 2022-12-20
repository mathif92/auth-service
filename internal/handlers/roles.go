package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/services"
)

type Roles struct {
	roleService *services.Roles
}

func (r *Roles) CreateRole(w http.ResponseWriter, req *http.Request) {

}

func (r *Roles) DeleteRole(w http.ResponseWriter, req *http.Request) {

}

func (r *Roles) AddActionToRole(w http.ResponseWriter, req *http.Request) {

}

func (r *Roles) AssignRole(w http.ResponseWriter, req *http.Request) {

}

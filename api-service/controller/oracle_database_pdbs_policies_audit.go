// Copyright (c) 2024 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package controller

import (
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) GetOraclePdbsPoliciesAudit(w http.ResponseWriter, r *http.Request) {
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]
	pdbname := mux.Vars(r)["pdbname"]

	resp, err := ctrl.Service.GetOracleDatabasePdbPoliciesAuditFlag(hostname, dbname, pdbname)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

func (ctrl *APIController) ListOraclePdbPoliciesAudit(w http.ResponseWriter, r *http.Request) {
	resp, err := ctrl.Service.ListOracleDatabasePdbPoliciesAudit()
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

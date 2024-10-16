// Copyright (c) 2022 Sorint.lab S.p.A.
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

package job

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"

	alert_service_client "github.com/ercole-io/ercole/v2/alert-service/client"
	"github.com/ercole-io/ercole/v2/data-service/database"
)

// FreshnessCheckJob is the job used to check the freshness of the current hosts
type FreshnessCheckJob struct {
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// AlertSvcClient
	AlertSvcClient alert_service_client.AlertSvcClientInterface
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Log contains logger formatted
	Log logger.Logger
	// NewObjectID return a new ObjectID
	NewObjectID func() primitive.ObjectID
}

// Run throws NO_DATA alert for each hosts that haven't sent a hostdata withing the host.Period (hours)
func (job *FreshnessCheckJob) Run() {
	if err := job.Database.DeleteAllNoDataAlerts(); err != nil {
		job.Log.Error(err)
		return
	}

	hosts, err := job.Database.GetActiveHostdata()

	if err != nil {
		job.Log.Error(err)
		return
	}

	for _, host := range hosts {
		var period time.Duration

		if host.Period <= 0 {
			period = 24
		} else {
			period = time.Duration(host.Period)
		}

		isOld, err := job.Database.FindOldCurrentHostdata(host.Hostname, job.TimeNow().Add(-(period)*time.Hour))
		if err != nil {
			job.Log.Error(err)
			continue
		}

		if isOld {
			elapsed := job.TimeNow().Sub(host.CreatedAt)
			elapsedDays := int(elapsed.Truncate(time.Hour*24).Hours() / 24)

			alert := model.Alert{
				ID:                      job.NewObjectID(),
				AlertAffectedTechnology: nil,
				AlertCategory:           model.AlertCategoryAgent,
				AlertCode:               model.AlertCodeNoData,
				AlertSeverity:           model.AlertSeverityCritical,
				AlertStatus:             model.AlertStatusNew,
				Date:                    job.TimeNow(),
				Description:             fmt.Sprintf("No data received from the host %s in the last %d day(s)", host.Hostname, elapsedDays),
				OtherInfo: map[string]interface{}{
					"hostname": host.Hostname,
				},
			}

			if job.Config.AlertService.Emailer.AlertType.NoData {
				errAlert := job.AlertSvcClient.ThrowNewAlert(alert)
				if errAlert != nil {
					job.Log.Error(errAlert)
					continue
				}
			}
		}
	}
}

// Copyright (c) 2019 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"time"

	"github.com/amreo/ercole-services/api-service/database"
	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/config"
)

// APIServiceInterface is a interface that wrap methods used to querying data
type APIServiceInterface interface {
	// Init initialize the service
	Init()
	// SearchCurrentHosts search current hosts
	SearchCurrentHosts(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetCurrentHost return the current host specified in the hostname param
	GetCurrentHost(hostname string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// SearchAlerts search alerts
	SearchAlerts(search string, sortBy string, sortDesc bool, page int, pageSize int, severity string, status string, from time.Time, to time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentClusters search current clusters
	SearchCurrentClusters(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentAddms search current addm
	SearchCurrentAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentSegmentAdvisors search current segment advisors
	SearchCurrentSegmentAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentPatchAdvisors search current patch advisors
	SearchCurrentPatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentDatabases search current databases
	SearchCurrentDatabases(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentExadata search current exadata
	SearchCurrentExadata(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// ListCurrentLicenses list current licenses
	ListCurrentLicenses(full bool, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)

	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOperatingSystemStats return a array containing the number of hosts per operating system
	GetOperatingSystemStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per type
	GetTypeStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
	GetDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseVersionStats return a array containing the number of databases per version
	GetDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopReclaimableDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
	GetTopReclaimableDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabasePatchStatusStats return a array containing the number of databases per patch status
	GetDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopWorkloadDatabaseStats return a array containing top databases by workload
	GetTopWorkloadDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
	GetDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseRACStatusStats return a array containing the number of databases per RAC status
	GetDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
	GetDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTotalDatabaseWorkStats return the total work of databases
	GetTotalDatabaseWorkStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalDatabaseMemorySizeStats return the total of memory size of databases
	GetTotalDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalDatabaseDatafileSizeStats return the total size of datafiles of databases
	GetTotalDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalDatabaseSegmentSizeStats return the total size of segments of databases
	GetTotalDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
	GetDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetTotalExadataMemorySizeStats return the total size of memory of exadata
	GetTotalExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetTotalExadataCPUStats return the total cpu of exadata
	GetTotalExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetAvegageExadataStorageUsageStats return the average usage of cell disks of exadata
	GetAvegageExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface)
	// GetExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
	GetExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetExadataPatchStatusStats return a array containing the number of exadata per patch status
	GetExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)

	// SetLicenseCount set the count of a certain license
	SetLicenseCount(name string, count int) utils.AdvancedErrorInterface
}

// APIService is the concrete implementation of APIServiceInterface.
type APIService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Version of the saved data
	Version string
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// Init initializes the service and database
func (as *APIService) Init() {
}

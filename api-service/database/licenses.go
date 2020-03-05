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

package database

import (
	"context"
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// ListLicenses list licenses
func (md *MongoDatabase) ListLicenses(full bool, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)
	//Find the informations

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APLookupPipeline("hosts", bson.M{
				"LicenseName": "$_id",
			}, "Used", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"Hostname": 1,
					"Databases": mu.APOReduce(
						mu.APOFilter(
							mu.APOMap("$Extra.Databases", "db", bson.M{
								"Name": "$$db.Name",
								"Count": mu.APOLet(
									bson.M{
										"val": mu.APOArrayElemAt(mu.APOFilter("$$db.Licenses", "lic", mu.APOEqual("$$lic.Name", "$$LicenseName")), 0),
									},
									"$$val.Count",
								),
							}),
							"db",
							mu.APOGreater("$$db.Count", 0),
						),
						bson.M{"Count": 0, "DBs": bson.A{}},
						bson.M{
							"Count": mu.APOMax("$$value.Count", "$$this.Count"),
							"DBs": bson.M{
								"$concatArrays": bson.A{
									"$$value.DBs",
									bson.A{"$$this.Name"},
								},
							},
						},
					),
				}),
				mu.APMatch(bson.M{
					"Databases.Count": bson.M{
						"$gt": 0,
					},
				}),
				mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname"}, "VM", mu.MAPipeline(
					FilterByOldnessSteps(olderThan),
					mu.APUnwind("$Extra.Clusters"),
					mu.APReplaceWith("$Extra.Clusters"),
					mu.APUnwind("$VMs"),
					mu.APReplaceWith("$VMs"),
					mu.APMatch(mu.QOExpr(mu.APOEqual("$Hostname", "$$hn"))),
					mu.APLimit(1),
				)),
				mu.APSet(bson.M{
					"VM": mu.APOArrayElemAt("$VM", 0),
				}),
				mu.APAddFields(bson.M{
					"ClusterName":  mu.APOIfNull("$VM.ClusterName", nil),
					"PhysicalHost": mu.APOIfNull("$VM.PhysicalHost", nil),
				}),
				mu.APUnset("VM"),
				mu.APGroup(mu.BsonOptionalExtension(full, bson.M{
					"_id": mu.APOCond(
						"$ClusterName",
						mu.APOConcat("cluster_§$#$§_", "$ClusterName"),
						mu.APOConcat("hostname_§$#$§_", "$Hostname"),
					),
					"License":    mu.APOMaxAggr("$Databases.Count"),
					"ClusterCpu": mu.APOMaxAggr("$ClusterCpu"),
				}, bson.M{
					"Hosts": mu.APOPush(bson.M{
						"Hostname":  "$Hostname",
						"Databases": "$Databases.DBs",
					}),
				})),
				mu.APSet(bson.M{
					"License": mu.APOCond(
						"$ClusterCpu",
						mu.APODivide("$ClusterCpu", 2),
						"$License",
					),
				}),
				mu.APGroup(mu.BsonOptionalExtension(full, bson.M{
					"_id":   0,
					"Value": mu.APOSum("$License"),
				}, bson.M{
					"Hosts": mu.APOPush("$Hosts"),
				})),
				mu.APOptionalStage(full, mu.MAPipeline(
					mu.APUnwind("$Hosts"),
					mu.APUnwind("$Hosts"),
					mu.APGroup(bson.M{
						"_id":   0,
						"Value": mu.APOMaxAggr("$Value"),
						"Hosts": mu.APOPush("$Hosts"),
					}),
				)),
			)),
			mu.APSet(bson.M{
				"Used": mu.APOArrayElemAt("$Used", 0),
			}),
			mu.APOptionalStage(full, mu.APSet(bson.M{
				"Hosts": mu.APOIfNull("$Used.Hosts", bson.A{}),
			})),
			mu.APSet(bson.M{
				"Used": mu.APOIfNull(mu.APOCeil("$Used.Value"), 0),
			}),
			mu.APSet(bson.M{
				"Compliance": mu.APOGreaterOrEqual("$Count", "$Used"),
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// SetLicenseCount set the count of a certain license
func (md *MongoDatabase) SetLicenseCount(name string, count int) utils.AdvancedErrorInterface {
	//Find the informations
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").UpdateOne(context.TODO(), bson.M{
		"_id": name,
	}, mu.UOSet(bson.M{
		"Count": count,
	}))
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Check the existance of the result
	if res.MatchedCount == 0 {
		return utils.AerrLicenseNotFound
	} else {
		return nil
	}
}

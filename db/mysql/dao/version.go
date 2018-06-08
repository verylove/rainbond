// Copyright (C) 2014-2018 Goodrain Co., Ltd.
// RAINBOND, Application Management Platform

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package dao

import (
	"github.com/goodrain/rainbond/db/model"

	"fmt"

	"github.com/jinzhu/gorm"
	"time"
)

//DeleteVersionByEventID DeleteVersionByEventID
func (c *VersionInfoDaoImpl) DeleteVersionByEventID(eventID string) error {
	version := &model.VersionInfo{
		EventID: eventID,
	}
	if err := c.DB.Where("event_id = ? ", eventID).Delete(version).Error; err != nil {
		return err
	}
	return nil
}

//DeleteVersionByServiceID DeleteVersionByServiceID
func (c *VersionInfoDaoImpl) DeleteVersionByServiceID(serviceID string) error {
	var version model.VersionInfo
	if err := c.DB.Where("service_id = ? ", serviceID).Delete(&version).Error; err != nil {
		return err
	}
	return nil
}

//AddModel AddModel
func (c *VersionInfoDaoImpl) AddModel(mo model.Interface) error {
	result := mo.(*model.VersionInfo)
	var oldResult model.VersionInfo
	if ok := c.DB.Where("event_id=?", result.EventID).Find(&oldResult).RecordNotFound(); ok {
		if err := c.DB.Create(result).Error; err != nil {
			return err
		}
	} else {
		fmt.Errorf("version is exist")
		return nil
	}
	return nil
}

//UpdateModel UpdateModel
func (c *VersionInfoDaoImpl) UpdateModel(mo model.Interface) error {
	result := mo.(*model.VersionInfo)
	if err := c.DB.Save(result).Error; err != nil {
		return err
	}
	return nil
}

//EventLogMessageDaoImpl EventLogMessageDaoImpl
type VersionInfoDaoImpl struct {
	DB *gorm.DB
}

//GetEventLogMessages get event log message
func (c *VersionInfoDaoImpl) GetVersionByEventID(eventID string) (*model.VersionInfo, error) {
	var result model.VersionInfo
	if err := c.DB.Where("event_id=?", eventID).Find(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//return messageRaw, nil
		}
		return nil, err
	}
	return &result, nil
}

//GetEventLogMessages get event log message
func (c *VersionInfoDaoImpl) GetVersionByDeployVersion(version, serviceID string) (*model.VersionInfo, error) {
	var result model.VersionInfo
	if err := c.DB.Where("build_version =? and service_id = ?", version, serviceID).Find(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &result, nil
}

//GetEventLogMessages get event log message
func (c *VersionInfoDaoImpl) GetVersionByServiceID(serviceID string) ([]*model.VersionInfo, error) {
	var result []*model.VersionInfo
	if err := c.DB.Where("service_id=?", serviceID).Find(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//return messageRaw, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *VersionInfoDaoImpl) GetVersionInfo(timePoint time.Time, deliveredType string) ([]*model.VersionInfo, error) {
	var result []*model.VersionInfo

	if err := c.DB.Where("create_time  < ? AND delivered_type = ?", timePoint, deliveredType).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil

}

func (c *VersionInfoDaoImpl) DeleteVersionInfo(obj *model.VersionInfo) error {
	if err := c.DB.Delete(obj).Error; err != nil {
		return nil
	} else {
		return err
	}
}

func (c *VersionInfoDaoImpl) DeleteFailureVersionInfo(timePoint time.Time, status string) error {
	if err := c.DB.Where("create_time  < ? AND final_status = ?", timePoint, status).Delete(&model.VersionInfo{}).Error; err != nil {
		return err
	}
	return nil
}


func (c *VersionInfoDaoImpl) SearchVersionInfo() ([]*model.VersionInfo, error) {
	var result []*model.VersionInfo
	if err := c.DB.Table("version_info").Select("service_id as ServiceId, count(ID) as total").Group("service_id").Having("count(ID) > ?", 5).Scan(&result).Error; err != nil {
		return nil,err
	}else {
		return result, nil

	}

}

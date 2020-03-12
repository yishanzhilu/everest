package workspace

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func registerRecordRoutes(r *gin.RouterGroup) {
	r.POST("/records", postRecord)
	r.GET("/records", getRecordList)
	r.GET("/record/:id", getRecord)
	r.DELETE("/record/:id", deleteRecord)
}

type postRecordBody struct {
	Content   string `json:"content" binding:"required,max=255"`
	Review    string `json:"review" binding:"max=255"`
	Mood      string `json:"mood" binding:"max=10"`
	Minutes   uint16 `json:"minutes" binding:"max=480"`
	GoalID    uint64 `json:"goalID"`
	MissionID uint64 `json:"missionID"`
}

func postRecord(c *gin.Context) {
	var body postRecordBody
	var err error
	if err = c.BindJSON(&body); err != nil {
		return
	}

	uid := c.MustGet(common.ContextUserID).(uint64)
	record := models.RecordModel{
		Content: body.Content,
		Review:  body.Review,
		Minutes: body.Minutes,
		Mood:    body.Mood,
		UserID:  uid,
	}

	if body.MissionID > 0 {
		mission, err := getMissionHelper(uid, body.MissionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"missionID": err.Error()})
			return
		}
		record.Mission.ID = mission.ID
		record.Mission.Title = mission.Title
		record.Goal.ID = mission.GoalID
		record.Goal.Title = mission.Goal.Title
	} else if body.GoalID > 0 {
		goal, err := getGoalHelper(uid, body.GoalID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goalID": err.Error()})
			return
		}
		record.Goal.ID = goal.ID
		record.Goal.Title = goal.Title
	}

	err = common.MySQLClient.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		if record.Minutes > 0 {
			if record.Goal.ID > 0 {
				goal := models.GoalModel{}
				goal.ID = record.Goal.ID
				if err := tx.Model(&goal).
					Update("minutes", gorm.Expr("minutes + ?", record.Minutes)).
					Error; err != nil {
					return err
				}
			}
			if record.Mission.ID > 0 {
				mission := models.MissionModel{}
				mission.ID = record.Mission.ID
				if err := tx.Model(&mission).
					Update("minutes", gorm.Expr("minutes + ?", record.Minutes)).
					Error; err != nil {
					return err
				}
			}
			user := models.UserModel{}
			user.ID = uid
			if err := tx.Model(&user).
				Update("minutes", gorm.Expr("minutes + ?", record.Minutes)).
				Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, record)
}

func getRecordList(c *gin.Context) {
	goalIDStr := c.DefaultQuery("goalID", "0")
	goalID, err := strconv.ParseUint(goalIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	cursorStr := c.DefaultQuery("cursor", "0")
	common.Logger.Debug("cursor ", cursorStr)
	cursor, err := strconv.ParseUint(cursorStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	missionIDStr := c.DefaultQuery("missionID", "0")
	missionID, err := strconv.ParseUint(missionIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var records []models.RecordModel

	query := &models.RecordModel{
		UserID:    c.MustGet(common.ContextUserID).(uint64),
		GoalID:    goalID,
		MissionID: missionID,
	}
	db := common.MySQLClient.
		Preload("Goal").
		Order("id desc").
		Preload("Mission").
		Where(query)

	if cursor > 0 {
		db = db.Where("ID < ?", cursor)
	}
	err = db.
		Limit(10).
		Find(&records).Error
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, records)
}

func getRecordHelper(uid, recordID interface{}) (record *models.RecordModel, err error) {
	record = &models.RecordModel{}
	err = common.MySQLClient.
		Where("user_id = ?", uid).
		// Preload("User").
		Preload("Goal").
		Preload("Mission").
		First(record, recordID).Error
	return
}

func getRecord(c *gin.Context) {
	mission, err := getRecordHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "record")
		return
	}
	c.JSON(http.StatusOK, *mission)
}

func deleteRecord(c *gin.Context) {
	uid := c.MustGet(common.ContextUserID).(uint64)
	record, err := getRecordHelper(uid, c.Param("id"))
	if err != nil {
		handleDBError(c, err, "record")
		return
	}

	err = common.MySQLClient.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Delete(record).Error; err != nil {
			return err
		}
		if record.Minutes > 0 {
			if record.Goal.ID > 0 {
				goal := models.GoalModel{}
				goal.ID = record.Goal.ID
				if err := tx.Model(&goal).
					Update("minutes", gorm.Expr("minutes - ?", record.Minutes)).
					Error; err != nil {
					return err
				}
			}
			if record.Mission.ID > 0 {
				mission := models.MissionModel{}
				mission.ID = record.Mission.ID
				if err := tx.Model(&mission).
					Update("minutes", gorm.Expr("minutes - ?", record.Minutes)).
					Error; err != nil {
					return err
				}
			}
			user := models.UserModel{}
			user.ID = uid
			if err := tx.Model(&user).
				Update("minutes", gorm.Expr("minutes - ?", record.Minutes)).
				Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err := common.MySQLClient.Delete(record).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Status(http.StatusOK)
}

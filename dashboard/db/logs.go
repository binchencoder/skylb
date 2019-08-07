package db

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/glog"
	"upper.io/db.v3/lib/sqlbuilder"

	pb "github.com/binchencoder/skylb/dashboard/proto"
	"github.com/binchencoder/skylb/dashboard/util"
)

const (
	tableLogs = "logs"

	sqlInsertLog = `INSERT INTO logs (operator, service_id, content, operation_time) VALUES (?, ?, ?, now())`
)

// Log represents a row in database table "logs".
type Log struct {
	Id        string    `db:"id"`
	Operator  string    `db:"operator"`
	ServiceId int32     `db:"service_id"`
	Content   string    `db:"content"`
	OpTime    time.Time `db:"operation_time"`
}

// ToLogInfo converts the Log struct to a pb.LogInfo struct.
func (l *Log) ToLogInfo() *pb.LogInfo {
	var svcName string
	if name, ok := util.ServiceIDsToNames[l.ServiceId]; ok {
		svcName = name
	} else {
		svcName = strconv.FormatInt(int64(l.ServiceId), 10)
	}
	info := pb.LogInfo{
		Operator: l.Operator,
		Service:  svcName,
		Content:  l.Content,
		OpTime:   l.OpTime.UnixNano() / int64(time.Millisecond),
	}
	return &info
}

// CreateLog inserts an operation log into database.
func CreateLog(operator string, serviceId int32, content string) error {
	tx, err := dbClient.NewTx(context.Background())
	if err != nil {
		glog.Errorf("Failed to begin transaction, %v", err)
		return err
	}

	if err := withTx(tx, func() error {
		return insertLog(tx, operator, serviceId, content)
	}); err != nil {
		glog.Errorf("Failed to insert log to database, %v", err)
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func insertLog(tx sqlbuilder.Tx, operator string, serviceId int32, content string) error {
	if _, err := tx.Exec(sqlInsertLog, operator, serviceId, content); err != nil {
		return err
	}
	return nil
}

// GetLogs returns log records.
func GetLogs(ctx context.Context, operator string, serviceId int32) ([]*Log, error) {
	tx, err := dbClient.NewTx(ctx)
	if err != nil {
		glog.Errorf("Failed to begin transaction, %v", err)
		return nil, err
	}

	filters := make([]interface{}, 0, 4)
	if len(operator) > 0 {
		filters = append(filters, "operator", operator)
	}
	if serviceId > -1 {
		filters = append(filters, "service_id", serviceId)
	}

	logs := []*Log{}
	err = withTx(tx, func() error {
		if err := tx.Collection(tableLogs).Find(filters...).OrderBy("operation_time DESC").Limit(100).All(&logs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		tx.Rollback()
		glog.Errorf("Failed to load logs from database, %v", err)
		return nil, err
	}
	tx.Commit()

	return logs, nil
}

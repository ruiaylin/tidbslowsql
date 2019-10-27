package tidbTools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	cp "tisql/components"
	"tisql/model"

	"github.com/jinzhu/gorm"
)

type DBHelper struct {
	// DBReader connection
	DBReader *gorm.DB
	// DBWriter writer
	DBWriter *gorm.DB
}

// NewDBH init database connection
func (dbh *DBHelper) NewDBH() {
	dbh.DBReader = cp.DBReader
	dbh.DBWriter = cp.DBWriter
}

// GetSlowSQL query slow sql
func GetSlowSQL(sqlID uint32) {
	fmt.Println("-----------------------------------------")
	var slowSQL model.SlowSQL
	// cp.DBReader.Where("sql_id = ?", sqlID).Find(&slowSQL)
	// cp.DBReader.Find(&slowSQL)
	notFound := cp.DBReader.Where(model.SlowSQL{Cluster: "feed", SQLID: sqlID}).First(&slowSQL).RecordNotFound()
	fmt.Println("notFound = ", notFound)
	if notFound {
		var slowSQL model.SlowSQL
		slowSQL.Cluster = "feed"
		slowSQL.SQLID = sqlID
		slowSQL.SQLText = "sq.SQLText xxx"
		slowSQL.Affirm = 0
		slowSQL.AffirmPerson = ""
		slowSQL.FirstSeen = time.Now()
		slowSQL.Counter = 1
		slowSQL.UpdateTime = time.Now()
		cp.DBWriter.Create(&slowSQL)
	} else {
		fmt.Println("slowSQL = ", slowSQL)
		fmt.Println(slowSQL.ID)
		fmt.Println(slowSQL.Cluster)
		fmt.Println(slowSQL.SQLID)
		fmt.Println(slowSQL.SQLText)
		fmt.Println(slowSQL.FirstSeen)
		fmt.Println(slowSQL.Affirm)
		fmt.Println(slowSQL.AffirmPerson)
		fmt.Println(slowSQL.Counter)
		fmt.Println(slowSQL.UpdateTime)
		counter := slowSQL.Counter + 1
		slowSQL.UpdateTime = time.Now()
		cp.DBWriter.Model(&slowSQL).Update(model.SlowSQL{Counter: counter, UpdateTime: time.Now()})
	}
}

// GetSlowQuery query slow sql
func GetSlowQuery(queryID int64) {
	var slowQuery model.SlowQuery
	cp.DBReader.Where("query_id = ?", queryID).Find(&slowQuery)
	// cp.DBReader.Find(&slowSQL)
	fmt.Println(slowQuery.ID)
	fmt.Println(slowQuery.Cluster)
	fmt.Println(slowQuery.SQLID)
	fmt.Println(slowQuery.QueryID)
	fmt.Println(slowQuery.SQLText)
	fmt.Println(slowQuery.SQLFullText)
	fmt.Println(slowQuery.Counter)
	fmt.Println(slowQuery.CreateTime)
	const TimeTMP = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dt, _ := time.ParseInLocation(TimeTMP, "2019-07-08 19:00:00", loc)
	var slowSQL1, ssql model.SlowSQL
	slowSQL1.Cluster = "feed"
	slowSQL1.SQLID = 12
	slowSQL1.SQLText = "sq.SQLText"
	slowSQL1.Affirm = 0
	slowSQL1.AffirmPerson = ""
	slowSQL1.FirstSeen = dt
	slowSQL1.Counter = 1
	slowSQL1.UpdateTime = dt
	cp.DBReader.Where(model.SlowSQL{Cluster: "feed", SQLID: 12}).Assign(slowSQL1).FirstOrCreate(&ssql)
}

// GetSlowQueryInfo query slow sql
func (dbh *DBHelper) GetSlowQueryInfo(queryID int64) {
	var slowQuery model.SlowQuery
	cp.DBReader.Where("query_id = ?", queryID).Find(&slowQuery)
	// cp.DBReader.Find(&slowSQL)
	fmt.Println(slowQuery.ID)
	fmt.Println(slowQuery.Cluster)
	fmt.Println(slowQuery.SQLID)
	fmt.Println(slowQuery.QueryID)
	fmt.Println(slowQuery.SQLText)
	fmt.Println(slowQuery.SQLFullText)
	fmt.Println(slowQuery.Counter)
	fmt.Println(slowQuery.CreateTime)
}

// HandleSQL handel base abstract sql info
func (dbh *DBHelper) HandleSQL(sq TiDBSlowQuery) {
	// handel time format
	const TimeTMP = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dt, _ := time.ParseInLocation(TimeTMP, sq.DayTime, loc)
	// struct for slow sql
	var slowSQL model.SlowSQL
	// check existing
	notFound := dbh.DBReader.Where(model.SlowSQL{Cluster: sq.Cluster, SQLID: sq.SQLID}).First(&slowSQL).RecordNotFound()
	fmt.Println("notFound = ", notFound)
	if notFound {
		// create one
		slowSQL.Cluster = sq.Cluster
		slowSQL.SQLID = sq.SQLID
		slowSQL.SQLText = sq.SQLText
		slowSQL.Affirm = 0
		slowSQL.AffirmPerson = ""
		slowSQL.FirstSeen = dt
		slowSQL.Counter = 1
		slowSQL.UpdateTime = time.Now()
		cp.DBWriter.Create(&slowSQL)
		if err := cp.DBWriter.Create(&slowSQL).Error; err != nil {
			fmt.Println("index  = ", strings.Index(err.Error(), "Error 1062"))
			if -1 != strings.Index(err.Error(), "Error 1062") {
				fmt.Println(" xx unique key ", err)
			}
		}
	}

}

// HandleQuery handle the real query
func (dbh *DBHelper) HandleQuery(sq TiDBSlowQuery) {
	// handel time format
	const TimeTMP = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dt, _ := time.ParseInLocation(TimeTMP, sq.DayTime, loc)
	// struct for slow query
	var slowQuery model.SlowQuery
	// check existing
	notFound := dbh.DBReader.Where(
		model.SlowQuery{Cluster: sq.Cluster, SQLID: sq.SQLID, QueryID: sq.QueryID}).First(
		&slowQuery).RecordNotFound()
	fmt.Println("notFound = ", notFound)
	instName := sq.IPAddr + ":" + strconv.FormatInt(sq.Port, 10)
	if notFound {
		// add a new query
		slowQuery.Cluster = sq.Cluster
		slowQuery.Host = sq.Hostname
		slowQuery.Port = sq.Port
		slowQuery.InstName = instName
		slowQuery.SQLID = sq.SQLID
		slowQuery.QueryID = sq.QueryID
		slowQuery.SQLText = sq.SQLText
		slowQuery.SQLFullText = sq.FullSQLText
		slowQuery.Counter = 1
		slowQuery.CreateTime = time.Now()
		slowQuery.UpdateTime = time.Now()
		cp.DBWriter.Create(&slowQuery)
		if err := cp.DBWriter.Create(&slowQuery).Error; err != nil {
			fmt.Println("index  = ", strings.Index(err.Error(), "1062"))
			if -1 != strings.Index(err.Error(), "Error 1062") {
				fmt.Println("unique key error ", err)
			}
		}
	}

	// insert the detail attribute to the query informantion table
	var qInfo model.SlowQueryInfo
	qInfo.Hostname = sq.Hostname
	qInfo.IP = sq.IPAddr
	qInfo.Port = sq.Port
	qInfo.InstName = instName
	qInfo.DayTime = dt
	qInfo.Cluster = sq.Cluster
	qInfo.SQLID = sq.SQLID
	qInfo.QueryID = sq.QueryID
	qInfo.SchemaName = sq.Schema
	qInfo.TabName = sq.Table
	qInfo.StartTS = int64(sq.TxnStartTS)
	qInfo.Username = sq.Username
	qInfo.ClientHost = sq.ClientHost
	qInfo.ConnID = int64(sq.ConnID)
	qInfo.QueryTime = sq.QueryTime
	qInfo.ProcessTime = sq.ProcessTime
	qInfo.WaitTime = sq.WaitTime
	qInfo.BackoffTime = sq.BackoffTime
	qInfo.RequestCount = int64(sq.RequestCount)
	qInfo.TotalKeys = int64(sq.TotalKeys)
	qInfo.ProcessKeys = int64(sq.ProcessKeys)
	qInfo.PrewriteTime = sq.PrewriteTime
	qInfo.CommitTime = sq.CommitTime
	qInfo.GetCommitTSTime = sq.GetCommitTST
	qInfo.TotalBackoffTime = sq.TotalBackoffTime
	qInfo.LatchWaitTime = sq.LatchWaitTime
	qInfo.ResolveLockTime = sq.ResolveLockTime
	qInfo.WriteKeys = int64(sq.WriteKeys)
	qInfo.WriteSize = int64(sq.WriteSize)
	qInfo.PreWriteRegion = int64(sq.PreWriteRegion)
	qInfo.TrxRetry = int64(sq.TrxRetry)
	qInfo.IndexIds = sq.IndexIDS
	qInfo.IsInternal = sq.IsInternal
	qInfo.Digest = sq.Digest
	qInfo.Stats = sq.Stats
	qInfo.CopProcAVG = sq.CopProcAVG
	qInfo.CopProc90 = sq.CopProcP90
	qInfo.CopProcMAX = sq.CopProcMAX
	qInfo.CopProcAddr = sq.CopProcAddr
	qInfo.CopWaitAVG = sq.CopWaitAVG
	qInfo.CopWait90 = sq.CopWaitP90
	qInfo.CopWaitMAX = sq.CopWaitMAX
	qInfo.CopWaitAddr = sq.CopWaitAddr
	qInfo.NumCopTasks = int64(sq.NumCopTasks)
	qInfo.MemMax = int64(sq.MemMax)
	qInfo.Succ = sq.Succ
	qInfo.CreateTime = time.Now()
	dbh.DBWriter.Create(&qInfo)
}

// HandleSQLWithCounter handel base abstract sql info
func (dbh *DBHelper) HandleSQLWithCounter(sq TiDBSlowQuery) {
	// handel time format
	const TimeTMP = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dt, _ := time.ParseInLocation(TimeTMP, sq.DayTime, loc)
	// struct for slow sql
	var slowSQL model.SlowSQL
	// check existing
	notFound := dbh.DBReader.Where(model.SlowSQL{Cluster: sq.Cluster, SQLID: sq.SQLID}).First(&slowSQL).RecordNotFound()
	fmt.Println("notFound = ", notFound)
	if notFound {
		// create one
		slowSQL.Cluster = sq.Cluster
		slowSQL.SQLID = sq.SQLID
		slowSQL.SQLText = sq.SQLText
		slowSQL.Affirm = 0
		slowSQL.AffirmPerson = ""
		slowSQL.FirstSeen = dt
		slowSQL.Counter = 1
		slowSQL.UpdateTime = time.Now()
		cp.DBWriter.Create(&slowSQL)
	} else {
		// increment the counter
		fmt.Println(slowSQL.ID)
		fmt.Println(slowSQL.Cluster)
		fmt.Println(slowSQL.SQLID)
		fmt.Println(slowSQL.SQLText)
		fmt.Println(slowSQL.FirstSeen)
		fmt.Println(slowSQL.Affirm)
		fmt.Println(slowSQL.AffirmPerson)
		fmt.Println(slowSQL.Counter)
		fmt.Println(slowSQL.UpdateTime)
		counter := slowSQL.Counter + 1
		slowSQL.UpdateTime = time.Now()
		// cp.DBWriter.Model(&slowSQL).Update(map[string]interface{}{"counter": counter, "update_time": time.Now()})
		dbh.DBWriter.Model(&slowSQL).Update(model.SlowSQL{Counter: counter, UpdateTime: time.Now()})
	}
}

// HandleQueryWithCounter handle the real query
func (dbh *DBHelper) HandleQueryWithCounter(sq TiDBSlowQuery) {
	// handel time format
	const TimeTMP = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dt, _ := time.ParseInLocation(TimeTMP, sq.DayTime, loc)
	// struct for slow query
	var slowQuery model.SlowQuery
	// check existing
	notFound := dbh.DBReader.Where(
		model.SlowQuery{Cluster: sq.Cluster, QueryID: sq.QueryID}).First(
		&slowQuery).RecordNotFound()
	fmt.Println("notFound = ", notFound)

	instName := sq.IPAddr + ":" + strconv.FormatInt(sq.Port, 10)

	if notFound {
		// add a new query
		slowQuery.Cluster = sq.Cluster
		slowQuery.Host = sq.Hostname
		slowQuery.Port = sq.Port
		slowQuery.InstName = instName
		slowQuery.SQLID = sq.SQLID
		slowQuery.QueryID = sq.QueryID
		slowQuery.SQLText = sq.SQLText
		slowQuery.SQLFullText = sq.FullSQLText
		slowQuery.Counter = 1
		slowQuery.CreateTime = time.Now()
		slowQuery.UpdateTime = time.Now()
		cp.DBWriter.Create(&slowQuery)
	} else {
		// increment the counter
		fmt.Println(slowQuery.ID)
		fmt.Println(slowQuery.Cluster)
		fmt.Println(slowQuery.SQLID)
		fmt.Println(slowQuery.SQLText)
		fmt.Println(slowQuery.SQLFullText)
		fmt.Println(slowQuery.Counter)
		counter := slowQuery.Counter + 1
		// cp.DBWriter.Model(&slowSQL).Update(map[string]interface{}{"counter": counter, "update_time": time.Now()})
		dbh.DBWriter.Model(&slowQuery).Update(
			model.SlowQuery{Counter: counter, UpdateTime: time.Now()})
	}

	// insert the detail attribute to the query informantion table
	var qInfo model.SlowQueryInfo
	qInfo.Hostname = sq.Hostname
	qInfo.IP = sq.IPAddr
	qInfo.Port = sq.Port
	qInfo.InstName = instName
	qInfo.DayTime = dt
	qInfo.Cluster = sq.Cluster
	qInfo.SQLID = sq.SQLID
	qInfo.QueryID = sq.QueryID
	qInfo.SchemaName = sq.Schema
	qInfo.TabName = sq.Table
	qInfo.StartTS = int64(sq.TxnStartTS)
	qInfo.Username = sq.Username
	qInfo.ClientHost = sq.ClientHost
	qInfo.ConnID = int64(sq.ConnID)
	qInfo.QueryTime = sq.QueryTime
	qInfo.ProcessTime = sq.ProcessTime
	qInfo.WaitTime = sq.WaitTime
	qInfo.BackoffTime = sq.BackoffTime
	qInfo.RequestCount = int64(sq.RequestCount)
	qInfo.TotalKeys = int64(sq.TotalKeys)
	qInfo.ProcessKeys = int64(sq.ProcessKeys)
	qInfo.PrewriteTime = sq.PrewriteTime
	qInfo.CommitTime = sq.CommitTime
	qInfo.GetCommitTSTime = sq.GetCommitTST
	qInfo.TotalBackoffTime = sq.TotalBackoffTime
	qInfo.LatchWaitTime = sq.LatchWaitTime
	qInfo.ResolveLockTime = sq.ResolveLockTime
	qInfo.WriteKeys = int64(sq.WriteKeys)
	qInfo.WriteSize = int64(sq.WriteSize)
	qInfo.PreWriteRegion = int64(sq.PreWriteRegion)
	qInfo.TrxRetry = int64(sq.TrxRetry)
	qInfo.IndexIds = sq.IndexIDS
	qInfo.IsInternal = sq.IsInternal
	qInfo.Digest = sq.Digest
	qInfo.Stats = sq.Stats
	qInfo.CopProcAVG = sq.CopProcAVG
	qInfo.CopProc90 = sq.CopProcP90
	qInfo.CopProcMAX = sq.CopProcMAX
	qInfo.CopProcAddr = sq.CopProcAddr
	qInfo.CopWaitAVG = sq.CopWaitAVG
	qInfo.CopWait90 = sq.CopWaitP90
	qInfo.CopWaitMAX = sq.CopWaitMAX
	qInfo.CopWaitAddr = sq.CopWaitAddr
	qInfo.NumCopTasks = int64(sq.NumCopTasks)
	qInfo.MemMax = int64(sq.MemMax)
	qInfo.Succ = sq.Succ
	qInfo.CreateTime = time.Now()
	dbh.DBWriter.Create(&qInfo)
}

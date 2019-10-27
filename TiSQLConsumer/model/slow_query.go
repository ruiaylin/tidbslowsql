package model

import "time"

// SlowSQL for abstract sql
type SlowSQL struct {
	ID           uint      `gorm:"primary_key"`
	Cluster      string    `gorm:"size:32;not null"`
	SQLID        uint32    `gorm:"column:sql_id" json:"sql_id"`
	SQLText      string    `gorm:"column:sql_text;type:longtext" json:"sql_text"`
	FirstSeen    time.Time `gorm:"column:first_seen;type:datetime;default:current_timestamp" json:"first_seen"`
	Affirm       int       `gorm:"column:affirm" json:"affirm" json:"affirm"`
	AffirmPerson string    `gorm:"column:affirm_person;size:128;not null;default:'' " json:"affirm_person"`
	Counter      int64     `gorm:"column:counter;type:bigint" json:"counter"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`
}

// TableName for the table's name
func (e *SlowSQL) TableName() string {
	return `tidb_slow_sql`
}

// SlowQuery for the real sql
type SlowQuery struct {
	ID          uint      `gorm:"primary_key"`
	Cluster     string    `gorm:"size:32;not null"`
	Host        string    `gorm:"size:32;not null"`
	Port        int64     `json:"port"`
	InstName    string    `gorm:"size:32;not null"`
	SQLID       uint32    `gorm:"column:sql_id" json:"sql_id"`
	QueryID     uint32    `gorm:"column:query_id" json:"query_id"`
	SQLText     string    `gorm:"column:sql_text;type:longtext" json:"sql_text"`
	SQLFullText string    `gorm:"column:sql_fulltext;type:longtext" json:"sql_fulltext"`
	Counter     int64     `gorm:"column:counter;type:bigint" json:"counter"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;default:current_timestamp" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;default:current_timestamp" json:"update_time"`
}

// TableName for the table's name
func (e *SlowQuery) TableName() string {
	return `tidb_slow_query`
}

// SlowQueryInfo for the real sql
type SlowQueryInfo struct {
	ID               uint      `gorm:"primary_key"`
	Hostname         string    `gorm:"size:32;not null"`
	IP               string    `gorm:"size:32;not null"`
	Port             int64     `json:"port"`
	InstName         string    `gorm:"size:32;not null"`
	DayTime          time.Time `gorm:"column:day_time;type:datetime;" json:"day_time"`
	Cluster          string    `gorm:"size:32;not null"`
	SQLID            uint32    `gorm:"column:sql_id" json:"sql_id"`
	QueryID          uint32    `gorm:"column:query_id" json:"query_id"`
	SchemaName       string    `gorm:"size:64;not null"`
	TabName          string    `gorm:"column:table_name;size:128;not null"`
	StartTS          int64     `gorm:"column:start_ts;type:bigint" json:"start_ts"`
	Username         string    `gorm:"size:32;not null"`
	ClientHost       string    `gorm:"size:16;not null"`
	ConnID           int64     `gorm:"column:conn_id;type:bigint" json:"conn_id"`
	QueryTime        float64   `gorm:"column:query_time;type:double" json:"query_time"`
	ProcessTime      float64   `gorm:"column:process_time;type:double" json:"process_time"`
	WaitTime         float64   `gorm:"column:wait_time;type:double" json:"wait_time"`
	BackoffTime      float64   `gorm:"column:backoff_time;type:double" json:"backoff_time"`
	RequestCount     int64     `gorm:"column:request_count;type:bigint" json:"request_count"`
	TotalKeys        int64     `gorm:"column:total_keys;type:bigint" json:"total_keys"`
	ProcessKeys      int64     `gorm:"column:process_keys;type:bigint" json:"process_keys"`
	PrewriteTime     float64   `gorm:"column:prewrite_time;type:double" json:"prewrite_time"`
	CommitTime       float64   `gorm:"column:commit_time;type:double" json:"commit_time"`
	GetCommitTSTime  float64   `gorm:"column:get_commit_ts_time;type:double" json:"get_commit_ts_time"`
	TotalBackoffTime float64   `gorm:"column:total_backoff_time;type:double" json:"total_backoff_time"`
	LatchWaitTime    float64   `gorm:"column:latch_wait_time;type:double" json:"latch_wait_time"`
	ResolveLockTime  float64   `gorm:"column:resolve_lock_time;type:double" json:"resolve_lock_time"`
	WriteKeys        int64     `gorm:"column:write_keys;type:bigint" json:"write_keys"`
	WriteSize        int64     `gorm:"column:write_size;type:bigint" json:"write_size"`
	PreWriteRegion   int64     `gorm:"column:prewrite_region;type:bigint" json:"prewrite_region"`
	TrxRetry         int64     `gorm:"column:trx_retry;type:int" json:"trx_retry"`
	IndexIds         string    `gorm:"column:index_ids;size:128;not null" json:"index_ids"`
	IsInternal       string    `gorm:"column:is_internal;size:10;not null" json:"is_internal"`
	Digest           string    `gorm:"column:digest;size:128;not null" json:"digest"`
	Stats            string    `gorm:"column:stats;size:128;not null" json:"stats"`
	CopProcAVG       float64   `gorm:"column:cop_proc_avg;type:double" json:"cop_proc_avg"`
	CopProc90        float64   `gorm:"column:cop_proc_90;type:double" json:"cop_proc_90"`
	CopProcMAX       float64   `gorm:"column:cop_proc_max;type:double" json:"cop_proc_max"`
	CopProcAddr      string    `gorm:"column:cop_proc_addr;size:128;not null" json:"cop_proc_addr"`
	CopWaitAVG       float64   `gorm:"column:cop_wait_avg;type:double" json:"cop_wait_avg"`
	CopWait90        float64   `gorm:"column:cop_wait_90;type:double" json:"cop_wait_90"`
	CopWaitMAX       float64   `gorm:"column:cop_wait_max;type:double" json:"cop_wait_max"`
	CopWaitAddr      string    `gorm:"column:cop_wait_addr;size:128;not null" json:"cop_wait_addr"`
	NumCopTasks      int64     `gorm:"column:num_cop_tasks;type:int" json:"num_cop_tasks"`
	MemMax           int64     `gorm:"column:mem_max;type:bigint" json:"mem_max"`
	Succ             string    `gorm:"column:succ;" json:"succ"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;" json:"create_time"`
}

// TableName for the table's name
func (e *SlowQueryInfo) TableName() string {
	return `tidb_slow_query_info`
}

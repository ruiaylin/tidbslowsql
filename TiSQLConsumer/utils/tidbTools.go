package tidbTools

import (
	"errors"
	"fmt"
	"hash/crc32"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tisql/message"

	_ "github.com/go-sql-driver/mysql"
	driver "github.com/pingcap/tidb/types/parser_driver"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

var (
	rmLineDelimiter, _ = regexp.Compile("|\n|\r\n|\t+|`")
	rmstartComment, _  = regexp.Compile(`(\/{2,}.*?(\r|\n))|(\/\*(\n|.)*?\*\/)`)
	reSpace, _         = regexp.Compile(`\s+`)
	reBackQuote, _     = regexp.Compile("`")
	ifUseWhere, _      = regexp.Compile(`(?i:where)`)
	// remove number and &
	rmVars, _         = regexp.Compile(`\(([0-9]+,?)+\)|\(('[^',]+',?)+\)|[0-9]+|'[^',]*'`)
	rmVstr, _         = regexp.Compile(`=(\s*).*((\".*\")|(\'.*\'))\s`)
	rmSelectFrom, _   = regexp.Compile(`(?i:select.*from)`)
	reSchemaTable1, _ = regexp.Compile(`(?i:( from )(.*)( force | inner |join|leftjoin|rightjoin))`)
	reSchemaTable2, _ = regexp.Compile(`(?i:( from )(.*)( where ))`)
	reSchemaTable3, _ = regexp.Compile(`(?i:( from )(.*)( limit ))`)
	reSchemaTable4, _ = regexp.Compile(`(?i:( from )(.*)(;))`)
)

var ops = map[string]string{
	"&&": "and",
	"||": "or",
	"&":  "and",
	"|":  "or",
	"eq": "=",
	"gt": ">",
	"lt": "<",
	"ge": ">=",
	"le": "<=",
	"ne": "<>",
}

// TiDBSlowQuery the structure of tidb slow query
type TiDBSlowQuery struct {
	// crc32 caculate sql_id & query_id
	SQLID    uint32
	QueryID  uint32
	Hostname string
	// for instance
	IPAddr  string
	Port    int64
	DayTime string
	Cluster string
	//schema & table
	Schema string
	Table  string
	//query SQL
	SQLText     string
	FullSQLText string
	// transaction start ts
	TxnStartTS float64
	//User username@host adress
	Username   string
	ClientHost string
	ConnID     float64
	QueryTime  float64
	// tikv
	ProcessTime  float64
	WaitTime     float64
	BackoffTime  float64
	RequestCount float64
	TotalKeys    float64
	ProcessKeys  float64
	PrewriteTime float64
	CommitTime   float64
	// -- Get_commit_ts_time
	GetCommitTST     float64
	TotalBackoffTime float64
	LatchWaitTime    float64
	ResolveLockTime  float64
	// keys
	WriteKeys      float64
	WriteSize      float64
	PreWriteRegion float64
	TrxRetry       float64
	// other
	IndexIDS   string
	IsInternal string
	Digest     string
	Stats      string
	// cop proc
	CopProcAVG  float64
	CopProcP90  float64
	CopProcMAX  float64
	CopProcAddr string
	// cop wait
	CopWaitAVG  float64
	CopWaitP90  float64
	CopWaitMAX  float64
	CopWaitAddr string
	// task
	NumCopTasks float64
	MemMax      float64
	Succ        string
}

// InitTiDBQuery init query
func InitTiDBQuery(msg *message.Message) *TiDBSlowQuery {
	logName := msg.GetLogger()
	fmt.Printf("logName : %s \n", logName)
	hpt := strings.Split(logName, ".")
	fmt.Printf("hpt[1] : %s \n", hpt[1])
	ips := strings.Split(hpt[1], "-")
	cluster := ips[5]
	port, _ := strconv.ParseInt(ips[4], 10, 64)
	ip := strings.Join(ips[:4], ".")
	fmt.Println("IP = ", ip)
	mQuery := &TiDBSlowQuery{
		Cluster:  cluster,
		Hostname: msg.GetHostname(),
		IPAddr:   ip,
		Port:     port,
	}
	return mQuery
}

// GenID gen crc32 id for the sql
func (tsql *TiDBSlowQuery) GenID() error {
	var IDERR = errors.New("General crc32 ID failed or sql is empty!!!")
	// fmt.Println("fullsql = ", tsql.FullSQLText)
	sql := strings.ToLower(tsql.FullSQLText)
	sql = rmLineDelimiter.ReplaceAllString(sql, "") // delete \n \t \r\n  in sql
	sql = rmstartComment.ReplaceAllString(sql, "")
	fullsql := reSpace.ReplaceAllString(sql, " ") // change multi "    " to one " "
	if fullsql == "" {
		return IDERR
	} else {
		tsql.FullSQLText = strings.Replace(fullsql, ";", " ", -1)
		tsql.QueryID = crc32.ChecksumIEEE([]byte(fullsql))
	}
	sql = rmVars.ReplaceAllString(fullsql, "?") // delete from literals
	sql = rmVstr.ReplaceAllString(sql, "= ? ")  // delete from literals
	tsql.SQLText = sql
	tsql.SQLID = crc32.ChecksumIEEE([]byte(sql))
	// fmt.Println("fullsql = ", tsql.FullSQLText)
	// fmt.Println("sql = ", tsql.SQLText)

	return nil
}

//SetFields decode the fields
func (tsql *TiDBSlowQuery) SetFields(fields []*message.Field) error {
	for _, f := range fields {
		switch f.GetName() {
		case "Time":
			if len(f.GetValueString()) == 1 {
				timeTemplateNanos := "2006-01-02 15:04:05"
				strTime := f.GetValueString()[0]
				t, err := time.ParseInLocation("2006-01-02T15:04:05+08:00", strTime, time.Local)
				fmt.Println(t, err)

				timestamp := fmt.Sprintf("%s", t.Format(timeTemplateNanos))
				fmt.Printf("Timestamp : %s \n", timestamp)
				tsql.DayTime = timestamp
				continue
			}
		case "Txn_start_ts":
			if len(f.GetValueDouble()) == 1 {
				tsql.TxnStartTS = f.GetValueDouble()[0]
				continue
			}
		case "User":
			if len(f.GetValueString()) == 1 {
				user := f.GetValueString()[0]
				uh := strings.Split(user, "@")
				tsql.Username = uh[0]
				tsql.ClientHost = uh[1]
				continue
			}
		case "Conn_ID":
			if len(f.GetValueDouble()) == 1 {
				tsql.ConnID = f.GetValueDouble()[0]
				continue
			}
		case "Query_time":
			if len(f.GetValueString()) == 1 {
				tsql.QueryTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Process_time":
			if len(f.GetValueString()) == 1 {
				tsql.ProcessTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Wait_time":
			if len(f.GetValueString()) == 1 {
				tsql.WaitTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Backoff_time":
			if len(f.GetValueString()) == 1 {
				tsql.BackoffTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Request_count":
			if len(f.GetValueDouble()) == 1 {
				tsql.RequestCount = f.GetValueDouble()[0]
				continue
			}
		case "Total_keys":
			if len(f.GetValueDouble()) == 1 {
				tsql.TotalKeys = f.GetValueDouble()[0]
				continue
			}
		case "Process_keys":
			if len(f.GetValueDouble()) == 1 {
				tsql.ProcessKeys = f.GetValueDouble()[0]
				continue
			}
		case "Prewrite_time":
			if len(f.GetValueString()) == 1 {
				tsql.PrewriteTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Commit_time":
			if len(f.GetValueString()) == 1 {
				tsql.CommitTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Get_commit_ts_time":
			if len(f.GetValueString()) == 1 {
				tsql.GetCommitTST = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Total_backoff_time":
			if len(f.GetValueString()) == 1 {
				tsql.TotalBackoffTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Local_latch_wait_time":
			if len(f.GetValueString()) == 1 {
				tsql.LatchWaitTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Resolve_lock_time":
			if len(f.GetValueString()) == 1 {
				tsql.ResolveLockTime = ToFloat64(f.GetValueString()[0])
				continue
			}
		case "Write_keys":
			if len(f.GetValueDouble()) == 1 {
				tsql.WriteKeys = f.GetValueDouble()[0]
				continue
			}
		case "Write_size":
			if len(f.GetValueDouble()) == 1 {
				tsql.WriteSize = f.GetValueDouble()[0]
				continue
			}
		case "Prewrite_region":
			if len(f.GetValueDouble()) == 1 {
				tsql.PreWriteRegion = f.GetValueDouble()[0]
				continue
			}
		case "Txn_retry":
			if len(f.GetValueDouble()) == 1 {
				tsql.TrxRetry = f.GetValueDouble()[0]
				continue
			}
		case "DB":
			if len(f.GetValueString()) == 1 {
				tsql.Schema = f.GetValueString()[0]
				continue
			}
		case "Index_ids":
			if len(f.GetValueString()) == 1 {
				tsql.IndexIDS = f.GetValueString()[0]
				continue
			}
		case "Is_internal":
			if len(f.GetValueString()) == 1 {
				tsql.IsInternal = f.GetValueString()[0]
				continue
			}
		case "Digest":
			if len(f.GetValueString()) == 1 {
				tsql.Digest = f.GetValueString()[0]
				continue
			}
		case "Stats":
			if len(f.GetValueString()) == 1 {
				tsql.Stats = f.GetValueString()[0]
				continue
			}
		case "Cop_proc_avg":
			if len(f.GetValueDouble()) == 1 {
				tsql.CopProcAVG = f.GetValueDouble()[0]
				continue
			}
		case "Cop_proc_p90":
			if len(f.GetValueDouble()) == 1 {
				tsql.CopProcP90 = f.GetValueDouble()[0]
				continue
			}
		case "Cop_proc_max":
			if len(f.GetValueDouble()) == 1 {
				tsql.CopProcMAX = f.GetValueDouble()[0]
				continue
			}
		case "Cop_proc_addr":
			if len(f.GetValueString()) == 1 {
				tsql.CopProcAddr = strings.TrimSpace(f.GetValueString()[0])
				continue
			}

		case "Cop_wait_avg":
			if len(f.GetValueDouble()) == 1 {
				tsql.CopWaitAVG = f.GetValueDouble()[0]
				continue
			}
		case "Cop_wait_p90":
			if len(f.GetValueDouble()) == 1 {
				tsql.CopWaitP90 = f.GetValueDouble()[0]
				continue
			}
		case "Cop_wait_max":
			if len(f.GetValueDouble()) == 1 {
				tsql.CopWaitMAX = f.GetValueDouble()[0]
				continue
			}
		case "Cop_wait_addr":
			if len(f.GetValueString()) == 1 {
				tsql.CopWaitAddr = strings.TrimSpace(f.GetValueString()[0])
				continue
			}

		case "Num_cop_tasks":
			if len(f.GetValueDouble()) == 1 {
				tsql.NumCopTasks = f.GetValueDouble()[0]
				continue
			}

		case "Mem_max":
			if len(f.GetValueDouble()) == 1 {
				tsql.MemMax = f.GetValueDouble()[0]
				continue
			}

		case "Succ":
			if len(f.GetValueString()) == 1 {
				tsql.Succ = f.GetValueString()[0]
				continue
			}
		case "Sql":
			if len(f.GetValueString()) == 1 {
				tsql.FullSQLText = f.GetValueString()[0]
				continue
			}
		default:
		}
	}
	return nil
}

//PrintInfo print the base information
func (tsql *TiDBSlowQuery) PrintInfo() {
	fmt.Printf("---------------------------------------------------------------------\n")
	fmt.Println("tsql.SQLID             = ", tsql.SQLID)
	fmt.Println("tsql.QueryID           = ", tsql.QueryID)
	fmt.Println("tsql.HostName          = ", tsql.Hostname)
	fmt.Println("tsql.IPAddr            = ", tsql.IPAddr)
	fmt.Println("tsql.Port              = ", tsql.Port)
	fmt.Println("tsql.Cluster           = ", tsql.Cluster)
	fmt.Println("tsql.DayTime           = ", tsql.DayTime)
	fmt.Println("tsql.Schema            = ", tsql.Schema)
	fmt.Println("tsql.Table             = ", tsql.Table)
	fmt.Println("tsql.SQLText           = ", tsql.SQLText)
	fmt.Println("tsql.FullSQLText       = ", tsql.FullSQLText)
	fmt.Printf("tsql.TxnStartTS        =  %.f \n", tsql.TxnStartTS)
	fmt.Println("tsql.Username          = ", tsql.Username)
	fmt.Println("tsql.ClientHost        = ", tsql.ClientHost)
	fmt.Printf("tsql.ConnID            =  %.f \n", tsql.ConnID)
	fmt.Println("tsql.QueryTime         = ", tsql.QueryTime)
	fmt.Println("tsql.ProcessTime       = ", tsql.ProcessTime)
	fmt.Println("tsql.WaitTime          = ", tsql.WaitTime)
	fmt.Println("tsql.BackoffTime       = ", tsql.BackoffTime)
	fmt.Println("tsql.RequestCount      = ", tsql.RequestCount)
	fmt.Println("tsql.TotalKeys         = ", tsql.TotalKeys)
	fmt.Println("tsql.ProcessKeys       = ", tsql.ProcessKeys)
	fmt.Println("tsql.PrewriteTime      = ", tsql.PrewriteTime)
	fmt.Println("tsql.CommitTime        = ", tsql.CommitTime)
	fmt.Println("tsql.GetCommitTST      = ", tsql.GetCommitTST)
	fmt.Println("tsql.TotalBackoffTime  = ", tsql.TotalBackoffTime)
	fmt.Println("tsql.LatchWaitTime     = ", tsql.LatchWaitTime)
	fmt.Println("tsql.ResolveLockTime   = ", tsql.ResolveLockTime)
	fmt.Println("tsql.WriteKeys         = ", tsql.WriteKeys)
	fmt.Println("tsql.WriteSize         = ", tsql.WriteSize)
	fmt.Println("tsql.PreWriteRegion    = ", tsql.PreWriteRegion)
	fmt.Println("tsql.TrxRetry          = ", tsql.TrxRetry)
	fmt.Println("tsql.IndexIDS          = ", tsql.IndexIDS)
	fmt.Println("tsql.IsInternal        = ", tsql.IsInternal)
	fmt.Println("tsql.Digest            = ", tsql.Digest)
	fmt.Println("tsql.Stats             = ", tsql.Stats)
	fmt.Println("tsql.CopProcAVG        = ", tsql.CopProcAVG)
	fmt.Println("tsql.CopProcP90        = ", tsql.CopProcP90)
	fmt.Println("tsql.CopProcMAX        = ", tsql.CopProcMAX)
	fmt.Println("tsql.CopProcAddr       = ", tsql.CopProcAddr)
	fmt.Println("tsql.CopWaitAVG        = ", tsql.CopWaitAVG)
	fmt.Println("tsql.CopWaitP90        = ", tsql.CopWaitP90)
	fmt.Println("tsql.CopWaitMAX        = ", tsql.CopWaitMAX)
	fmt.Println("tsql.CopWaitAddr       = ", tsql.CopWaitAddr)
	fmt.Println("tsql.NumCopTasks       = ", tsql.NumCopTasks)
	fmt.Printf("tsql.MemMax            =  %.f \n", tsql.MemMax)
	fmt.Println("tsql.Succ              = ", tsql.Succ)
}

// FormatSQL and SQLID
func (tsql *TiDBSlowQuery) FormatSQL() {
	fmt.Printf("before-sql: %s\n", tsql.FullSQLText)
	sql, err := ParserALLSQL(tsql.FullSQLText)
	fmt.Printf("before-sql: %s\n", sql)
	if err != nil {
		fmt.Println("AbSQL: ", sql)
	}
	tsql.SQLText = sql
	tsql.SQLID = crc32.ChecksumIEEE([]byte(sql))
}

// ParserALLSQL and SQLID
func ParserALLSQL(sql string) (absql string, err error) {
	abSQL := ""
	parser := parser.New()
	// stmts, err := parser.Parse(src, "", "")
	stmt, err := parser.ParseOneStmt(sql, "", "")
	// fmt.Printf(" xx-1 :  %s  \n", reflect.TypeOf(stmt))
	if err != nil {
		fmt.Println("... err: ", err)
		return sql, err
	}
	switch node := stmt.(type) {
	case *ast.InsertStmt:
		abSQL = AbstractINSSQL(node)
	case *ast.DeleteStmt:
		abSQL = AbstractDELSQL(node)
	case *ast.UpdateStmt:
		abSQL = AbstractUPDSQL(node)
	case *ast.SelectStmt:
		abSQL = AbstractSELSQL(node)
	case *ast.UnionStmt:
		abSQL = AbstractUnionSQL(node)
	default:
		abSQL = sql
	}
	return abSQL, nil
}

// AbstractUnionSQL union
func AbstractUnionSQL(unn *ast.UnionStmt) (sql string) {
	tsql := ""
	// for SelectList
	// fmt.Println(" len(unn.SelectList.Selects) ", len(unn.SelectList.Selects))
	for i, sel := range unn.SelectList.Selects {
		tsql += "(" + AbstractSELSQL(sel) + ")"
		if i < len(unn.SelectList.Selects)-1 {
			tsql += " UNION "
		}
	}
	//OrderBy
	if unn.OrderBy != nil {
		tsql += " OrderBy " + getByItemNames(unn.OrderBy.Items)
	}
	if unn.Limit != nil {
		tsql += " LIMIT "
	}
	return tsql
}

// AbstractUPDSQL for update
func AbstractUPDSQL(upd *ast.UpdateStmt) (sql string) {
	tsql := "UPDATE "
	// tableIdent
	tsql += analyzeFrom(upd.TableRefs.TableRefs)
	// From
	if upd.Where != nil {
		tsql += " Where " + analyzeWhere(upd.Where)
	}
	//for order by
	if upd.Order != nil {
		tsql += " Orderby " + getByItemNames(upd.Order.Items)
	}
	//
	return tsql
}

// AbstractDELSQL for delete
func AbstractDELSQL(del *ast.DeleteStmt) (sql string) {
	tsql := "DELETE"
	tabs := ""
	tabs1 := ""
	if del.IsMultiTable {
		for _, tab := range del.Tables.Tables {
			if len(tab.Schema.String()) > 0 {
				tabs += fmt.Sprintf("%s.%s", tab.Schema, tab.Name) + ","
			} else {
				tabs += fmt.Sprintf("%s", tab.Name) + ","
			}
		}
		tabs = tabs[:len(tabs)-1]
		// format from
	}
	if del.TableRefs != nil {
		tabs1 = analyzeFrom(del.TableRefs.TableRefs)
	}
	if del.IsMultiTable && del.BeforeFrom {
		tsql += " " + tabs + " FROM " + tabs1
	} else if del.IsMultiTable && !del.BeforeFrom {
		tsql += " " + tabs
	} else {
		tsql += " " + tabs1
	}
	// where
	if del.Where != nil {
		tsql += " WHERE " + analyzeWhere(del.Where)
	}
	//OrderBy
	if del.Order != nil {
		tsql += " OrderBy " + getByItemNames(del.Order.Items)
	}
	return tsql
}

// AbstractINSSQL for insert
func AbstractINSSQL(inst *ast.InsertStmt) (sql string) {
	tsql := "INSERT "
	// ColNames
	tsql += analyzeFrom(inst.Table.TableRefs)
	// From
	if inst.Select != nil {
		switch nod := inst.Select.(type) {
		case *ast.SelectStmt:
			tsql += " " + AbstractSELSQL(nod)
		case *ast.UnionStmt:
			tsql += " " + AbstractUnionSQL(nod)
		default:
			fmt.Printf("xx-4 : %s %s  \n", nod, reflect.TypeOf(nod))
		}
	} else {
		tsql = fmt.Sprintf("%s(...*%v)", tsql, len(inst.Lists))
	}
	return tsql
}

// AbstractSELSQL for select
func AbstractSELSQL(sel *ast.SelectStmt) (sql string) {
	tsql := "SELECT"
	// ditinct
	if sel.Distinct {
		tsql += " DISTINCT"
	}
	//From
	if sel.From != nil {
		tsql += " " + analyzeFrom(sel.From.TableRefs)
	}
	//where
	if sel.Where != nil {
		tsql += " WHERE " + analyzeWhere(sel.Where)
	}
	//group by
	if sel.GroupBy != nil {
		tsql += " Groupby " + getByItemNames(sel.GroupBy.Items)
	}

	// having
	if sel.Having != nil {
		names := ""
		switch nod := sel.Having.Expr.(type) {
		case *ast.ColumnNameExpr:
			names += getColName(nod) + ","
		case *ast.AggregateFuncExpr:
			names += nod.F + ","
		case *ast.FuncCallExpr:
			names += getFuncInfo(nod) + ","
		default:
			fmt.Sprintf("[[[[[4]]]] : %s %s  \n", nod, reflect.TypeOf(nod))
		}
		tsql += " Having " + names
	}

	//OrderBy
	if sel.OrderBy != nil {
		tsql += " OrderBy " + getByItemNames(sel.OrderBy.Items)
	}
	// limit
	if sel.Limit != nil {
		tsql += " LIMIT "
		if sel.Limit.Offset != nil {
			tsql += "? "
		}
		if sel.Limit.Count != nil {
			tsql += "? "
		}
	}
	return tsql
}

// dofor where
func analyzeWhere(node interface{}) (wher string) {
	switch oper := node.(type) {
	case *ast.PatternInExpr:
		colName := ""
		if oper.Expr != nil {
			colName = analyzeExpr(oper.Expr)
		}
		// in or not in
		instr := ""
		if oper.Not {
			instr += "NOT "
		}
		instr += "IN"
		// sub select
		substr := ""
		if oper.Sel != nil {
			switch mnode := oper.Sel.(*ast.SubqueryExpr).Query.(type) {
			case *ast.SelectStmt:
				substr += "(" + AbstractSELSQL(mnode) + ")"
			case *ast.UnionStmt:
				substr += "(" + AbstractUnionSQL(mnode) + ")"
			default:
				substr += fmt.Sprintf("oper.Sel.(*ast.SubqueryExpr).Query : %s ", reflect.TypeOf(oper.Sel.(*ast.SubqueryExpr).Query))
			}
		}
		if len(substr) > 0 {
			return fmt.Sprintf("%s %s %s", colName, instr, substr)
		} else {
			return fmt.Sprintf("%s %s (..*%v)", colName, instr, len(oper.List))
		}
	case *ast.BinaryOperationExpr:
		lname := ""
		rname := ""
		if typeCheck(oper.L) == "ColumnNameExpr" {
			lname = getColName(oper.L.(*ast.ColumnNameExpr))
		} else if typeCheck(oper.L) == "ValueExpr" {
			lname = "?"
		}
		if typeCheck(oper.R) == "ColumnNameExpr" {
			rname = getColName(oper.R.(*ast.ColumnNameExpr))
		}
		if typeCheck(oper.R) == "ValueExpr" {
			rname = "?"
		}
		if typeCheck(oper.R) == "FuncCallExpr" {
			rname = getFuncInfo(oper.R.(*ast.FuncCallExpr))
		}
		op1 := oper.Op.String()
		opsItem, ok := ops[oper.Op.String()] // 假如key存在,则name = 李四 ，ok = true,否则，ok = false
		if ok {
			op1 = opsItem
		}
		if len(lname) > 0 {
			return fmt.Sprintf("%s %s %s", lname, op1, rname)
		}
		lstr := analyzeWhere(oper.L)
		rstr := analyzeWhere(oper.R)
		return fmt.Sprintf("%s %s %s", lstr, op1, rstr)
	case *ast.ParenthesesExpr:
		str := analyzeWhere(oper.Expr)
		return fmt.Sprintf("%s", str)
	case *ast.BetweenExpr:
		colName := ""
		if oper.Expr != nil {
			colName = analyzeExpr(oper.Expr)
		}
		return fmt.Sprintf("%s BETWEEN", colName)
	case *ast.PatternLikeExpr:
		// fmt.Println("PatternLikeExpr oper type ", reflect.TypeOf(oper))
		colName := ""
		if oper.Expr != nil {
			// fmt.Println("oper.Expr type ", reflect.TypeOf(oper.Expr))
			switch mnode := oper.Expr.(type) {
			case *ast.FuncCallExpr:
				funcName := analyzeWhere(mnode)
				return fmt.Sprintf("%s LIKE", funcName)
			case *ast.ColumnNameExpr:
				tmpName := oper.Expr.(*ast.ColumnNameExpr).Name
				if len(tmpName.Table.String()) > 0 {
					colName += tmpName.Table.String() + "."
				}
				colName += tmpName.Name.String()
				return fmt.Sprintf("%s LIKE", colName)
			default:
				return " - Like "
			}
		}
	case *ast.IsNullExpr:
		colName := ""
		if oper.Expr != nil {
			switch isNode := oper.Expr.(type) {
			case *ast.FuncCallExpr:
				colName = analyzeWhere(isNode)
			case *ast.ColumnNameExpr:
				tmpName := isNode.Name
				if len(tmpName.Table.String()) > 0 {
					colName += tmpName.Table.String() + "."
				}
				colName += tmpName.Name.String()
			default:
				colName = " ISNULL "
			}
		}
		if oper.Not {
			return fmt.Sprintf("%s IS NOT NULL", colName)
		}
		return fmt.Sprintf("%s IS NULL", colName)
	case *ast.FuncCallExpr:
		return getFuncInfo(oper)
	case *ast.ValueExpr:
		return fmt.Sprintf("?")
	case *driver.ValueExpr:
		return fmt.Sprintf("?")
	case *ast.ExistsSubqueryExpr:
		return genExistString(oper)
	case *ast.RowExpr:
		rstr := ""
		for _, expr := range oper.Values {
			rstr += analyzeExpr(expr) + ","
		}
		return fmt.Sprintf("(%s)", rstr[:len(rstr)-1])
	case *ast.SubqueryExpr:
		return fmt.Sprintf("(%s)", AbstractSELSQL(oper.Query.(*ast.SelectStmt)))
	case *ast.ColumnNameExpr:
		colName := oper.Name
		return fmt.Sprintf(" %s ", colName)
	default:
		fmt.Println("default execution where ", reflect.TypeOf(oper))
	}
	return ""
}

func analyzeExpr(node interface{}) (valName string) {
	colStr := ""
	switch nod := node.(type) {
	case *ast.ColumnNameExpr:
		tmpName := nod.Name
		if len(tmpName.Table.String()) > 0 {
			colStr += tmpName.Table.String() + "."
		}
		colStr += tmpName.Name.String()
		return colStr
	case *ast.ParenthesesExpr:
		return analyzeExpr(nod.Expr)
	case *ast.RowExpr:
		rstr := ""
		for _, expr := range nod.Values {
			rstr += analyzeExpr(expr) + ","
		}
		return fmt.Sprintf("(%s)", rstr[:len(rstr)-1])
	case *ast.FuncCallExpr:
		return getFuncInfo(nod)
	case *ast.ValueExpr:
		return "?"
	case *driver.ValueExpr:
		return fmt.Sprintf("?")
	default:
		fmt.Sprintf("%s", nod)
		return "oops"
	}
}

func genExistString(node *ast.ExistsSubqueryExpr) (exists string) {
	sts := ""
	if node.Not {
		sts += "NOT "
	}
	sts += "EXISTS"
	if node.Sel != nil {
		sts += " ( " + AbstractSELSQL(node.Sel.(*ast.SubqueryExpr).Query.(*ast.SelectStmt)) + " )"
	}
	return sts
}

func getByItemNames(bys []*ast.ByItem) (cols string) {
	names := ""
	for _, x := range bys {
		switch nod := x.Expr.(type) {
		case *ast.ColumnNameExpr:
			names += getColName(nod) + ","
		case *ast.AggregateFuncExpr:
			names += nod.F + ","
		case *ast.FuncCallExpr:
			names += getFuncInfo(nod) + ","
		case *driver.ValueExpr:
			val := nod.GetValue()
			sval := fmt.Sprintf("%d", val)
			names += sval + ","
		default:
			fmt.Sprintf("[[[[[4]]]] : %s %s  \n", nod, reflect.TypeOf(nod))
		}
	}
	if len(names) == 0 {
		return names
	}
	return names[:len(names)-1]
}

// get column name
func getFuncInfo(node *ast.FuncCallExpr) (name string) {
	fname := node.FnName
	// arg
	var argss []string
	for _, x := range node.Args {
		fmt.Println("x = ", x, reflect.TypeOf(x))
		switch nod := x.(type) {
		case *ast.ColumnNameExpr:
			argss = append(argss, getColName(nod))
		case *ast.AggregateFuncExpr:
			argss = append(argss, nod.F+"(?)")
		case *driver.ValueExpr:
			argss = append(argss, "?")
		case *ast.FuncCallExpr:
			argss = append(argss, getFuncInfo(nod))
		case *ast.TimeUnitExpr:
			fmt.Sprintln("nodx = ", node)
		default:
			fmt.Printf("[[[[[4]]]] : %s %s \n", nod, reflect.TypeOf(nod))
		}
	}
	argss1 := strings.Join(argss, ",")
	return fmt.Sprintf("%s(%s)", fname, argss1)
}

// get column name
func getColName(cols *ast.ColumnNameExpr) (name string) {
	colName := "" + cols.Name.Schema.String()
	if len(colName) > 0 {
		colName += "."
	}
	colName += cols.Name.Table.String()
	if len(colName) > 0 {
		colName += "."
	}
	colName += cols.Name.Name.String()
	return colName
}

// type check
func typeCheck(node interface{}) (tp string) {
	switch nod := node.(type) {
	case *ast.ColumnNameExpr:
		return "ColumnNameExpr"
	case *ast.ValueExpr:
		return "ValueExpr"
	case *driver.ValueExpr:
		return "ValueExpr"
	case *ast.FuncCallExpr:
		return "FuncCallExpr"
	default:
		fmt.Sprintf("%s", nod)
		// fmt.Println("typeCheck check", reflect.TypeOf(node))
		return "oops"
	}
}

//analyze From clouses
func analyzeFrom(node interface{}) (tabs string) {
	switch t := node.(type) {
	case *ast.TableSource:
		switch x := t.Source.(type) {
		case *ast.TableName:
			if len(x.Schema.String()) > 0 {
				// fmt.Println("...xxx...", x.Schema)
				return fmt.Sprintf("%s.%s", x.Schema, x.Name)
			}
			return fmt.Sprintf("%s", x.Name)
		case *ast.SelectStmt:
			tmpsql := AbstractSELSQL(x)
			return "(" + tmpsql + ")"
		case *ast.UnionStmt:
			tmpsql := AbstractUnionSQL(x)
			return "(" + tmpsql + ")"
		default:
			fmt.Println("default execution , type is :  ", reflect.TypeOf(x))
		}
	case *ast.Join:
		left := analyzeFrom(t.Left)
		if t.Right != nil {
			tp := t.Tp
			jstr := ""
			switch tp {
			case 1:
				jstr = "CrossJoin"
			case 2:
				jstr = "LeftJoin"
			case 3:
				jstr = "RightJoin"

			default:
				jstr = "CrossJoin"
				fmt.Sprintf("%d", tp)
			}
			right := analyzeFrom(t.Right)
			return left + " " + jstr + " " + right
		}
		if t.On != nil {
			// todo:
			fmt.Printf(" on.Expr :  %s ", reflect.TypeOf(t.On.Expr))
		}
		return left
	default:
		return fmt.Sprintf("default exec ~~~~ %s ", reflect.TypeOf(t))
	}
	return ""
}

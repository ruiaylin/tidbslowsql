package main

import (
	"fmt"
	"reflect"

	"github.com/ruiaylin/sqlparser/ast"
	"github.com/ruiaylin/sqlparser/parser"
)

var ops = map[string]string{
	"&&": "and",
	"||": "or",
	"&":  "and",
	"|":  "or",
}

func main() {
	TestParser0()
}

// TODO: table 43 and 50 parse failed
func TestParser0() {
	table := []struct {
		src string
		ok  bool
	}{
		{"select substring(created_time,1,10) as dt,count(distinct leads_id ) as num from order_sales_leads where substring(created_time,1,10)>='2015-10-01' and substring(created_time,1,10)<='2015-12-20' and stage_id <> 0 group by substring(created_time,1,10)  ", true},
		{"Select * From user Where DATE_FORMAT(birthday,'%m-%d') >= '06-03' and DATE_FORMAT(birthday,'%m-%d') <= '07-08'", true},
		{"select * from mysql.test where a.name like 'ruiaylin%' ", true},
		{"SELECT DISTINCT contractId,companyId,advertType from CampaignProperty where dates = '20170822' and contractId in (474582,474296,471479,475423,474496,474280,475469,475539,472497,475647,471923,472591,475937,475691,473926,476124,471068,476500,472979,477293,472970,472104,477287,477389,477403,477030,477500,477029,477668,471012,477662,477666,477675,475781,477674,477676,475782,477679,475784,477686,475982,477697,475816,475783,477760,475786,477761,475785,477762,476785,477776,472386,473932,478734,478840,479050,475520,479502,472996,477302,477152,476783,477258,474644,478863,473084,474402,479258,476632,474481,479262,474398,474474,479589,475852,474278,480260,480265,477321,473577,475839,473381,473443,474771,473977,474638,479419,480525,479426,479428,480510,477520,473937,476569,478962,479424,480107,481111,481114,481116)", true},
		{"select * from test t where t.name between  'ruiaylin1' and 'ruiaylin2' ", true},
		{"select * from test where  aget = 3 and name is null ", true},
		{"select * from test where name is not null and age in (1, 20 ) ", true},
		{"select * from test t where t.name in ('111','232232','1212')   ", true},
		{"select * from (select * FROM TEST union select 2 FROM TEST3) as a", true},
		{"insert into t select c1 from a.t1  where id = 2 ", true},
		{"insert into t select c1 from t1 union select c2 from t2", true},
		{"insert into t (c) select c1 from t1 union select c2 from t2", true},
		// from join
		{"SELECT * from t1, t2, t3", true},
		{"select * from t1 join t2 left join t3 on t2.id = t3.id", true},
		{"select * from t1 right join t2 on t1.id = t2.id left join t3 on t3.id = t2.id", true},
		{"select id, zan_type,f_id, zan_num, zan_user,is_flag , zan_users,zan_source from zan_info  where f_id = '46f2275476a81261c5e8117eed12b980'", false},
		{`select
    atr.id, atr.topicId, atr.topicSrc, atr.platType, atr.recommendType, atr.topicTitleLong, atr.topicTitleShort,
    atr.imgSrc, atr.topicURL, atr.recommendPosition, atr.recommendTime, atr.isFlag, atr.createTime, atr.createUserId,atr.tdate,avc.replyCount,avc.viewCount,tt.tagCode
    from autoTopicRecommend atr
    left join autotopicViewCount avc on atr.topicId = avc.topicId
    left join (
      select topicId,'1' topictype,tagCode from TopicTag
      union
      select topic_id,'2',tag_code from information_tag
      union
      select id,'0',tagdict from travel_cover
    ) tt on tt.topictype = atr.topicSrc and tt.topicId = atr.topicId
    where
      platType = 3 and
      recommendType = 32 and
      isFlag = '1' and
      recommendPosition = 3 and
      recommendTime <=now() and
    1=1
    order by recommendTime desc limit 0,1`, false},
		{`insert into travel_day (id, day_num, start_date,
        travel_id)
        values
          
        ('d498054fb73f57323ea84bd7a39f4daa', 1, '2016-07-27',
        485)
         , 
        ('678f4e7acfd272c2d312ea1cd70afad0', 2, '2016-07-28',
        485)
         , 
        ('086561b5541db2600d053d696d1e1636', 3, '2016-07-29',
        485)
         , 
        ('0466cedf06cda11fa7f9a90e933921a0', 4, '2016-07-30',
        485)`, false},
		{"SELECT count(1) from travel t where  t.user_id = '666666' and t.id = '1240'", false},
		{"delete from travel_content  where travel_id = 721", true},
		{"SELECT count(1) from travel t where  t.user_id = '666666' and t.id <= '1240'", false},
		{"delete from travel_content  where travel_id >= 721 and truncate(avg(a.full_price),2) = 2", true},
		{"select distinct series_id ,spec_id,city_id from car_price c where status = 1 and not exists (select member_id from car_salesman_info m where `status` = 1 and m.member_id= c.member_id ) and gmt_create > '2019-02-20 00:00:00' and gmt_create <= '2019-02-25 23:59:59' order by series_id limit 102000 ,1000", true},
		{"select distinct series_id ,spec_id,city_id from car_price c where status = 1 and exists (select member_id from car_salesman_info m where `status` = 1 and m.member_id= c.member_id ) and gmt_create > '2019-02-20 00:00:00' and gmt_create <= '2019-02-25 23:59:59' order by series_id limit 102000 ,1000", true},
		{"select max(a.naked_price) as npmax ,min(a.naked_price) as npmin ,truncate(avg(a.naked_price),2) npavg, max(a.full_price) as fpmax ,min(a.full_price) as fpmin ,truncate(avg(a.full_price),2) fpavg, count(1) as total from car_price a where 1=1 and a.is_hide = 0 and a.series_id=442 and a.`status` = 1", true},
		{`select high_priority host,user,password,select_priv,insert_priv,update_priv,delete_priv,create_priv,drop_priv,process_priv,grant_priv,references_priv,alter_priv,show_db_priv,super_priv,execute_priv,create_view_priv,show_view_priv,index_priv,create_user_priv,trigger_priv,create_role_priv,drop_role_priv,account_locked from mysql.user `,true},
	}

	for _, t := range table {

		fmt.Println("source >>>> : ", t.src)
		fmt.Println("       >>>> : ", ParserALLSQL(t.src))
	}

}

func ParserALLSQL(sql string) (absql string) {
	abSQL := ""
	parser := parser.New()
	// stmts, err := parser.Parse(src, "", "")
	stmt, err := parser.ParseOneStmt(sql, "", "")
	// fmt.Printf(" xx-1 :  %s  \n", reflect.TypeOf(stmt))
	if err != nil {
		fmt.Println("... err: ", err)
		return sql
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
	return abSQL
}

func AbstractUnionSQL(unn *ast.UnionStmt) (sql string) {
	tsql := ""
	// ditinct
	if unn.Distinct {
		tsql += "DISTINCT "
	}
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

// for update
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

// for delete
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

// for insert
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
		tsql = fmt.Sprintf("%s(...)*%v", tsql, len(inst.Lists))
	}
	return tsql
}

// for select
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
	//OrderBy
	if sel.OrderBy != nil {
		tsql += " OrderBy " + getByItemNames(sel.OrderBy.Items)
	}

	/* Having is the having condition.
	Having *HavingClause
	// OrderBy is the ordering expression list.
	OrderBy *OrderByClause
	// Limit is the limit clause.
	Limit *Limit
	*/

	if sel.Having != nil {
		tsql += " OrderBy " + getByItemNames(sel.OrderBy.Items)
	}

	if sel.Limit != nil {
		tsql += " LIMIT "
		fmt.Println(sel.Limit)
		if sel.Limit.Offset != 0 {
			tsql += "? "
		}
		if sel.Limit.Count != 0 {
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
			tmpName := oper.Expr.(*ast.ColumnNameExpr).Name
			if len(tmpName.Table.String()) > 0 {
				colName += tmpName.Table.String() + "."
			}
			colName += tmpName.Name.String()
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
			substr += "(" + AbstractSELSQL(oper.Sel.(*ast.SubqueryExpr).Query.(*ast.SelectStmt)) + ")"
		}
		if len(substr) > 0 {
			return fmt.Sprintf("%s %s%s", colName, instr, substr)
		} else {
			return fmt.Sprintf("%s %s(..*%v)", colName, instr, len(oper.List))
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
		op1 := oper.Op.String()
		ops_item, ok := ops[oper.Op.String()] // 假如key存在,则name = 李四 ，ok = true,否则，ok = false
		if ok {
			op1 = ops_item
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
			tmpName := oper.Expr.(*ast.ColumnNameExpr).Name
			if len(tmpName.Table.String()) > 0 {
				colName += tmpName.Table.String() + "."
			}
			colName += tmpName.Name.String()
		}
		return fmt.Sprintf("%s BETWEEN", colName)
	case *ast.PatternLikeExpr:
		colName := ""
		if oper.Expr != nil {
			tmpName := oper.Expr.(*ast.ColumnNameExpr).Name
			if len(tmpName.Table.String()) > 0 {
				colName += tmpName.Table.String() + "."
			}
			colName += tmpName.Name.String()
		}
		return fmt.Sprintf("%s LIKE", colName)
	case *ast.IsNullExpr:
		colName := ""
		if oper.Expr != nil {
			tmpName := oper.Expr.(*ast.ColumnNameExpr).Name
			if len(tmpName.Table.String()) > 0 {
				colName += tmpName.Table.String() + "."
			}
			colName += tmpName.Name.String()
		}
		if oper.Not {
			return fmt.Sprintf("%s IS NOT NULL", colName)
		}
		return fmt.Sprintf("%s IS NULL", colName)
	case *ast.FuncCallExpr:
		fname := oper.FnName
		return fmt.Sprintf("%s()", fname)
	case *ast.ValueExpr:

		return fmt.Sprintf("?")
	default:
		fmt.Println("default execution where ", reflect.TypeOf(oper))
	}
	return ""
}

func getByItemNames(bys []*ast.ByItem) (cols string) {
	names := ""
	for _, x := range bys {
		switch nod := x.Expr.(type) {
		case *ast.ColumnNameExpr:
			names += getColName(nod) + ","
		case *ast.ValueExpr:
			val := nod.GetValue()
			sval := fmt.Sprintf("%d", val)
			names += sval + ","
		case *ast.AggregateFuncExpr:
			names += nod.F + ","
			//		case *ast.FuncSubstringExpr:
			//			fmt.Printf("[[[[[3]]]] : %s %s  \n", nod, reflect.TypeOf(nod))
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
	default:
		fmt.Sprintf("%s", nod)
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
				fmt.Println("...xxx...", x.Schema)
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
				jstr = "LeftJoin"
			case 2:
				jstr = "RightJoin"
			default:
				jstr = "CrossJoin"
				fmt.Sprintf("%d", tp)
			}
			right := analyzeFrom(t.Right)
			return left + " " + jstr + " " + right
		}
		if t.On != nil {
			fmt.Printf(" on.Expr :  %s ", reflect.TypeOf(t.On.Expr))
		}
		return left
	default:
		return fmt.Sprintf("default exec ~~~~ %s ", reflect.TypeOf(t))
	}
	return ""
}

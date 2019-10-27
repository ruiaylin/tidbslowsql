package main

import (
	"fmt"
	"reflect"
	"strings"

	tidbTools "tisql/utils"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	driver "github.com/pingcap/tidb/types/parser_driver"
	// "reflect"
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

func main() {
	TestParser0()
}

// TestParser0 : table 43 and 50 parse failed
func TestParser0() {
	table := []struct {
		src string
		ok  bool
	}{
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
		{`insert into travel_day (id, day_num, start_date, travel_id)
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
		{"SELECT count(1) from travel t where  t.user_id = '666666' and t.id <= '1240'", false},
		{"update travel t  set x = 1 , y= 2 where  t.user_id = '666666' and t.id <= '1240'", false},
		{"select distinct series_id ,spec_id,city_id from car_price c where status = 1 and not exists (select member_id from car_salesman_info m where `status` = 1 and m.member_id= c.member_id ) and gmt_create > '2019-02-20 00:00:00' and gmt_create <= '2019-02-25 23:59:59' order by series_id limit 102000 ,1000", true},
		{"select max(a.naked_price) as npmax ,min(a.naked_price) as npmin ,truncate(avg(a.naked_price),2) npavg, max(a.full_price) as fpmax ,min(a.full_price) as fpmin ,truncate(avg(a.full_price),2) fpavg, count(1) as total from car_price a where 1=1 and a.is_hide = 0 and a.series_id=442 and a.`status` = 1", true},
	}

	for _, t := range table {
		fmt.Println("source >>>> : ", t.src)
		sql, _ := tidbTools.ParserALLSQL(t.src)
		fmt.Println("       >>>> 01 : ", sql)
		normalized, digest := parser.NormalizeDigest(t.src)
		fmt.Println("       >>>> 02 : ", normalized, digest)
	}

}

// ParserALLSQL entry of sql parse
func ParserALLSQL1(sql string) (absql string) {
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
	//OrderBy
	if sel.OrderBy != nil {
		tsql += " OrderBy " + getByItemNames(sel.OrderBy.Items)
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
		// arg
		var argss []string
		for y, x := range oper.Args {
			fmt.Println(" y = ", y, x, " type = ", reflect.TypeOf(x))
			switch nod := x.(type) {
			case *ast.ColumnNameExpr:
				argss = append(argss, getColName(nod))
			case *ast.AggregateFuncExpr:
				argss = append(argss, nod.F+"(?)")
			case *driver.ValueExpr:
				argss = append(argss, "?")
			default:
				fmt.Sprintf("[[[[[4]]]] : %s %s  \n", nod, reflect.TypeOf(nod))
			}
		}
		argss1 := strings.Join(argss, ",")
		return fmt.Sprintf("%s(%s)", fname, argss1)
	case *ast.ValueExpr:
		return fmt.Sprintf("?")
	case *driver.ValueExpr:
		return fmt.Sprintf("?")
	case *ast.ExistsSubqueryExpr:
		return genExistString(oper)
	case *ast.ColumnNameExpr:
		name := oper.Name
		return fmt.Sprintf("%s", name)
	default:
		fmt.Println("default execution where ", reflect.TypeOf(oper))
	}
	return ""
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
	for y, x := range node.Args {
		fmt.Println(" y = ", y, x, " type = ", reflect.TypeOf(x))
		switch nod := x.(type) {
		case *ast.ColumnNameExpr:
			argss = append(argss, getColName(nod))
		case *driver.ValueExpr:
			argss = append(argss, "?")
		default:
			fmt.Sprintf("[[[[[4]]]] : %s %s  \n", nod, reflect.TypeOf(nod))
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
	default:
		fmt.Sprintf("%s", nod)
		//fmt.Println("typeCheck check ", reflect.TypeOf(node))
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

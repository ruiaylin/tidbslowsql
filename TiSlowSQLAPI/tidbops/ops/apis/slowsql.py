# from .auth import auth
from ops.models import db
from flask import jsonify
from flask import current_app as app
from sqlalchemy.sql import text
# from sqlalchemy.sql import func
from datetime import datetime, timedelta
from flask_restful import Resource, reqparse
from ops.models.slowsql import SlowQueryInfo,\
    SlowQueryInfoViewSchema, SlowQuery, SlowSql
from ops.utils.slowsql import get_slowsql_rpt,\
    get_realtime_slowsql_rpt, get_tidb_slow_query_info


# for event_API
class SlowSQLAPI(Resource):
    """
        docstring for SlowSQLAPI
    """

    def __init__(self):
        # for query
        self.getparse = reqparse.RequestParser()
        self.getparse.add_argument(
            'cluster',
            type=str,
            required=True)
        self.getparse.add_argument(
            'instName',
            type=str,
            required=False)
        self.getparse.add_argument(
            'startTime', type=str,
            required=False)
        self.getparse.add_argument(
            'endTime', type=str,
            required=False)
        self.getparse.add_argument(
            'pageSize', type=int,
            required=False)
        self.getparse.add_argument(
            'page', type=int,
            required=False)
        self.putparse = reqparse.RequestParser()
        self.putparse.add_argument(
            'cluster',
            type=str, location="json",
            required=True)
        self.putparse.add_argument(
            'sqlID',
            type=str, location="json",
            required=True)
        self.putparse.add_argument(
            'affirmPerson',
            type=str, location="json",
            required=False, default='')
        super(SlowSQLAPI, self).__init__()

    def get(self):
        RSYES = {'code': 1, 'msg': 'Successful.'}
        # RSNO = {'code': 0, 'msg': 'Fail'}
        data = {}
        data['status'] = RSYES
        data['results'] = []
        argss = self.getparse.parse_args()
        cluster = argss.get('cluster')
        inst_name = argss.get('instName')
        st = argss.get('startTime')
        et = argss.get('endTime')
        page = argss.get('page')
        PAGE_SIZE = argss.get('pageSize')
        if not page:
            page = 1
        if not PAGE_SIZE:
            PAGE_SIZE = 50
        if not st:
            stime = datetime.now() - timedelta(hours=2)
            etime = datetime.now()
            stime = datetime.strptime(
                datetime.strftime(
                    stime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(
                datetime.strftime(
                    etime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
        else:
            stime = datetime.strptime(st, '%Y-%m-%d %H:%M:%S')
            if et is None:
                et = datetime.strftime(datetime.now(), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(et, '%Y-%m-%d %H:%M:%S')
        if not inst_name:
            condition_time = (
                SlowQueryInfo.cluster == SlowQuery.cluster,
                SlowQueryInfo.cluster == SlowSql.cluster,
                SlowQueryInfo.day_time >= stime,
                SlowQueryInfo.day_time <= etime,
                SlowQueryInfo.cluster == cluster
            )
        else:
            condition_time = (
                SlowQueryInfo.cluster == SlowQuery.cluster,
                SlowQueryInfo.cluster == SlowSql.cluster,
                SlowQueryInfo.day_time >= stime,
                SlowQueryInfo.day_time <= etime,
                SlowQueryInfo.cluster == cluster,
                SlowQueryInfo.inst_name == inst_name
            )
        querysqllist = SlowQueryInfo.query.join(
            SlowQuery, SlowQueryInfo.query_id == SlowQuery.query_id
        ).join(
            SlowSql, SlowQueryInfo.sql_id == SlowSql.sql_id
        ).add_columns(
            SlowQueryInfo.id,
            SlowQueryInfo.sql_id,
            SlowQueryInfo.query_id,
            SlowQueryInfo.cluster,
            SlowQueryInfo.hostname,
            SlowQueryInfo.port,
            SlowQueryInfo.inst_name,
            SlowQueryInfo.day_time,
            SlowQueryInfo.schema_name,
            SlowQueryInfo.table_name,
            SlowQueryInfo.start_ts,
            SlowQueryInfo.username,
            SlowQueryInfo.client_host,
            SlowQuery.sql_text,
            SlowQuery.sql_fulltext,
            SlowQuery.counter,
            SlowQueryInfo.query_time,
            SlowQueryInfo.process_time,
            SlowQueryInfo.wait_time,
            SlowQueryInfo.get_commit_ts_time,
            SlowQueryInfo.commit_time,
            SlowQueryInfo.mem_max,
            SlowQueryInfo.total_keys,
            SlowQueryInfo.process_keys,
            SlowQueryInfo.write_keys,
            SlowQueryInfo.write_size,
            SlowSql.first_seen,
            SlowSql.affirm,
            SlowQueryInfo.is_internal
        ).filter(
            *condition_time
        ).order_by(
            SlowQueryInfo.day_time.desc()
        ).paginate(page, PAGE_SIZE, False)
        queryinfo_schemas = SlowQueryInfoViewSchema(many=True)
        yz = queryinfo_schemas.dump(querysqllist.items)
        # print(yz.data)
        data['results'] = yz.data
        data['total'] = querysqllist.total
        data['pages'] = querysqllist.pages
        data['page'] = page
        if querysqllist.has_prev:
            data['prev_num'] = querysqllist.prev_num
        if querysqllist.has_next:
            data['next_num'] = querysqllist.next_num
        return jsonify(data)

    def put(self):
        RSYES = {'code': 1, 'msg': 'Successful.'}
        RSNO = {'code': 0, 'msg': 'Fail'}
        data = {}
        data['status'] = RSYES
        argss = self.putparse.parse_args()
        cluster = argss.get('cluster')
        sql_id = argss.get('sqlID')
        affirm_person = argss.get('affirmPerson')
        # 慢查询SQL
        param = {}
        param['cluster'] = cluster
        param['sql_id'] = sql_id
        param['affirm_person'] = affirm_person
        slowsql = '''
            update tidb_slow_sql
            set affirm = 1,
            affirm_person = :affirm_person
            where
                cluster = :cluster
                and sql_id = :sql_id
        '''
        try:
            eg = db.get_engine(app, 'tidbops')
            eg.execute(text(slowsql), **param)
            db.session.commit()
            return jsonify(data)
        except Exception as e:
            db.session.rollback()
            data['status'] = RSNO
            data['info'] = str(e)
            return jsonify(data)


# for event_API
class SlowSqlTrendAPI(Resource):
    """
        docstring for SlowReportAPI
    """

    def __init__(self):
        # for query
        self.getparse = reqparse.RequestParser()
        self.getparse.add_argument(
            'cluster',
            type=str,
            required=True)
        self.getparse.add_argument(
            'instName',
            type=str,
            required=False)
        self.getparse.add_argument(
            'startTime',
            type=str)
        self.getparse.add_argument(
            'endTime',
            type=str,
            required=False)
        super(SlowSqlTrendAPI, self).__init__()

    def get(self):
        RSYES = {'code': 1, 'msg': 'Successful.'}
        # RSNO = {'code': 0, 'msg': 'Fail'}
        data = {}
        data['status'] = RSYES
        data['results'] = []
        argss = self.getparse.parse_args()
        cluster = argss.get('cluster')
        inst_name = argss.get('instName')
        st = argss.get('startTime')
        et = argss.get('endTime')
        if not st:
            stime = datetime.now() - timedelta(hours=2)
            etime = datetime.now()
            stime = datetime.strptime(
                datetime.strftime(
                    stime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(
                datetime.strftime(
                    etime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
        else:
            stime = datetime.strptime(st, '%Y-%m-%d %H:%M:%S')
            if et is None:
                et = datetime.strftime(datetime.now(), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(et, '%Y-%m-%d %H:%M:%S')
        # 慢查询SQL
        param = {}
        param['cluster'] = cluster
        param['start_time'] = stime
        param['end_time'] = etime
        slowsql = '''
            select
                count(*) as counter,
                v_day_time as day_time

            from tidb_slow_query_info
            where
                cluster = :cluster
                and day_time >= :start_time
                and day_time <= :end_time
        '''
        if inst_name:
            param['inst_name'] = inst_name
            slowsql += " and inst_name = :inst_name"
        slowsql += " group by v_day_time"
        eg = db.get_engine(app, 'tidbops')
        slow_set = eg.execute(text(slowsql), **param)
        slow_list = []
        for row in slow_set.fetchall():
            tmp = {
                'counter': row[0],
                'day_time': datetime.strftime(row[1], '%Y-%m-%d %H:%M:%S'),
                'cluster': cluster
            }
            if inst_name:
                tmp['inst_name'] = inst_name
            slow_list.append(tmp)
        slow_set.close()
        data['results'] = slow_list
        return jsonify(data)


# for event_API
class SlowQueryTrendAPI(Resource):
    """
        docstring for SlowReportAPI
    """

    def __init__(self):
        # for query
        self.getparse = reqparse.RequestParser()
        self.getparse.add_argument(
            'cluster',
            type=str,
            required=True)
        self.getparse.add_argument(
            'sqlID',
            type=str,
            required=True)
        self.getparse.add_argument(
            'startTime',
            type=str)
        self.getparse.add_argument(
            'endTime',
            type=str,
            required=False)
        super(SlowQueryTrendAPI, self).__init__()

    def get(self):
        RSYES = {'code': 1, 'msg': 'Successful.'}
        # RSNO = {'code': 0, 'msg': 'Fail'}
        data = {}
        data['status'] = RSYES
        data['results'] = []
        argss = self.getparse.parse_args()
        cluster = argss.get('cluster')
        sql_id = argss.get('sqlID')
        st = argss.get('startTime')
        et = argss.get('endTime')
        if not st:
            stime = datetime.now() - timedelta(hours=2)
            etime = datetime.now()
            stime = datetime.strptime(
                datetime.strftime(
                    stime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(
                datetime.strftime(
                    etime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
        else:
            stime = datetime.strptime(st, '%Y-%m-%d %H:%M:%S')
            if et is None:
                et = datetime.strftime(datetime.now(), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(et, '%Y-%m-%d %H:%M:%S')
        # 慢查询SQL
        param = {}
        param['cluster'] = cluster
        param['sql_id'] = sql_id
        param['start_time'] = stime
        param['end_time'] = etime
        slowsql = '''
            select
                count(*) as counter,
                round(avg(query_time), 4) query_time,
                round(max(query_time), 4) max_query_time,
                v_day_time as day_time
            from tidb_slow_query_info
            where
                cluster = :cluster
                and sql_id = :sql_id
                and day_time >= :start_time
                and day_time <= :end_time
                group by v_day_time
        '''
        eg = db.get_engine(app, 'tidbops')
        slow_set = eg.execute(text(slowsql), **param)
        print(slowsql)
        print(param)
        slow_list = []
        for row in slow_set.fetchall():
            tmp = {
                'cluster': cluster,
                'counter': row[0],
                'query_time': row[1],
                'max_query_time': row[2],
                'day_time': datetime.strftime(row[3], '%Y-%m-%d %H:%M:%S')
            }
            slow_list.append(tmp)
        slow_set.close()
        data['results'] = slow_list
        return jsonify(data)


# for event_API
class SlowReportAPI(Resource):
    """
        docstring for SlowReportAPI
    """

    def __init__(self):
        # for query
        self.getparse = reqparse.RequestParser()
        self.getparse.add_argument(
            'cluster',
            type=str,
            required=True)
        self.getparse.add_argument(
            'instName',
            type=str,
            required=False)
        self.getparse.add_argument(
            'startTime',
            type=str,
            required=True)
        self.getparse.add_argument(
            'endTime',
            type=str,
            required=True)
        self.getparse.add_argument(
            'type',
            type=str,
            required=False, default='minute')
        super(SlowReportAPI, self).__init__()

    def get(self):
        RSYES = {'code': 1, 'msg': 'Successful.'}
        # RSNO = {'code': 0, 'msg': 'Fail'}
        data = {}
        data['status'] = RSYES
        data['results'] = []
        argss = self.getparse.parse_args()
        cluster = argss.get('cluster')
        inst_name = argss.get('instName')
        st = argss.get('startTime')
        et = argss.get('endTime')
        rpt_type = argss.get('type')
        if not st:
            rpt_type = "hour"
            stime = datetime.now() - timedelta(hours=5)
            etime = datetime.now()
            stime = datetime.strptime(
                datetime.strftime(
                    stime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(
                datetime.strftime(
                    etime, '%Y-%m-%d %H:%M:%S'), '%Y-%m-%d %H:%M:%S')
            stime = datetime.strftime(stime, '%Y-%m-%d %H:00:00')
            etime = datetime.strftime(etime, '%Y-%m-%d %H:00:00')
        else:
            stime = datetime.strptime(st, '%Y-%m-%d %H:%M:%S')
            if et is None:
                et = datetime.strftime(datetime.now(), '%Y-%m-%d %H:%M:%S')
            etime = datetime.strptime(et, '%Y-%m-%d %H:%M:%S')
            seconds_delta = (etime - stime).total_seconds()
            if seconds_delta <= 3600:
                rpt_type = "minute"
                stime = datetime.strftime(stime, '%Y-%m-%d %H:%M:%S')
                etime = datetime.strftime(etime, '%Y-%m-%d %H:%M:%S')
            elif seconds_delta > 3600 and seconds_delta <= 86400:
                rpt_type = "hour"
                stime = datetime.strftime(stime, '%Y-%m-%d %H:00:00')
                etime = datetime.strftime(etime, '%Y-%m-%d %H:00:00')
            elif seconds_delta > 86400 and seconds_delta <= 604800:
                rpt_type = "daily"
                stime = datetime.strftime(stime, '%Y-%m-%d 00:00:00')
                etime = datetime.strftime(etime, '%Y-%m-%d 00:00:00')
            else:
                rpt_type = "daily"
                stime = etime - timedelta(hours=168)
                stime = datetime.strftime(stime, '%Y-%m-%d 00:00:00')
                etime = datetime.strftime(etime, '%Y-%m-%d 00:00:00')
        # 慢查询SQL
        if rpt_type != 'minute':
            slow_list = get_slowsql_rpt(cluster,
                                        inst_name, stime, etime, rpt_type)
        else:
            slow_list = get_realtime_slowsql_rpt(cluster,
                                                 inst_name, stime, etime)
        data['results'] = slow_list
        return jsonify(data)


# for event_API
class SlowInfoAPI(Resource):
    """
        docstring for SlowReportAPI
    """

    def __init__(self):
        # for query
        self.getparse = reqparse.RequestParser()
        self.getparse.add_argument(
            'cluster',
            type=str,
            required=True)
        self.getparse.add_argument(
            'queryID',
            type=str,
            required=True)
        super(SlowInfoAPI, self).__init__()

    def get(self):
        RSYES = {'code': 1, 'msg': 'Successful.'}
        # RSNO = {'code': 0, 'msg': 'Fail'}
        data = {}
        data['status'] = RSYES
        data['results'] = []
        argss = self.getparse.parse_args()
        cluster = argss.get('cluster')
        query_id = argss.get('queryID')
        # 慢查询SQL
        query_info = get_tidb_slow_query_info(cluster, query_id)
        data['results'] = query_info
        return jsonify(data)

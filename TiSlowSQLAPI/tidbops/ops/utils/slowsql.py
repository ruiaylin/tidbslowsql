from ops.models import db
# from ops.models.base import Instance
# from sqlalchemy.sql import text
from flask import current_app as app
# from datetime import datetime, timedelta
from ops.utils.common import get_mysql_connection
from datetime import datetime, timedelta


# 获取实时慢sql报表
def get_realtime_slowsql_rpt(cluster, inst, st, et,):
    try:
        eg = db.get_engine(app, 'tidbops')
        select_inst = ""
        if inst:
            groupby = """ and a.cluster='%s' and a.inst_name='%s'
            GROUP BY a.cluster,a.inst_name,a.sql_id
            """ % (cluster, inst)
            select_inst = ",a.inst_name"
        else:
            groupby = """ and a.cluster='%s' GROUP BY a.cluster,a.sql_id
            """ % (cluster)
        sc_list = []
        slowsql = """
            SELECT
            a.sql_id,
            a.cluster,
            count(1) AS count,
            b.first_seen,
            round(avg(a.query_time), 4) query_time,
            round(max(a.query_time), 4) max_query_time,
            round(avg(a.process_time), 4) process_time,
            round(max(a.process_time), 4) max_process_time,
            round(avg(a.wait_time), 4) wait_time,
            round(max(a.wait_time), 4) max_wait_time,
            round(avg(a.get_commit_ts_time), 4) get_commit_ts_time,
            round(max(a.get_commit_ts_time), 4) max_get_commit_ts_time,
            round(avg(a.commit_time), 4) commit_time,
            round(max(a.commit_time), 4) max_commit_time,
            round(avg(a.total_keys), 0) total_keys,
            round(max(a.total_keys), 0) max_total_keys,
            round(avg(a.process_keys), 0) process_keys,
            round(max(a.process_keys), 0) max_process_keys,
            round(avg(a.mem_max), 0) mem_max,
            round(max(a.mem_max), 0) max_mem_max,
            round(avg(a.write_keys), 0) write_keys,
            round(max(a.write_keys), 0) max_write_keys,
            round(avg(a.write_size), 0) write_size,
            round(max(a.write_size), 0) max_write_size,
            round(max(a.cop_proc_max), 3) cop_proc_max,
            round(max(a.cop_wait_max), 3) cop_wait_max,
            a.query_id,b.affirm,b.sql_text,a.schema_name
            %s
            FROM
            tidb_slow_query_info a inner join tidb_slow_sql b
            on a.sql_id=b.sql_id and a.cluster=b.cluster
            WHERE a.day_time BETWEEN '%s' AND '%s' %s order by count desc;
        """ % (select_inst, st, et, groupby)
        print(slowsql)
        slowrpt_set = eg.execute(slowsql)
        for row in slowrpt_set.fetchall():
            tmp = {
                'sql_id': row[0],
                'cluster': row[1],
                'count': row[2],
                'first_seen': datetime.strftime(
                    row[3], '%Y-%m-%d %H:%M:%S'),
                'query_time': row[4],
                'max_query_time': row[5],
                'process_time': row[6],
                'max_process_time': row[7],
                'wait_time': row[8],
                'max_wait_time': row[9],
                'get_commit_ts_time': row[10],
                'max_get_commit_ts_time': row[11],
                'commit_time': row[12],
                'max_commit_time': row[13],
                'total_keys': row[14],
                'max_total_keys': row[15],
                'process_keys': row[16],
                'max_process_keys': row[17],
                'mem_max': row[18],
                'max_mem_max': row[19],
                'write_keys': row[20],
                'max_write_keys': row[21],
                'write_size': row[22],
                'max_write_size': row[23],
                'cop_proc_max': row[24],
                'cop_wait_max': row[25],
                'query_id': row[26],
                'affirm': row[27],
                'sql_text': row[28],
                'schema_name': row[29]
            }
            if inst:
                tmp['inst_name'] = row[30]
            sc_list.append(tmp)
        slowrpt_set.close()
        return sc_list
    except Exception as e:
        raise e


# 获取慢sql报表
def get_slowsql_rpt(cluster, inst, st, et, rpt_type):
    try:
        eg = db.get_engine(app, 'tidbops')
        select_inst = ""
        if inst:
            groupby = """ and cluster='%s' and inst_name='%s' group by cluster, inst_name, sql_id
            """ % (cluster, inst)
            select_inst = ",inst_name"
            if rpt_type == "daily":
                tab_name = "tidb_inst_slow_sql_daily_rpt"
            else:
                tab_name = "tidb_inst_slow_sql_hour_rpt"
        else:
            groupby = """ and cluster='%s' group by cluster, sql_id
            """ % (cluster)
            if rpt_type == "daily":
                tab_name = "tidb_cluster_slow_sql_daily_rpt"
            else:
                tab_name = "tidb_cluster_slow_sql_hour_rpt"
        sc_list = []
        slowsql = '''
            SELECT
            sql_id,
            cluster,
            sum(count) AS count,
            first_seen,
            round(avg(query_time), 4) query_time,
            round(max(max_query_time), 4) max_query_time,
            round(avg(process_time), 4) process_time,
            round(max(max_process_time), 4) max_process_time,
            round(avg(wait_time), 4) wait_time,
            round(max(max_wait_time), 4) max_wait_time,
            round(avg(get_commit_ts_time), 4) get_commit_ts_time,
            round(max(max_get_commit_ts_time), 4) max_get_commit_ts_time,
            round(avg(commit_time), 4) commit_time,
            round(max(max_commit_time), 4) max_commit_time,
            round(avg(total_keys), 0) total_keys,
            round(max(max_total_keys), 0) max_total_keys,
            round(avg(process_keys), 0) process_keys,
            round(max(max_process_keys), 0) max_process_keys,
            round(avg(mem_max), 0) mem_max,
            round(max(max_mem_max), 0) max_mem_max,
            round(avg(write_keys), 0) write_keys,
            round(max(max_write_keys), 0) max_write_keys,
            round(avg(write_size), 0) write_size,
            round(max(max_write_size), 0) max_write_size,
            round(max(cop_proc_max), 3) cop_proc_max,
            round(max(cop_wait_max), 3) cop_wait_max,query_id,affirm,sql_text,
            schema_name
            %s
            FROM
            %s
            WHERE
            daytime BETWEEN '%s' AND '%s' %s order by count desc
        ''' % (select_inst, tab_name, st, et, groupby)
        print(slowsql)
        slowrpt_set = eg.execute(slowsql)
        for row in slowrpt_set.fetchall():
            tmp = {
                'sql_id': row[0],
                'cluster': row[1],
                'count': row[2],
                'first_seen': datetime.strftime(
                    row[3], '%Y-%m-%d %H:%M:%S'),
                'query_time': row[4],
                'max_query_time': row[5],
                'process_time': row[6],
                'max_process_time': row[7],
                'wait_time': row[8],
                'max_wait_time': row[9],
                'get_commit_ts_time': row[10],
                'max_get_commit_ts_time': row[11],
                'commit_time': row[12],
                'max_commit_time': row[13],
                'total_keys': row[14],
                'max_total_keys': row[15],
                'process_keys': row[16],
                'max_process_keys': row[17],
                'mem_max': row[18],
                'max_mem_max': row[19],
                'write_keys': row[20],
                'max_write_keys': row[21],
                'write_size': row[22],
                'max_write_size': row[23],
                'cop_proc_max': row[24],
                'cop_wait_max': row[25],
                'query_id': row[26],
                'affirm': row[27],
                'sql_text': row[28],
                'schema_name': row[29]
            }
            if inst:
                tmp['inst_name'] = row[30]
            sc_list.append(tmp)
        slowrpt_set.close()
        return sc_list
    except Exception as e:
        raise e


def get_tidb_slow_query_info(cluster, query_id):
    sql = """
    select a.sql_text, a.sql_fulltext, b.cluster, b.inst_name,\
    b.start_ts, b.sql_id, b.query_id,b.day_time,\
    b.schema_name, b.table_name, b.username, b.client_host,\
    b.query_time, b.process_time, b.wait_time, b.get_commit_ts_time,\
    b.commit_time,b.total_keys, b.process_keys, b.mem_max, b.write_keys,\
    b.write_size,b.cop_proc_max,b.cop_wait_max,b.mvcc_data,\
    b.backoff_time,b.request_count,b.prewrite_time,b.total_backoff_time,\
    b.latch_wait_time,b.resolve_lock_time,b.prewrite_region,b.trx_retry,\
    b.index_ids,b.is_internal,b.digest,b.stats,b.cop_proc_avg,b.cop_proc_90,\
    b.cop_proc_addr,b.cop_wait_avg,b.cop_wait_90,b.cop_wait_addr,b.num_cop_tasks,\
    b.succ from tidb_slow_query a, tidb_slow_query_info b \
    where a.cluster=b.cluster and \
    a.query_id= b.query_id and \
    b.cluster='%s' and b.query_id=%s order by b.day_time desc limit 1
    """ % (cluster, query_id)
    queryinfo = {}
    try:
        eg = db.get_engine(app, 'tidbops')
        query_set = eg.execute(sql)
        for sql_text, sql_fulltext, cluster, inst_name, \
                start_ts, sql_id, query_id, \
                day_time, schema_name, table_name, \
                username, client_host, \
                query_time, process_time, \
                wait_time, get_commit_ts_time, \
                commit_time, total_keys, \
                process_keys, mem_max, write_keys, \
                write_size, cop_proc_max, cop_wait_max, \
                mvcc_data, backoff_time, request_count,\
                prewrite_time, total_backoff_time,\
                latch_wait_time, resolve_lock_time, \
                prewrite_region, trx_retry,\
                index_ids, is_internal, digest, stats, \
                cop_proc_avg, cop_proc_90,\
                cop_proc_addr, cop_wait_avg, cop_wait_90, \
                cop_wait_addr, num_cop_tasks,\
                succ in query_set.fetchall():
            queryinfo['sql_text'] = sql_text
            queryinfo['sql_fulltext'] = sql_fulltext
            queryinfo['cluster'] = cluster
            queryinfo['inst_name'] = inst_name
            queryinfo['start_ts'] = start_ts
            queryinfo['sql_id'] = sql_id
            queryinfo['query_id'] = query_id
            queryinfo['day_time'] = datetime.strftime(
                day_time, '%Y-%m-%d %H:%M:%S')
            queryinfo['schema_name'] = schema_name
            queryinfo['table_name'] = table_name
            queryinfo['username'] = username
            queryinfo['client_host'] = client_host
            queryinfo['query_time'] = query_time
            queryinfo['process_time'] = process_time
            queryinfo['wait_time'] = wait_time
            queryinfo['get_commit_ts_time'] = get_commit_ts_time
            queryinfo['commit_time'] = commit_time
            queryinfo['total_keys'] = total_keys
            queryinfo['process_keys'] = process_keys
            queryinfo['mem_max'] = mem_max
            queryinfo['write_keys'] = write_keys
            queryinfo['write_size'] = write_size
            queryinfo['cop_proc_max'] = cop_proc_max
            queryinfo['cop_wait_max'] = cop_wait_max
            queryinfo['mvcc_data'] = mvcc_data
            queryinfo['backoff_time'] = backoff_time
            queryinfo['request_count'] = request_count
            queryinfo['prewrite_time'] = prewrite_time
            queryinfo['total_backoff_time'] = total_backoff_time
            queryinfo['latch_wait_time'] = latch_wait_time
            queryinfo['resolve_lock_time'] = resolve_lock_time
            queryinfo['prewrite_region'] = prewrite_region
            queryinfo['trx_retry'] = trx_retry
            queryinfo['index_ids'] = index_ids
            queryinfo['is_internal'] = is_internal
            queryinfo['digest'] = digest
            queryinfo['stats'] = stats
            queryinfo['cop_proc_avg'] = cop_proc_avg
            queryinfo['cop_proc_90'] = cop_proc_90
            queryinfo['cop_proc_addr'] = cop_proc_addr
            queryinfo['cop_wait_avg'] = cop_wait_avg
            queryinfo['cop_wait_90'] = cop_wait_90
            queryinfo['cop_wait_addr'] = cop_wait_addr
            queryinfo['num_cop_tasks'] = num_cop_tasks
            queryinfo['succ'] = succ
        query_set.close()
        if 'inst_name' in queryinfo.keys():
            queryinfo['plan'] = get_tidb_query_explain(
                inst_name, schema_name, sql_fulltext)
        return queryinfo
    except Exception as e:
        print(e)
        raise e


def get_tidb_query_explain(inst_name, db_name, sql_text):
    data = []
    # get target conn
    try:
        conn = get_mysql_connection(inst_name, db_name)
    except Exception as e:
        print(e)
        raise e
    # get explain info
    try:
        cur = conn.cursor()
        cur.execute('desc %s' % sql_text)
        for x in cur:
            dict = {}
            dict['id'] = x[0]
            dict['count'] = x[1]
            dict['task'] = x[2]
            dict['operator_info'] = x[3]
            data.append(dict)
        conn.commit()
        return data
    except Exception as e:
        conn.rollback()
        raise e
    finally:
        cur.close()
        conn.close()


def get_stime_etime(st=None, et=None):
    if not st:
        stime = datetime.now() - timedelta(hours=5)
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
    return stime, etime, st, et

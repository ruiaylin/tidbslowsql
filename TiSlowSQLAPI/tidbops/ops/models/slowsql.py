# coding: utf-8
from datetime import datetime
from ..models import db, ma


# slow sql information
class SlowSql(db.Model):
    """slow sql informations for SlowSql"""
    __tablename__ = 'tidb_slow_sql'
    __bind_key__ = 'tidbops'
    id = db.Column(db.Integer, primary_key=True)
    cluster = db.Column(db.String(32), index=True)
    sql_id = db.Column(db.Integer, nullable=False)
    sql_text = db.Column(db.TEXT, nullable=False)
    first_seen = db.Column(db.DateTime, default=datetime.now)
    affirm = db.Column(db.Integer)
    affirm_person = db.Column(db.String(128), nullable=False, index=True)
    counter = db.Column(db.Integer, nullable=True)
    update_time = db.Column(
        db.DateTime,
        default=datetime.now,
    )

    __table_args__ = (
        db.Index(
            'idx_slow_sql_count_updateTime',
            'counter', 'update_time'
        ),
    )

    def __str__(self):
        return self.sql_id

    def __repr__(self):
        return '<SlowSql: %s>' % self.id

    def save(self):
        # update slow sql info & insert
        db.session.add(self)
        db.session.commit()
        return self


# slow query information
class SlowQuery(db.Model):
    """slow query informations for SlowQuery"""
    __tablename__ = 'tidb_slow_query'
    __bind_key__ = 'tidbops'
    id = db.Column(db.Integer, primary_key=True)
    cluster = db.Column(db.String(32), index=True)
    host = db.Column(db.String(32), index=True)
    port = db.Column(db.Integer, nullable=False)
    inst_name = db.Column(db.String(32), index=True)
    sql_id = db.Column(db.Integer, nullable=False)
    query_id = db.Column(db.Integer, nullable=False)
    sql_text = db.Column(db.TEXT, nullable=False)
    sql_fulltext = db.Column(db.TEXT, nullable=False)
    counter = db.Column(db.Integer, nullable=True)
    createTime = db.Column(db.DateTime, default=datetime.now)
    updateTime = db.Column(
        db.DateTime,
        default=datetime.now,
        index=True
    )

    __table_args__ = (
        db.Index(
            'idx_sqlid_instname',
            'sql_id', 'inst_name'
        ),
    )

    def __str__(self):
        return self.query_id

    def __repr__(self):
        return '<SlowQuery: %s>' % self.id

    def save(self):
        # update slow query info & insert
        db.session.add(self)
        db.session.commit()
        return self


# slow query detail information
class SlowQueryInfo(db.Model):
    """slow query detail informations for SlowQueryInfo"""
    __tablename__ = 'tidb_slow_query_info'
    __bind_key__ = 'tidbops'
    id = db.Column(db.Integer)
    hostname = db.Column(db.String(32), nullable=False)
    ip = db.Column(db.String(24), index=True)
    port = db.Column(db.Integer, nullable=False)
    inst_name = db.Column(db.String(32), index=True)
    day_time = db.Column(db.DateTime)
    cluster = db.Column(db.String(32), index=True)
    sql_id = db.Column(db.Integer, nullable=False)
    query_id = db.Column(db.Integer, nullable=False)
    schema_name = db.Column(db.String(64))
    table_name = db.Column(db.String(128))
    start_ts = db.Column(db.Integer)
    username = db.Column(db.String(32))
    client_host = db.Column(db.String(16))
    conn_id = db.Column(db.Integer)
    query_time = db.Column(db.Float)
    process_time = db.Column(db.Float)
    wait_time = db.Column(db.Float)
    backoff_time = db.Column(db.Float)
    request_count = db.Column(db.Integer)
    total_keys = db.Column(db.Integer)
    process_keys = db.Column(db.Integer)
    prewrite_time = db.Column(db.Float)
    commit_time = db.Column(db.Float)
    get_commit_ts_time = db.Column(db.Float)
    total_backoff_time = db.Column(db.Float)
    latch_wait_time = db.Column(db.Float)
    resolve_lock_time = db.Column(db.Float)
    write_keys = db.Column(db.Integer)
    write_size = db.Column(db.Integer)
    prewrite_region = db.Column(db.Integer)
    trx_retry = db.Column(db.Integer)
    index_ids = db.Column(db.String(128), nullable=True)
    is_internal = db.Column(db.String(10), nullable=True)
    digest = db.Column(db.String(64), nullable=True)
    stats = db.Column(db.String(128), nullable=True)
    cop_proc_avg = db.Column(db.Float)
    cop_proc_90 = db.Column(db.Float)
    cop_proc_max = db.Column(db.Float)
    cop_proc_addr = db.Column(db.String(32))
    cop_wait_avg = db.Column(db.Float)
    cop_wait_90 = db.Column(db.Float)
    cop_wait_max = db.Column(db.Float)
    cop_wait_addr = db.Column(db.String(32))
    num_cop_tasks = db.Column(db.Integer, nullable=True)
    mem_max = db.Column(db.Integer, nullable=True)
    succ = db.Column(db.String(30), nullable=True)
    create_time = db.Column(db.DateTime, default=datetime.now)

    __table_args__ = (
        db.PrimaryKeyConstraint(
            'id', 'day_time', name='pk_tidb_slow_query_info'
        ),
        db.Index(
            'idx_tidb_query_cluster_dt',
            'cluster', 'day_time'
        ),
    )

    def __str__(self):
        return self.query_id

    def __repr__(self):
        return '<slow_query_info: %s>' % self.id

    def save(self):
        # update slow query detail info & insert
        db.session.add(self)
        db.session.commit()
        return self


class SlowQueryInfoView(db.Model):
    """slow query detail informations for SlowQueryInfo"""
    __tablename__ = 'tidb_slow_query_info_view'
    __bind_key__ = 'tidbops'
    id = db.Column(db.Integer, primary_key=True)
    hostname = db.Column(db.String(32), nullable=False)
    port = db.Column(db.Integer, nullable=False)
    inst_name = db.Column(db.String(32), index=True)
    day_time = db.Column(db.DateTime)
    cluster = db.Column(db.String(32), index=True)
    sql_id = db.Column(db.Integer, nullable=False)
    query_id = db.Column(db.Integer, nullable=False)
    schema_name = db.Column(db.String(64))
    table_name = db.Column(db.String(128))
    username = db.Column(db.String(32))
    client_host = db.Column(db.String(16))
    sql_text = db.Column(db.TEXT, nullable=False)
    sql_fulltext = db.Column(db.TEXT, nullable=False)
    counter = db.Column(db.Integer, nullable=True)
    query_time = db.Column(db.Float)
    process_time = db.Column(db.Float)
    wait_time = db.Column(db.Float)
    total_keys = db.Column(db.Integer)
    process_keys = db.Column(db.Integer)
    commit_time = db.Column(db.Float)
    get_commit_ts_time = db.Column(db.Float)
    write_keys = db.Column(db.Integer)
    write_size = db.Column(db.Integer)
    mem_max = db.Column(db.Integer, nullable=True)
    is_internal = db.Column(db.String(10), nullable=True)


class SlowQueryInfoViewSchema(ma.ModelSchema):
    class Meta:
        model = SlowQueryInfoView

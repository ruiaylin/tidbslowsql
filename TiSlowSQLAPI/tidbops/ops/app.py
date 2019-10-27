# coding: utf-8
import logging
from .models import db, ma
from .apis import api
from flask import Flask
from ops.configure import config
from flask_cors import CORS
from flask_cas import CAS


cas = CAS()
cors = CORS()


# create application
def create_app(config_name):
    app = Flask(
        __name__,
        static_folder=config[config_name].STATIC_PATH,
        template_folder=config[config_name].TEMPLATE_PATH
    )
    app.config.from_object(config[config_name])
    config[config_name].init_app(app)
    # login manager
    register_principals(app)
    register_api(app)
    register_database(app)
    cas.init_app(app)
    ma.init_app(app)
    CORS(app, resources={"/api/*": {"origins": "*"}})
    return app


def register_principals(app):
    return app


# register dblogical
def register_database(app):
    """Database related configuration."""
    # prepare for database
    db.init_app(app)
    db.app = app
    return app


# register dblogical
def register_api(app):
    """Database related configuration."""
    #: prepare for database
    api.init_app(app)
    api.app = app
    from .apis import slowsql
    api.add_resource(
        slowsql.SlowSQLAPI, '/api/v2/tislow')
    api.add_resource(
        slowsql.SlowSqlTrendAPI, '/api/v2/tislowcnt')
    api.add_resource(
        slowsql.SlowQueryTrendAPI, '/api/v2/tiquerycnt')
    api.add_resource(
        slowsql.SlowReportAPI, '/api/v2/tislowrpt')
    api.add_resource(
        slowsql.SlowInfoAPI, '/api/v2/tiqueryinfo')
    return app

# coding: utf-8
import os
import sys
from ops.app import create_app
from ops.models import db
from ops.models.slowsql import SlowQueryInfoView
from flask_script import Manager, Shell
from flask_migrate import Migrate, MigrateCommand

# set default encoding when python is version2.x
# reload(sys)
# sys.setdefaultencoding('utf8')
print (os.getenv('TiDBSlowSQL'))
# create app object
app = create_app(os.getenv('TiDBSlowSQL') or 'default')

app.debug = True

# manager
manager = Manager(app)
# migrate
migrate = Migrate(app, db)


def make_shell_context():
    return dict(
        app=app, db=db, SlowQueryInfoView=SlowQueryInfoView)


manager.add_command("shell", Shell(make_context=make_shell_context))
manager.add_command('db', MigrateCommand)


@manager.command
def createdb():
    """Create database for june."""
    from databus.models import db
    db.create_all()


@manager.command
def deploy():
    """Run deployment tasks."""
    from flask_migrate import upgrade

    # migrate database to latest revision
    upgrade()


if __name__ == '__main__':
    manager.run()

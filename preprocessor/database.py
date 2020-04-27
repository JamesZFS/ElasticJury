import MySQLdb


class MySQLWrapper:
    db_host = 'localhost'
    db_init_table_script = '../database/init-tables.sql'
    freq_commit = 50

    def __init__(self, drop=True):
        # Connection will not include password
        # Please manually add user/password or set your MySQL root settings to login without password
        self.connection = MySQLdb.connect(self.db_host)
        self.commands_not_committed = 0
        if drop:
            self.drop()
        self.create_and_switch()
        self.execute_script(self.db_init_table_script)
        print('SQL database initialized')

    def drop(self):
        self.execute('DROP DATABASE IF EXISTS ElasticJury')

    def create_and_switch(self):
        self.execute('CREATE DATABASE IF NOT EXISTS ElasticJury DEFAULT CHARACTER SET utf8')
        self.execute('USE ElasticJury')

    def execute(self, command, commit=False):
        cursor = self.connection.cursor()
        cursor.execute(command)
        cursor.close()
        if commit or (self.commands_not_committed % self.freq_commit == 0):
            self.commit()

    def commit(self):
        self.connection.commit()
        self.commands_not_committed = 0

    def execute_script(self, path):
        with open(path, 'r') as file:
            script = file.read()
            commands = script.split(';')
            commands = map(lambda x: x.strip(), commands)
            commands = filter(lambda x: len(x) > 0, commands)
            for command in commands:
                self.execute(command)


if __name__ == '__main__':
    db = MySQLWrapper()

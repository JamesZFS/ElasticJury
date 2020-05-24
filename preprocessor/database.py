# -*- coding: UTF-8 -*-

import MySQLdb

from entries import *
from utility import *


class MySQLWrapper:
    db_host = 'localhost'
    # db_host = 'cdb-f0b6x25m.cd.tencentcdb.com'
    db_port = 3306
    # db_port = 10104
    db_user = 'root'
    db_init_table_script = '../database/init-tables.sql'
    freq_commit = 50

    def __init__(self, drop=False, password=''):
        # Connection will not include password
        # Please manually add user/password or set your MySQL root settings to login without password
        log_info('Database', 'Connecting to MySQL database ...')
        self.connection = MySQLdb.connect(self.db_host, user=self.db_user, port=self.db_port,
                                          password=password, charset='utf8')
        self.commands_not_committed = 0
        if drop:
            self.drop()
        self.create_and_switch()
        self.execute_script(self.db_init_table_script)
        log_info('Database', 'MySQL database initialized')

    def drop(self):
        log_info('Database', 'Dropping original database ...')
        self.execute('DROP DATABASE IF EXISTS ElasticJury')

    def create_and_switch(self):
        self.execute('CREATE DATABASE IF NOT EXISTS ElasticJury DEFAULT CHARACTER SET utf8')
        self.execute('USE ElasticJury')

    def execute(self, command, values=None, commit=False):
        # log_info('Database', 'Executing command={} with values={}'.format(command, values))
        cursor = self.connection.cursor()
        if values:
            cursor.execute(command, values)
        else:
            cursor.execute(command)
        self.commands_not_committed += 1
        if commit or self.commands_not_committed > self.freq_commit == 0:
            self.commit()

    def execute_many(self, command, values, commit=False):
        cursor = self.connection.cursor()
        cursor.executemany(command, values)
        self.commands_not_committed += len(values)
        if commit or self.commands_not_committed > self.freq_commit == 0:
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

    def insert(self, entry: MySQLEntry):
        command, values = entry.generate_insert_command()
        self.execute(command, values)
        return self.connection.insert_id()

    def insert_many(self, entries: [MySQLEntry]):
        if len(entries) == 0:
            return

        command, _ = entries[0].generate_insert_command()
        values = []
        for entry in entries:
            entry_command, value = entry.generate_insert_command()
            assert entry_command == command
            values.append(value)
        self.execute_many(command, values)

    def close(self):
        self.commit()
        self.connection.close()
        log_info('Database', 'MySQL Database closed')


if __name__ == '__main__':
    db = MySQLWrapper()

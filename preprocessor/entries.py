# -*- coding: UTF-8 -*-

from utility import *


class MySQLEntry:
    def generate_insert_command(self):
        unimplemented(self)


class CaseEntry(MySQLEntry):
    separator = '#'

    def __init__(self, judges, laws, tags, keywords, detail, tree):
        self.judges = judges
        self.laws = laws
        self.tags = tags
        self.detail = detail
        self.tree = tree
        self.keywords = keywords

    def generate_insert_command(self):
        # Field id will automatically increase
        command = 'INSERT INTO Cases(judges, laws, tags, keywords, detail, tree) VALUES (%s,%s,%s,%s,%s,%s)'
        values = (self.separator.join(self.judges), self.separator.join(self.laws),
                  self.separator.join(self.tags), self.separator.join(self.keywords), self.detail, self.tree)
        return command, values


class IndexEntry(MySQLEntry):

    def __init__(self, table_name, key_name, key, case_id, weight):
        self.table_name = table_name
        self.key_name = key_name
        self.key = key
        self.case_id = case_id
        self.weight = weight

    def generate_insert_command(self):
        command = 'INSERT INTO {}({}, caseID, weight) VALUES (%s,%s,%s)' \
            .format(self.table_name, self.key_name)
        values = (self.key, self.case_id, self.weight)
        return command, values

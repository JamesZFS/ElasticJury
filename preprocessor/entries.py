from utility import *


class MySQLEntry:
    def generate_insert_command(self):
        unimplemented(self)


class CaseEntry(MySQLEntry):
    separator = '#'

    def __init__(self, title, judges, laws, tags, detail, tree):
        self.title = title
        self.judges = judges
        self.laws = laws
        self.tags = tags
        self.detail = detail
        self.tree = tree

    def generate_insert_command(self):
        # Field id will automatically increase
        command = 'INSERT INTO Cases(title, judge, law, tag, detail) VALUES (%s,%s,%s,%s,%s)'
        values = (self.title, self.separator.join(self.judges), self.separator.join(self.laws),
                  self.separator.join(self.tags), self.detail)
        return command, values


class WordEntry(MySQLEntry):

    def __init__(self, word, case_id, weight):
        self.word = word
        self.case_id = case_id
        self.weight = weight

    def generate_insert_command(self):
        command = 'INSERT INTO WordIndex(word, caseID, weight) VALUES (%s,%s,%s)'
        values = (self.word, self.case_id, self.weight)
        return command, values


class JudgeEntry(MySQLEntry):

    def __init__(self, judge, case_id, weight):
        self.judge = judge
        self.case_id = case_id
        self.weight = weight

    def generate_insert_command(self):
        command = 'INSERT INTO JudgeIndex(judge, caseID, weight) VALUES (%s,%s,%s)'
        values = (self.judge, self.case_id, self.weight)
        return command, values


class LawEntry(MySQLEntry):

    def __init__(self, law, case_id, weight):
        self.law = law
        self.case_id = case_id
        self.weight = weight

    def generate_insert_command(self):
        command = 'INSERT INTO LawIndex(law, caseID, weight) VALUES (%s,%s,%s)'
        values = (self.law, self.case_id, self.weight)
        return command, values


class TagEntry(MySQLEntry):

    def __init__(self, law, case_id, weight):
        self.tag = tag
        self.case_id = case_id
        self.weight = weight

    def generate_insert_command(self):
        command = 'INSERT INTO TagIndex(tag, caseID, weight) VALUES (%s,%s,%s)'
        values = (self.tag, self.case_id, self.weight)
        return command, values
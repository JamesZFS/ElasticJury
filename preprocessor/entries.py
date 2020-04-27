import utility


class MySQLEntry:
    def generate_insert_command(self):
        utility.unimplemented()


class Case(MySQLEntry):
    table = 'Cases'
    separator = '#'

    def __init__(self, title='', judges=[], laws=[], tags=[], link='', detail='', tree=''):
        self.title = title
        self.judges = judges
        self.laws = laws
        self.tags = tags
        self.link = link
        self.detail = detail
        self.tree = tree

    def generate_insert_command(self):
        return 'INSERT INTO Cases(id, title, judge, law, tag, link, detail) VALUES ({},{},{},{},{},{},{})' \
            .format(0, self.title, self.separator.join(self.judges),
                    self.separator.join(self.laws), self.separator.join(self.tags), self.link, self.detail)

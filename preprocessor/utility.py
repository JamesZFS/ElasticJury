import os


class UnimplementedError(Exception):
    def __init__(self):
        self.message = 'Unimplemented code'


def unimplemented():
    raise UnimplementedError()


def get_all_xml_files(path):
    xmls = []
    for home, dirs, files in os.walk(path):
        for file in files:
            if file.endswith('.xml'):
                xmls.append(os.path.join(home, file))
    return xmls

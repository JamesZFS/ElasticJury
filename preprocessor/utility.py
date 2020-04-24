import os


def get_all_xml_files(path):
    xmls = []
    for home, dirs, files in os.walk(path):
        for file in files:
            if file.endswith('.xml'):
                xmls.append(os.path.join(home, file))
    return xmls

# -*- coding: UTF-8 -*-

import argparse
import json

import idf_dict_collector
import mapping_collector
import processor

from utility import *

parser = argparse.ArgumentParser(description='ElasticJury Data Preprocessor')

# Configs
parser.add_argument('--path', type=str, default='../dataset/use', help='Relative path for data to process')
parser.add_argument('--mapping', type=str, default='mapping.json', help='Tag to nameCN mapping file path')
parser.add_argument('--idf_dict', type=str, default='idf_dict.json', help='Idf dict file path')
parser.add_argument('--password', type=str, default='', help='Password to login your local mysql account')
parser.add_argument('--clean-mapping', dest='clean_mapping', action='store_true', help='Clean the mapping file')
parser.add_argument('--clean-idf-dict', dest='clean_idf_dict', action='store_true', help='Clean the idf dict file')

# For debugging
parser.add_argument('--no-db', dest='do_db', action='store_false', help='Whether run database processing')
parser.add_argument('--no-mapping', dest='do_mapping', action='store_false', help='Whether run mapping')
parser.add_argument('--no-idf-dict', dest='do_idf_dict', action='store_false', help='Whether generate idf word dict')

parser.set_defaults(do_db=True, do_mapping=True, clean_mapping=False, do_idf_dict=True, clean_idf_dict=False)


def run(path, mapping_path, idf_dict_path, db_password, do_db, do_mapping, clean_mapping, do_idf_dict, clean_idf_dict):
    # Check if xml path exists
    if not os.path.exists(path):
        log_exit('Main', 'XML data directory path {} does not exist'.format(path))

    # Clean mapping
    if os.path.exists(mapping_path) and clean_mapping:
        os.remove(mapping_path)

    # Create mapping if not exists
    if (not os.path.exists(mapping_path)) and do_mapping:
        mapping = mapping_collector.get_mapping(path)
        mapping_collector.dump(mapping, mapping_path)

    # Clean idf dict
    if os.path.exists(idf_dict_path) and clean_idf_dict:
        os.remove(idf_dict_path)

    # Create idf dict if not exist
    if (not os.path.exists(idf_dict_path)) and do_idf_dict:
        idf_dict = idf_dict_collector.get_idf_dict(path)
        idf_dict_collector.dump(idf_dict, idf_dict_path)

    # Read mapping
    if not os.path.exists(mapping_path):
        log_exit('Main', 'Mapping file {} does not exist'.format(mapping_path))
    with open(mapping_path, 'r') as file:
        mapping = json.load(file)

    if not os.path.exists(idf_dict_path):
        log_exit('Main', 'Idf dict file {} does not exist'.format(idf_dict_path))
    with open(idf_dict_path, 'r') as file:
        idf_dict = json.load(file)

    # Processing and dump into database
    if do_db:
        processor.process(mapping, idf_dict, path, db_password)


if __name__ == '__main__':
    args = parser.parse_args()
    run(args.path, args.mapping, args.idf_dict, args.password, args.do_db,
        args.do_mapping, args.clean_mapping, args.do_idf_dict, args.clean_idf_dict)

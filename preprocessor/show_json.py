import argparse
import json

parser = argparse.ArgumentParser(description='Show json file into console')
parser.add_argument('--path', type=str, default='mapping.json', help='Json file path')
args = parser.parse_args()

if __name__ == '__main__':
    with open(args.path, 'r') as file:
        j = json.load(file)
        print(j)

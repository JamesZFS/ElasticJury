import json

if __name__ == '__main__':
    with open('mapping.json') as file:
        j = json.load(file)
        print(j)
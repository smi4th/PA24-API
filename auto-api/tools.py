import os, math, json

def SQL_to_dict() -> dict:

    tables = {}

    with open('createBDD.sql', "r") as bdd:
        sql = bdd.read()

        for query in sql.split(';')[0:-1]:
            if 'CREATE TABLE IF NOT EXISTS' in query:
            
                tableSeparator = query.split('CREATE TABLE IF NOT EXISTS ')[1].split(' ')[0].split('`')[1]
            
                tables[tableSeparator] = {}
                tables[tableSeparator]['columns'] = {}
                tables[tableSeparator]['foreignKeys'] = {}
                tables[tableSeparator]['primaryKeys'] = []
                tables[tableSeparator]['uniqueKeys'] = []
                tables[tableSeparator]['autoGen'] = []

                for line in query.split(f'CREATE TABLE IF NOT EXISTS `{tableSeparator}` (')[1].split('\n')[0:-1]:
                    if line != "":
                        
                        columnName = line.split('`')[1]
                        columnType = line.split('`')[2].split(' ')[1]
                            
                        tables[tableSeparator]['columns'][columnName] = columnType

                        if 'PRIMARY KEY' in line:
                            if 'PRIMARY KEY (' in line:
                                for key in line.split('PRIMARY KEY (')[1].split(')')[0].split(', '):
                                    key = key.replace('`', '')
                                    tables[tableSeparator]['primaryKeys'].append(key)
                                    tables[tableSeparator]['uniqueKeys'].append(key)
                            else:
                                tables[tableSeparator]['primaryKeys'].append(columnName)
                                tables[tableSeparator]['uniqueKeys'].append(columnName)
                        elif 'UNIQUE' in line:
                            tables[tableSeparator]['uniqueKeys'].append(columnName)
                        elif 'FOREIGN KEY' in line:
                            tables[tableSeparator]['foreignKeys'][columnName] = {}
                            tables[tableSeparator]['foreignKeys'][columnName]['references'] = line.split('REFERENCES `')[1].split('`')[0]
                            tables[tableSeparator]['foreignKeys'][columnName]['referencesColumn'] = line.split('REFERENCES `')[1].split('`')[2]
                        elif 'AUTO GEN' in line:
                            tables[tableSeparator]['autoGen'].append(columnName)

    return tables

def removeOldFiles() -> None:
    for path in ["requests", "pytest/jsonFiles"]:
        for (dirpath, dirnames, filenames) in os.walk(path):
            for filename in filenames:
                if 'go.mod' not in filename:
                    os.remove(f"{dirpath}/{filename}")
                
            if 'requests' not in path:
                import shutil; shutil.rmtree(dirpath)

def passwordCheck(file: str) -> str:
    file = file.replace('{{passwordCheck}}', f"if !tools.ValueIsEmpty(password_) {{\n\t\tif tools.PasswordNotStrong(password_) {{\n\t\t\ttools.JsonResponse(w, 400, `{{\"error\": \"Password is not strong enough\"}}`) \n\t\t\treturn\n\t\t}} else {{\n\t\t\tpassword_ = tools.HashPassword(password_)\n\t\t}}\n\t}}")
    return file

def emailCheck(file: str) -> str:
    file = file.replace('{{emailCheck}}', f"if !tools.ValueIsEmpty(email_) {{\n\t\tif !tools.EmailIsValid(email_) {{\n\t\t\ttools.JsonResponse(w, 400, `{{\"error\": \"Email is not valid\"}}`) \n\t\t\treturn\n\t\t}}\n\t}}")
    return file

def foreignKeyCheck(file: str, tables: dict, table: str) -> str:
    for foreignKey in tables[table]['foreignKeys']:
        file = file.replace('{{foreignKeyTests}}', f"if !tools.ValueIsEmpty({foreignKey}_) {{\n\t\tif !tools.ElementExists(db, \"{tables[table]['foreignKeys'][foreignKey]['references']}\", \"{tables[table]['foreignKeys'][foreignKey]['referencesColumn']}\", {foreignKey}_) {{\n\t\t\ttools.JsonResponse(w, 400, `{{\"error\": \"This {foreignKey} does not exist\"}}`) \n\t\t\treturn\n\t\t}}\n\t}}\n\t{{{{foreignKeyTests}}}}")
    return file

def uniqueCheck(file: str, table: str, unique: str) -> str:
    file = file.replace('{{uniqueFieldsTests}}', f"if tools.ElementExists(db, \"{table}\", \"{unique}\", {unique}_) {{\n\t\ttools.JsonResponse(w, 400, `{{\"error\": \"This {unique} already exists\"}}`) \n\t\treturn\n\t}}\n\t{{{{uniqueFieldsTests}}}}")

    return file

class Json:
    def __init__(self : 'Json', tables: dict, table: str, method: str, url: str, urlParams: dict, requestBody: dict, statusCode: int, responseBody: dict) -> None:
        self.json = {
            "request": {
                "url": "",
                "urlParams" : {},
                "method": "",
                "body": {}
            },
            "response": {
                "status_code": 0,
                "body": {}
            }
        }
        self.name = self.getJsonName(tables, table, method)
        self.json["request"]["url"] = url.lower()
        self.json["request"]["urlParams"] = urlParams
        self.json["request"]["method"] = method
        self.setRequestBody(requestBody)
        self.json["response"]["status_code"] = statusCode
        self.json["response"]["body"] = responseBody

        self.createJson(tables, table, method)

        self.save()

    def setRequestBody(self: 'Json', requestBody: dict) -> None:

        asoc = {
            "RANDOMIZED": ["CHAR", "VARCHAR", "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT", "BINARY", "VARBINARY", "TINYBLOB", "MEDIUMBLOB", "LONGBLOB", "BLOB", "ENUM", "SET"],
        }

        # for key, value in requestBody.items():
        #     print(list((k, v) for k, v in asoc.items()))
        #     print(list((k2 for k2, v2 in list(((k, v) for k, v in asoc.items()))))[0])
            #self.json["request"]["body"][key] = (k2 for k2, v2 in list(((k, v) for k, v in asoc.items())) if v2 in value)

    def save(self: 'Json') -> None:
        with open(self.name, 'w') as file:
            file.write(json.dumps(self.json, indent=4))

    def getJsonName(self: 'Json', tables: dict, table: str, method: str) -> str:
        # The files/folders need to be organized like so:
        # pytest/jsonFiles/{n table}_{table}/{n method}{method}/test_{method}_{n test}.json
        # where:
        # - n table is the number of the table in the list of tables
        # - table is the name of the table in lowercase and without spaces or special characters
        # - n method is the number of the method in the list of methods (create, read, update, delete)
        # - method is the name of the method in lowercase and without spaces or special characters
        # - n test is the number of the test in the list of tests

        # get the number of existing tables
        n_table = str(list(tables.keys()).index(table)).zfill(math.ceil(math.log10(len(tables))))

        # get the number of existing methods
        method = {"POST" : "create", "GET" : "read", "PUT" : "update", "DELETE" : "delete"}[method]
        
        n_method = list({"create" : "001", "read" : "002", "update" : "003", "delete" : "004"}.keys()).index(method) + 1

        # get the number of existing tests
        n_test = "001"
        for (dirpath, dirnames, filenames) in os.walk(f"pytest/jsonFiles/{n_table}_{''.join(table)}"):
            n_test = str(len(filenames)).zfill(3)

        return f"pytest/jsonFiles/{n_table}_{''.join(table.lower())}/{n_method}{method}/test_{method}_{n_test}.json"

    def createJson(self: 'Json', tables: dict, table: str, method: str) -> None:
        jsonName = self.getJsonName(tables, table, method)

        # the path may not exist yet
        if not os.path.exists(os.path.dirname(jsonName)):
            os.makedirs(os.path.dirname(jsonName))
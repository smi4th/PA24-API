import tools

def replacePost(file: str, tables: dict, table: str) -> str:
    if len(tables[table]['primaryKeys']) == 1:
        file = file.replace('{{fields}}', f"`{'`, `'.join((field for field in tables[table]['columns'].keys() if field not in (tables[table]['primaryKeys'] + tables[table]['autoGen'])))}`")
        file = file.replace('{{fieldsVar}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['columns'].keys() if fieldVar not in (tables[table]['primaryKeys'] + tables[table]['autoGen'] + list(tables[table]['foreignKeys'].keys()))))}")
        file = file.replace('{{fieldsVarSQL}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['columns'].keys() if fieldVar not in (tables[table]['primaryKeys'] + tables[table]['autoGen'])))}")
    else: 
        file = file.replace('{{fields}}', f"`{'`, `'.join((field for field in tables[table]['columns'].keys() if field not in tables[table]['autoGen']))}`")
        file = file.replace('{{fieldsVar}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['columns'].keys() if fieldVar not in (tables[table]['autoGen'] + list(tables[table]['foreignKeys'].keys()))))}")
        file = file.replace('{{fieldsVarSQL}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['columns'].keys() if fieldVar not in tables[table]['autoGen']))}")
    file = file.replace('{{primaryKeys}}', f"`{'`, `'.join(tables[table]['primaryKeys'])}`")
    file = file.replace('{{primaryKeysVar}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['primaryKeys'] if fieldVar not in tables[table]['autoGen']))}")

    for var in tables[table]['columns']:
        if var not in (tables[table]['primaryKeys'] + tables[table]["autoGen"]) or (len(tables[table]['primaryKeys']) != 1 and var not in tables[table]['autoGen']):
            file = file.replace('{{fields to var from body}}', f"{var}_ := tools.BodyValueToString(body, \"{var}\")\n\t" + '{{fields to var from body}}')
    file = file.replace('{{fields to var from body}}', "")

    for unique in tables[table]['uniqueKeys']:
        if unique not in (tables[table]['primaryKeys'] + tables[table]['autoGen']):
            file = tools.uniqueCheck(file, table, unique)
            json = tools.Json(tables, table, 'POST', table, {}, tables[table]['columns'], 400, {"error": f"This {unique} already exists"})

    file = file.replace('{{uniqueFieldsTests}}', "")

    if len(tables[table]['primaryKeys']) == 1:
        file = file.replace('{{primaryKeyUUID}}', f"{tables[table]['primaryKeys'][0]}_ := tools.GenerateUUID()")
    file = file.replace('{{primaryKeyUUID}}', "")

    file = file.replace('{{questionMarks}}', f"{', '.join(('?' for _ in tables[table]['columns'] if _ not in tables[table]['autoGen']))}")

    if 'password' in tables[table]['columns']:
        file = tools.passwordCheck(file)
    if 'email' in tables[table]['columns']:
        file = tools.emailCheck(file)
    if len(list(tables[table]['foreignKeys'].keys())) > 0:
        file = tools.foreignKeyCheck(file, tables, table)
        
    file = file.replace('{{foreignKeyTests}}', "")
    file = file.replace('{{emailCheck}}', "")
    file = file.replace('{{passwordCheck}}', "")

    return file
    
def replaceGet(file: str, tables: dict, table: str) -> str:
    file = file.replace('{{fields}}', f"`{'`, `'.join((field for field in tables[table]['columns'].keys() if field != 'password'))}`")

    for var in tables[table]['columns']:
        if var not in tables[table]['primaryKeys']:
            file = file.replace('{{fields to var from body}}', f"{var}_ := tools.BodyValueToString(body, \"{var}\")\n\t" + '{{fields to var from body}}')
    file = file.replace('{{fields to var from body}}', "")

    for unique in tables[table]['uniqueKeys']:
        file = tools.uniqueCheck(file, table, unique)
    file = file.replace('{{uniqueFieldsTests}}', "")

    file = file.replace('{{questionMarks}}', f"{', '.join(('?' for _ in tables[table]['columns']))}")

    return file

def replacePut(file: str, tables: dict, table: str) -> str:
    file = file.replace('{{fields}}', f"`{'`, `'.join((field for field in tables[table]['columns'].keys() if field not in (tables[table]['primaryKeys'] + tables[table]['autoGen'])))}`")
    file = file.replace('{{fieldsVar}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['columns'].keys() if fieldVar not in (tables[table]['primaryKeys'] + tables[table]['autoGen'] + list(tables[table]['foreignKeys'].keys()))))}")
    file = file.replace('{{primaryKeys}}', f"`{'`, `'.join(tables[table]['primaryKeys'])}`")
    file = file.replace('{{primaryKeysVar}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['primaryKeys'] if fieldVar not in tables[table]['autoGen']))}")

    for var in tables[table]['columns']:
        if var not in (tables[table]['primaryKeys'] + tables[table]['autoGen']):
            file = file.replace('{{fields to var from body}}', f"{var}_ := tools.BodyValueToString(body, \"{var}\")\n\t" + '{{fields to var from body}}')
        elif var not in tables[table]['autoGen']:
            file = file.replace('{{primaryKeys to var from query}}', f"{var}_ := query[\"{var}\"]\n\t" + '{{primaryKeys to var from query}}')
    file = file.replace('{{fields to var from body}}', "")
    file = file.replace('{{primaryKeys to var from query}}', "")

    for unique in tables[table]['uniqueKeys']:
        if unique not in (tables[table]['primaryKeys'] + tables[table]['autoGen']):
            file = tools.uniqueCheck(file, table, unique)
        elif unique not in tables[table]['autoGen']:
            file = file.replace('{{uniqueFieldsTests}}', f"if !tools.ElementExists(db, \"{table}\", \"{unique}\", {unique}_) {{\n\t\ttools.JsonResponse(w, 400, `{{\"error\": \"This " + "{{routeName}}" + f" does not exist\"}}`) \n\t\treturn\n\t}}\n\t{{{{uniqueFieldsTests}}}}")
    file = file.replace('{{uniqueFieldsTests}}', "")

    if 'password' in tables[table]['columns']:
        file = tools.passwordCheck(file)
    if 'email' in tables[table]['columns']:
        file = tools.emailCheck(file)
    if len(list(tables[table]['foreignKeys'].keys())) > 0:
        file = tools.foreignKeyCheck(file, tables, table)
        
    file = file.replace('{{foreignKeyTests}}', "")
    file = file.replace('{{emailCheck}}', "")
    file = file.replace('{{passwordCheck}}', "")

    file = file.replace('{{primaryKeysSQL}}', f"{', '.join((f'{key} = ?' for key in tables[table]['primaryKeys']))}")

    file = file.replace('{{questionMarks}}', f"{', '.join(('?' for _ in tables[table]['columns']))}")

    return file

def replaceDelete(file: str, tables: dict, table: str) -> str:
    file = file.replace('{{primaryKeys}}', f"`{'`, `'.join(tables[table]['primaryKeys'])}`")
    file = file.replace('{{primaryKeysVar}}', f"{', '.join((fieldVar + '_' for fieldVar in tables[table]['primaryKeys']))}")

    for var in tables[table]['primaryKeys']:
        file = file.replace('{{primaryKeys to var from query}}', f"{var}_ := query[\"{var}\"]\n\t" + '{{primaryKeys to var from query}}')
        file = file.replace('{{uniqueFieldsTests}}', f"if !tools.ElementExists(db, \"{table}\", \"{var}\", {var}_) {{\n\t\ttools.JsonResponse(w, 400, `{{\"error\": \"This " + "{{routeName}}" + f" does not exist\"}}`) \n\t\treturn\n\t}}\n\t{{{{uniqueFieldsTests}}}}")
    file = file.replace('{{uniqueFieldsTests}}', "")
    file = file.replace('{{primaryKeys to var from query}}', "")

    for foreignKey in tables[table]['foreignKeys']:
        file = file.replace('{{foreignKeyTests}}', f"if !tools.ElementExists(db, \"{tables[table]['foreignKeys'][foreignKey]['references']}\", \"{tables[table]['foreignKeys'][foreignKey]['referencesColumn']}\", {foreignKey}_) {{\n\t\ttools.JsonResponse(w, 400, `{{\"error\": \"This {foreignKey} does not exist\"}}`) \n\t\treturn\n\t}}\n\t{{{{foreignKeyTests}}}}")
    file = file.replace('{{foreignKeyTests}}', "")

    file = file.replace('{{primaryKeysSQL}}', f"{', '.join((f'{key} = ?' for key in tables[table]['primaryKeys']))}")

    file = file.replace('{{primaryKeysJSON}}', ', '.join((f'"{key}": "` + {key}_ + `"' for key in tables[table]['primaryKeys'])))

    return file

def replaceGetAll(file: str, tables: dict, table: str) -> str:
    file = file.replace('{{fields}}', f"`{'`, `'.join((field for field in tables[table]['columns'].keys() if field != 'password'))}`")
    file = file.replace('{{fieldsVar}}', ', '.join((field + '_' for field in tables[table]['columns'].keys() if field != 'password')))
    file = file.replace('{{fieldsVarPointer}}', ', '.join((f"&{field}_" for field in tables[table]['columns'].keys() if field != 'password')))
    file = file.replace('{{primaryKeys}}', ', '.join(f"{field}_" for field in tables[table]['primaryKeys']))
    file = file.replace('{{primaryKeysSQL}}', f"{', '.join((f'{key} = ?' for key in tables[table]['primaryKeys']))}")
    file = file.replace('{{primaryKeysVarFunc}}', ', '.join(f"{field}_ string" for field in tables[table]['primaryKeys']))

    for var in tables[table]['columns']:
        if var not in (tables[table]['primaryKeys'] + tables[table]['autoGen']):
            file = file.replace('{{fields to var from body}}', f"{var}_ := tools.BodyValueToString(body, \"{var}\")\n\t" + '{{fields to var from body}}')
    file = file.replace('{{fields to var from body}}', "")
    
    file = file.replace('{{fieldsJson}}', ', '.join((f'"{field}": "` + {field}_ + `"' for field in tables[table]['columns'].keys() if field != 'password')))
    file = file.replace('{{fieldsJson}}', "")

    return file

def main() -> None:

    tables = tools.SQL_to_dict()

    with open("auto-api/base.txt", "r") as base:
        base = base.read()

    with open("auto-api/post.txt", "r") as post:
        post = post.read()

    with open("auto-api/get.txt", "r") as get:
        get = get.read()

    with open("auto-api/put.txt", "r") as put:
        put = put.read()

    with open("auto-api/delete.txt", "r") as delete:
        delete = delete.read()

    with open("auto-api/getAll.txt", "r") as getAll:
        getAll = getAll.read()

    for table in tables:

        file = ""

        postResult = replacePost(post, tables, table)
        getResult = replaceGet(get, tables, table)
        putResult = replacePut(put, tables, table)
        deleteResult = replaceDelete(delete, tables, table)
        getAllResult = replaceGetAll(getAll, tables, table)

        file = base + "\n\n" + postResult + "\n\n" + getResult + "\n\n" + putResult + "\n\n" + deleteResult + "\n\n" + getAllResult
        file = file.replace('{{routeName}}', ''.join((letter.capitalize() for letter in table.split('_'))))
        file = file.replace('{{routeTable}}', table)

        with open(f"requests/{table.lower()}.go", "w") as go:
            go.write(file)


if __name__ == "__main__":
    tools.removeOldFiles()
    main()
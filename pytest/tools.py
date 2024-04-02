import json, random, string, ast

def formatJson(jsonPath):

    with open(jsonPath, "r") as f:
        data = json.load(f)
        for r in ["request", "response"]:
            # if it is not a list
            if len(str(data[r]["body"]).split("]")) == 1:
                for key in data[r]["body"]:
                    for key, value in data[r]["body"].items():
                        if value == "RANDOMIZED":
                            data[r]["body"][key] = ''.join(random.choices(string.ascii_uppercase + string.digits, k = 10)) + "@gmail.com"
                            if data["response"]["status_code"] in [200, 201]:
                                writeJson(jsonPath, data[r]["body"])
                        elif value == "INJECT":
                            with open("pytest/temp.json", "r") as F:
                                data[r]["body"][key] = json.load(F)[''.join(jsonPath.split("/")[-3].split("_")) + "_" + key]
            else:
                data["response"]["body"] = ast.literal_eval(data["response"]["body"])
                for element in data["response"]["body"]:
                    for key, value in element.items():
                        if value == "RANDOMIZED":
                            element[key] = ''.join(random.choices(string.ascii_uppercase + string.digits, k = 10)) + "@gmail.com"
                            if data["response"]["status_code"] in [200, 201]:
                                writeJson(jsonPath, element)
                        elif value == "INJECT":
                            with open("pytest/temp.json", "r") as F:
                                element[key] = json.load(F)[''.join(jsonPath.split("/")[-3].split("_")) + "_" + key]

    return data

def writeJson(jsonPath, data):
    try:
        with open("pytest/temp.json", "x") as F:
            json.dump({"empty": "empty"}, F)
    except:
        pass
    with open("pytest/temp.json", "r+") as f:
        temp = json.load(f)
        for key, value in data.items():
            temp[''.join(jsonPath.split("/")[-3].split("_")) + "_" + key] = value
        f.seek(0)
        json.dump(temp, f)
        f.truncate()

def testValues(data, response, jsonPath):
    if len(str(data).split("]")) == 1:
        for key, value in data.items():
            if value not in ["UNPREDEFINED", "INJECT"]:
                assert response.json()[key] == value
            elif value == "INJECT":
                with open("pytest/temp.json", "r") as f:
                    assert response.json()[key] == json.load(f)[''.join(jsonPath.split("/")[-3].split("_")) + "_" + key]
    else:
        newResponse = ast.literal_eval(response.text)
        for element in data:
            for key, value in element.items():
                if value not in ["UNPREDEFINED", "INJECT"]:
                    assert newResponse[data.index(element)][key] == value
                elif value == "INJECT":
                    with open("pytest/temp.json", "r") as f:
                        assert newResponse[data.index(element)][key] == json.load(f)[''.join(jsonPath.split("/")[-3].split("_")) + "_" + key]
    
    if response.status_code in [200, 201]:
        if len(str(data).split("]")) == 1:
            writeJson(jsonPath, response.json())
        else:
            for element in data:
                writeJson(jsonPath, newResponse[data.index(element)])
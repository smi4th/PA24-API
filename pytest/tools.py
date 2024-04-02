import json, random, string, ast, requests

def callAPI(data, headers, url):
    print(data, headers, url)
    match data["request"]["method"]:
        case "GET":
            response = requests.get(url, headers=headers, json=data["request"]["body"])
        case "POST":
            response = requests.post(url, headers=headers, json=data["request"]["body"])
        case "PUT":
            response = requests.put(url, headers=headers, json=data["request"]["body"])
        case "DELETE":
            response = requests.delete(url, headers=headers, json=data["request"]["body"])
        case _:
            raise ValueError("Invalid method")
        
    print(response.text)
        
    return response

def formatJson(jsonPath):

    with open(jsonPath, "r") as f:
        data = json.load(f)
        for r in ["request", "response"]:
            # if it is not a list
            if len(str(data[r]["body"]).split("]")) == 1:

                # For the urlParams
                for u in data["request"]["urlParams"]:
                    if data["request"]["urlParams"][u] == "RANDOMIZED":
                        data["request"]["urlParams"][u] = ''.join(random.choices(string.ascii_uppercase + string.digits, k = 10)) + "@gmail.com"
                    elif data["request"]["urlParams"][u] == "INJECT":
                        with open("pytest/temp.json", "r") as F:
                            data["request"]["urlParams"][u] = json.load(F)[jsonPath.split("/")[-3].split("_")[1] + "_" + u]
                    elif "INJECT_FOREIGN" in data["request"]["urlParams"][u]:
                        with open("pytest/temp.json", "r") as F:
                            data["request"]["urlParams"][u] = json.load(F)[data["request"]["urlParams"][u].split("INJECT_FOREIGN:")[1]]
                
                # For the body
                for key, value in data[r]["body"].items():
                    if value == "RANDOMIZED":
                        data[r]["body"][key] = ''.join(random.choices(string.ascii_uppercase + string.digits, k = 10)) + "@gmail.com"
                        if data["response"]["status_code"] in [200, 201]:
                            writeJson(jsonPath, data[r]["body"])
                    elif value == "INJECT":
                        with open("pytest/temp.json", "r") as F:
                            data[r]["body"][key] = json.load(F)[jsonPath.split("/")[-3].split("_")[1] + "_" + key]
                    elif "INJECT_FOREIGN" in data[r]["body"][key]:
                        with open("pytest/temp.json", "r") as F:
                            data[r]["body"][key] = json.load(F)[data[r]["body"][key].split("INJECT_FOREIGN:")[1]]
            else:
                data[r]["body"] = ast.literal_eval(data[r]["body"])
                for element in data[r]["body"]:
                    for key, value in element.items():
                        if value == "RANDOMIZED":
                            element[key] = ''.join(random.choices(string.ascii_uppercase + string.digits, k = 10)) + "@gmail.com"
                            if data[r]["status_code"] in [200, 201]:
                                writeJson(jsonPath, element)
                        elif value == "INJECT":
                            with open("pytest/temp.json", "r") as F:
                                element[key] = json.load(F)[jsonPath.split("/")[-3].split("_")[1] + "_" + key]

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
            temp[jsonPath.split("/")[-3].split("_")[1] + "_" + key] = value
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
                    assert response.json()[key] == json.load(f)[jsonPath.split("/")[-3].split("_")[1] + "_" + key]
    else:
        newResponse = ast.literal_eval(response.text)
        for element in data:
            for key, value in element.items():
                if value not in ["UNPREDEFINED", "INJECT"]:
                    assert newResponse[data.index(element)][key] == value
                elif value == "INJECT":
                    with open("pytest/temp.json", "r") as f:
                        assert newResponse[data.index(element)][key] == json.load(f)[jsonPath.split("/")[-3].split("_")[1] + "_" + key]
    
    if response.status_code in [200, 201]:
        if len(str(data).split("]")) == 1:
            writeJson(jsonPath, response.json())
        else:
            for element in data:
                writeJson(jsonPath, newResponse[data.index(element)])
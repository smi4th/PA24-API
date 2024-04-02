import requests, json, tools

def create_account(
    api_url: str,
    working_dir: str,

    body: dict
):
    url = f"{api_url}/api/account"
    return requests.post(url, json=body)

def test_create_account_1(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_1.json")
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_1.json")

def test_create_account_2(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_2.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_2.json")

def test_create_account_3(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_3.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_3.json")

def test_create_account_4(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_4.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_4.json")

def test_create_account_5(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_5.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_5.json")

def test_create_account_6(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_6.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_6.json")

def test_create_account_7(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_7.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_7.json")

def test_create_account_8(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_8.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_8.json")

def test_create_account_9(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_9.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_9.json")

def test_create_account_10(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/create/test_create_10.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["account_type"] = json.load(f)["accounttype_id"]
    response = create_account(api_url, working_dir, data["request"]["body"])

    print(data, response.json())
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account/create/test_create_10.json")
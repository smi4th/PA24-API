import requests, json, pytest, tools

def update_accountType(
    api_url: str,
    working_dir: str,

    body: dict,
    accountTypeID: str = ""
):
    url = f"{api_url}/api/account_type?id={accountTypeID}"
    return requests.put(url, json=body)

def test_update_accountType_1(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_1.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        accountTypeID = json.load(f)["accounttype_id"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_1.json")

def test_update_accountType_2(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_2.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        accountTypeID = json.load(f)["accounttype_id"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_2.json")

def test_update_accountType_3(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_3.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        accountTypeID = json.load(f)["accounttype_id"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_3.json")

def test_update_accountType_4(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_4.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        accountTypeID = json.load(f)["accounttype_id"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)

    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_4.json")

# @pytest.mark.skip(reason="No ID passed test. Not implemented yet.")
def test_update_accountType_5(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_5.json")
    accountTypeID = data["request"]["body"]["type"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_5.json")

def test_update_accountType_6(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_6.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        accountTypeID = json.load(f)["accounttype_id"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_6.json")

def test_update_accountType_7(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/update/test_update_7.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        accountTypeID = json.load(f)["accounttype_id"]
    response = update_accountType(api_url, working_dir, data["request"]["body"], accountTypeID)

    print(data, response.json())
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/update/test_update_7.json")
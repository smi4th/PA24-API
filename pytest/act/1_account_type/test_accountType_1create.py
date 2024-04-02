import requests, json, tools

def create_accountType(
    api_url: str,

    body: dict
):
    url = f"{api_url}/api/account_type"
    return requests.post(url, json=body)

def test_create_accountType_1(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/create/test_create_1.json")
    response = create_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/create/test_create_1.json")

def test_create_accountType_2(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/create/test_create_2.json")
    response = create_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/create/test_create_2.json")

def test_create_accountType_3(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/create/test_create_3.json")
    response = create_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/create/test_create_3.json")

def test_create_accountType_4(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/create/test_create_4.json")
    response = create_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/create/test_create_4.json")

def test_create_accountType_5(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/create/test_create_5.json")
    response = create_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/create/test_create_5.json")


def test_create_accountType_6(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/create/test_create_6.json")
    response = create_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/create/test_create_6.json")
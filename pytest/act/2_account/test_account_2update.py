import requests, json, tools

def update_account(
    api_url: str,
    working_dir: str,

    body: dict
):
    with open(working_dir + "/pytest/temp.json", "r") as f:
        url = f"{api_url}/api/account?id={json.load(f)['account_id']}"
    return requests.put(url, json=body)

def test_update_account_1(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/update/test_update_1.json")
    response = update_account(api_url, working_dir, data["request"]["body"])
    assert response.status_code == data["response"]["status_code"]
    for key, value in data["response"]["body"].items():
        if value != "UNPREDEFINED":
            assert response.json()[key] == value

def test_update_account_2(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/update/test_update_2.json")
    response = update_account(api_url, working_dir, data["request"]["body"])
    assert response.status_code == data["response"]["status_code"]
    for key, value in data["response"]["body"].items():
        if value != "UNPREDEFINED":
            assert response.json()[key] == value

def test_update_account_3(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/update/test_update_3.json")
    response = update_account(api_url, working_dir, data["request"]["body"])
    assert response.status_code == data["response"]["status_code"]
    for key, value in data["response"]["body"].items():
        if value != "UNPREDEFINED":
            assert response.json()[key] == value

def test_update_account_4(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/update/test_update_4.json")
    response = update_account(api_url, working_dir, data["request"]["body"])
    assert response.status_code == data["response"]["status_code"]
    for key, value in data["response"]["body"].items():
        if value != "UNPREDEFINED":
            assert response.json()[key] == value

def test_update_account_5(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account/update/test_update_5.json")
    response = update_account(api_url, working_dir, data["request"]["body"])
    assert response.status_code == data["response"]["status_code"]
    for key, value in data["response"]["body"].items():
        if value != "UNPREDEFINED":
            assert response.json()[key] == value
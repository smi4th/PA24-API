import requests, json, tools

def delete_account(
    api_url: str,
    urlParams: dict
):
    url = f"{api_url}/api/account_type?" + "&".join([f"{key}={value}" for key, value in urlParams.items()])
    return requests.delete(url)

def test_delete_accountType_1(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/delete/test_delete_1.json")
    response = delete_account(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/delete/test_delete_1.json")

def test_delete_accountType_2(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/delete/test_delete_2.json")
    response = delete_account(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/delete/test_delete_2.json")

def test_delete_accountType_3(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/delete/test_delete_3.json")
    with open(working_dir + "/pytest/temp.json", "r") as f:
        data["request"]["body"]["id"] = json.load(f)["accounttype_id"]
    response = delete_account(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/delete/test_delete_3.json")
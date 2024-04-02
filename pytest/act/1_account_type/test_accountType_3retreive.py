import requests, json, tools

def retreive_accountType(
    api_url: str,
    urlParams: dict
):
    url = f"{api_url}/api/account_type?" + "&".join([f"{key}={value}" for key, value in urlParams.items()])
    return requests.get(url)

def test_retreive_accountType_1(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/retreive/test_retreive_1.json")
    response = retreive_accountType(api_url, data["request"]["body"])
    
    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/retreive/test_retreive_1.json")

def test_retreive_accountType_2(api_url, working_dir):
    data = tools.formatJson(working_dir + "/pytest/jsonFiles/account_type/retreive/test_retreive_2.json")
    response = retreive_accountType(api_url, data["request"]["body"])

    assert response.status_code == data["response"]["status_code"]

    tools.testValues(data["response"]["body"], response, working_dir + "/pytest/jsonFiles/account_type/retreive/test_retreive_2.json")
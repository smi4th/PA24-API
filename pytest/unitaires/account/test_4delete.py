import requests, json

def delete_account(
    api_url: str,
    urlParams: dict
):
    url = f"{api_url}/api/account?" + "&".join([f"{key}={value}" for key, value in urlParams.items()])
    return requests.delete(url)

def test_delete_account_1(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/delete/test_delete_1.json") as f:
        data = json.load(f)
        response = delete_account(api_url, data["request"]["urlParams"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_delete_account_2(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/delete/test_delete_2.json") as f:
        data = json.load(f)
        response = delete_account(api_url, data["request"]["urlParams"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_delete_account_3(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/delete/test_delete_3.json") as f:
        data = json.load(f)
        with open(working_dir + "/pytest/temp.txt", "r") as f:
            data["request"]["urlParams"]["id"] = f.read()
        response = delete_account(api_url, data["request"]["urlParams"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value
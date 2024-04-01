import requests, json

def retreive_account(
    api_url: str,
    urlParams: dict
):
    url = f"{api_url}/api/account?" + "&".join([f"{key}={value}" for key, value in urlParams.items()])
    return requests.get(url)

def test_retreive_account_1(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/retreive/test_retreive_1.json") as f:
        data = json.load(f)
        response = retreive_account(api_url, data["request"]["urlParams"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_retreive_account_2(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/retreive/test_retreive_2.json") as f:
        data = json.load(f)
        response = retreive_account(api_url, data["request"]["urlParams"])
        assert response.status_code == data["response"]["status_code"]
        data["response"]["body"] = data["response"]["body"][0:1].split("}")[0:-1]
        for element in data["response"]["body"]:
            element += "}"
            for key, value in json.loads(element).items():
                if value != "UNPREDEFINED":
                    assert response.json()[key] == value
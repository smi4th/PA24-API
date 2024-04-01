import requests, json

def create_account(
    api_url: str,

    body: dict
):
    url = f"{api_url}/api/account"
    return requests.post(url, json=body)

def test_create_account_1(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_1.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_2(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_2.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_3(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_3.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_4(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_4.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_5(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_5.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_6(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_6.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_7(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_7.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_8(api_url, working_dir): # This test is expected to work
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_8.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value
            if key == "id":
                with open(working_dir + "/pytest/temp.txt", "w") as f:
                    f.write(str(response.json()["id"]))

def test_create_account_9(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_9.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value

def test_create_account_10(api_url, working_dir):
    with open(working_dir + "/pytest/jsonFiles/account/create/test_create_10.json") as f:
        data = json.load(f)
        response = create_account(api_url, data["request"]["body"])
        assert response.status_code == data["response"]["status_code"]
        for key, value in data["response"]["body"].items():
            if value != "UNPREDEFINED":
                assert response.json()[key] == value
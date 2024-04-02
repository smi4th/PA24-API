import pytest

@pytest.fixture(scope="session")
def api_url():
    return "http://localhost/api/"

@pytest.fixture(scope="session")
def working_dir():
    import os; return os.getcwd()

@pytest.fixture(scope="session")
def headers():
    return {"Content-Type": "application/json; charset=utf-8"}
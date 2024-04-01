import pytest

@pytest.fixture(scope="session")
def api_url():
    return "http://localhost"

@pytest.fixture(scope="session")
def working_dir():
    import os; return os.getcwd()
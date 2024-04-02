def test_removeTemp():
    import os; os.remove("pytest/temp.json")
    
    assert True
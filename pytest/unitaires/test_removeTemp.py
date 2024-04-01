def test_removeTemp():
    import os; os.remove("pytest/temp.txt")
    
    assert True
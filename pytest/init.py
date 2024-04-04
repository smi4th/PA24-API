import os

files = dict()

for (dirpath, dirnames, filenames) in os.walk("pytest"):
    for filename in filenames:
        if '.json' in filename:
            files[dirpath.replace("\\", "/").replace("pytest/jsonFiles/", "")] = files.get(dirpath.replace("\\", "/").replace("pytest/jsonFiles/", ""), []) + [filename]

ROOT = "pytest/pyTests"

# create the new files
for key, value in files.items():

    path = ''

    if "delete" in key:
        path = f"{ROOT.split('/')[0]}/pyTestsCleanup/test_{list(files).index(list(files)[-int(key.split('/')[0].split('_')[0])])}_{key.split('/')[0].split('_')[1]}_{value[0].split('_')[1]}.py"
    else:
        path = f"{ROOT}/{key.split('/')[0]}/test_{key.split('/')[0].split('_')[1]}_{key.split('/')[1]}.py"

    os.makedirs(os.path.dirname(path), exist_ok=True)
    
    with open(path, "w") as f:
        f.write("import tools\n\n")
        with open("pytest/model.txt", "r") as m:
            model = m.read()
            value = sorted(value)

            for v in value:
                modelModified = model
                modelModified = modelModified.replace("JSONFILE", f"test_{key.split('/')[0].split('_')[1]}_{'_'.join(v.split('_')[1::]).replace('.json', '')}")
                modelModified = modelModified.replace("JSONPATH", f'"/pytest/jsonFiles/{key}/{v}"')
                f.write(modelModified)
                f.write("\n\n")
                
with open(f"{ROOT.split('/')[0]}/pyTestsCleanup/test_removeTemp.py", "w") as f:
    f.write("def test_removeTemp(): import os; os.remove('pytest/temp.json'); assert True")
    # f.write("assert True")
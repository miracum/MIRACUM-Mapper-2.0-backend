import http.client
import random

conn = http.client.HTTPConnection("localhost:8080")

headers = {
    'accept': "application/json",
    'Content-Type': "application/json",
    'Authorization': "Bearer eyJhbGciOiJFUzI1NiIsImtpZCI6ImZha2Uta2V5LWlkIiwidHlwIjoiSldUIn0.eyJhdWQiOlsiZXhhbXBsZS11c2VycyJdLCJpc3MiOiJmYWtlLWlzc3VlciIsInBlcm0iOlsibm9ybWFsIl19.icXYWDt_Vf-eoQXlqWiqLzhuurJQj7sdfk1DxkLhnMSbClsLvz5T4AT7yCOhrk3BsczcYK6OukZ-vfD5-PChWA"
}

for _ in range(10000):
    payload = f"""
    {{
      "equivalence": "related-to",
      "status": "inactive",
      "comment": "Test",
      "elements": [
            {{
              "codeSystemRole": 115,
              "concept": 1
            }},
            {','.join([f'{{"codeSystemRole": {116 + i}, "concept": {random.randint(2, 50)}}}' for i in range(5)])}
      ]
    }}
    """

    conn.request("POST", "/projects/56/mappings", payload, headers)

    res = conn.getresponse()
    data = res.read()

    # print(data.decode("utf-8"))
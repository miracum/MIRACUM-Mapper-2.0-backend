import re
import csv

# Read SQL data from file
with open('tools/codesystem-import/sql_converter/mapper.sql', 'r') as file:
    sql_data = file.read()

# Regular expression to extract code and meaning
pattern = re.compile(r"INSERT INTO public\.sourceterms VALUES \('([^']*)', '([^']*)'\);")

# Extract data
data = pattern.findall(sql_data)

# Write to CSV
with open('sourceterms.csv', 'w', newline='') as csvfile:
    csvwriter = csv.writer(csvfile)
    csvwriter.writerow(['code', 'meaning'])
    csvwriter.writerows(data)

print("CSV file 'sourceterms.csv' created successfully.")
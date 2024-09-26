# CSV Parser

This script parses a `csv` file of e.g. `LOINC` codes and descriptions and convert it into a `csv` file that can be used to import the codes into the `MIRACUM MAPPER 2.0`.

## Requirements

- `Python 3.x`
- dependencies from the `requirements.txt` file installed

## Setup

### 1. Create a Virtual Environment and Install Dependencies

To create and start a virtual environment, open your terminal and run the following command:

```sh
python3 -m venv .venv
source .venv/bin/activate
```

Then install the dependencies:

```sh
pip install -r requirements.txt
```

### 2. Run the Script

To run the script, execute the following command:

```sh
python csv_parser.py -i <path/to/your/input.csv> -o <path/to/your/output.csv>
```

Example usage:

```sh
python parser.py --input=sample-input.csv --config=sample-config.json --output=sample-output.csv
```

--input="tools/codesystem-import/sample/sample-input.csv" --config="tools/codesystem-import/sample/sample-config.json" --output="tools/codesystem-import/sample/sample-output.csv"

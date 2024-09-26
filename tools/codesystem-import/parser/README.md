# CSV Parser

This script parses a `csv` file of e.g. `LOINC` codes and descriptions and convert it into a `csv` file that can be used to import the codes into the `MIRACUM MAPPER 2.0`. A config file has to be provided in order to tell the parser which columns of the input file have to be mapped to which columns of the output file.

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
python csv_parser.py -i <path/to/your/input.csv> -c <path/to/your/config.json> -o <path/to/your/output.csv>
```

Example usage:

```sh
python parser.py --input=sample-input.csv --config=sample-config.json --output=sample-output.csv
```

## Configuration

The configuration file is a `json` file that has to be provided to the script. The keys of the json define the columns of the output file and the values define the columns of the input with a template syntax. An example is given below:

```json
{
  "code": "$LOINC_NUM$",
  "meaning": "Long name: $LONG_COMMON_NAME$ | Short name: $SHORTNAME$"
}
```

In this example the `code` column of the output file is mapped to the `LOINC_NUM` column of the input file. Note that the template syntax is `$<column_name>$` for the columns in the input file. The `meaning` column of the output is handled respectively.

### Example usage

Given the following input file `sample-input.csv`:

| LOINC_NUM | OTHER_COLUMN | LONG_COMMON_NAME | SHORTNAME    |
| --------- | ------------ | ---------------- | ------------ |
| 12345     | other        | Long name 1      | Short name 1 |
| 67890     | other        | Long name 2      | Short name 2 |

And the configuration file from above, the output file `sample-output.csv` will look like this:

| code  | meaning                                            |
| ----- | -------------------------------------------------- |
| 12345 | Long name: Long name 1 \| Short name: Short name 1 |
| 67890 | Long name: Long name 2 \| Short name: Short name 2 |

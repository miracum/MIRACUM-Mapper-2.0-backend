import pandas as pd
import argparse
import json
import re
from tqdm import tqdm

# Regular expression to match placeholders in the template string
templateVariableRegex = r'\$(.*?)\$'

def compile_template(template):
    """
    Compiles a template string by identifying placeholders and creating a format string.

    Parameters:
    template (str): The template string with placeholders.

    Returns:
    tuple: A compiled template format string and a list of placeholders.
    """
    placeholders = re.findall(templateVariableRegex, template)
    format_string = re.sub(templateVariableRegex, '{}', template)
    return format_string, placeholders

def apply_template(format_string, placeholders, row):
    """
    Applies a compiled template to a row by replacing placeholders with corresponding column values.

    Parameters:
    format_string (str): The compiled template format string.
    placeholders (list): A list of placeholders.
    row (namedtuple): A named tuple from the DataFrame containing the data.

    Returns:
    str: The processed template string with placeholders replaced by actual values.
    """
    values = [str(getattr(row, placeholder)) for placeholder in placeholders]
    return format_string.format(*values)

def parse_csv(input_file, output_file, config_file):
    """
    Parses the input CSV file based on the provided configuration and writes the result to the output CSV file.

    Parameters:
    input_file (str): Path to the input CSV file.
    output_file (str): Path to the output CSV file.
    config_file (str): Path to the JSON configuration file.
    """
    # Read the input CSV file
    df = pd.read_csv(input_file)

    # Read the JSON configuration file
    with open(config_file, 'r') as file:
        config = json.load(file)

    # Validate the configuration
    required_keys = ['code', 'meaning']
    for key in required_keys:
            if key not in config:
                raise ValueError(f"Configuration must contain '{key}' key.")

    # Initialize the output DataFrame
    output_df = pd.DataFrame()

    # compile the templates for the required keys and heck if all placeholders exist in the DataFrame columns. Store the placeholders in a map
    placeholders_map, format_string_map = {}, {}
    for key in required_keys:
        format_string, placeholders = compile_template(config[key])
        format_string_map[key], placeholders_map[key] = format_string, placeholders
        for placeholder in placeholders:
            if placeholder not in df.columns:
                raise ValueError(f"Placeholder '{placeholder}' not found in input CSV columns.")

    # Initialize the output DataFrame
    output_df = pd.DataFrame()

    # Compile the templates for the required keys and apply them to each row with a progress bar
    for key in required_keys:
        format_string, placeholders = format_string_map[key], placeholders_map[key]
        output_df[key] = [apply_template(format_string, placeholders, row) for row in tqdm(df.itertuples(index=False, name='Pandas'), total=len(df), desc=f"Processing {key}")]

    # Write the result to the output CSV file
    output_df.to_csv(output_file, index=False)

if __name__ == "__main__":
    # Set up argument parser
    parser = argparse.ArgumentParser(description="Parse a CSV file and extract specific columns based on a JSON configuration.")
    parser.add_argument('-i', '--input', required=True, help="Path to the input CSV file")
    parser.add_argument('-o', '--output', required=True, help="Path to the output CSV file")
    parser.add_argument('-c', '--config', required=True, help="Path to the JSON configuration file")

    # Parse arguments
    args = parser.parse_args()

    # Call the parse_csv function with the provided arguments
    parse_csv(args.input, args.output, args.config)

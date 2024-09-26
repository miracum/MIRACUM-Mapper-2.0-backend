import pandas as pd
import argparse

def parse_csv(input_file, output_file):
    """
    Parses the input CSV file and extracts only the 'LOINC_NUM' and 'LONG_COMMON_NAME' columns.
    Renames these columns to 'code' and 'meaning' respectively, and writes the result to the output CSV file.

    Parameters:
    input_file (str): Path to the input CSV file.
    output_file (str): Path to the output CSV file.
    """
    # Read the input CSV file
    df = pd.read_csv(input_file)

    # Select the required columns
    selected_columns = df[['LOINC_NUM', 'LONG_COMMON_NAME']]

    # Rename the columns
    selected_columns = selected_columns.rename(columns={'LOINC_NUM': 'code', 'LONG_COMMON_NAME': 'meaning'})

    # Write the result to the output CSV file
    selected_columns.to_csv(output_file, index=False)

if __name__ == "__main__":
    # Set up argument parser
    parser = argparse.ArgumentParser(description="Parse a CSV file and extract specific columns.")
    parser.add_argument('-i', '--input', required=True, help="Path to the input CSV file")
    parser.add_argument('-o', '--output', required=True, help="Path to the output CSV file")

    # Parse arguments
    args = parser.parse_args()

    # Call the parse_csv function with the provided arguments
    parse_csv(args.input, args.output)
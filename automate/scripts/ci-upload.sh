#!/bin/bash

# Load environment variables from the .env file
env_file_path=".env"

if [ -f "$env_file_path" ]; then
    while IFS='=' read -r key value; do
        # Skip empty lines and comments
        if [[ -z "$key" || "$key" =~ ^# ]]; then
            continue
        fi
        # Skip empty lines and comments
        if [[ -z "$key" || "$key" =~ ^# ]]; then
            continue
        fi

        # Trim leading and trailing whitespace from the key and value
        key=$(echo "$key" | xargs)
        value=$(echo "$value" | xargs)

        # Validate key name format
        if [[ ! "$key" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]]; then
            echo "Invalid environment variable name: '$key'"
            continue
        fi
        
        # Export the key-value pair as an environment variable
        export "$key=$value"
    done < "$env_file_path"
else
    echo "ERROR: .env file not found."
    exit 1
fi

# Check if ORG_ID or CIRCLECI_TOKEN is empty
if [[ -z "$CIRCLECI_TOKEN" || -z "$ORG_ID"  ]]; then
    echo "ERROR: CIRCLECI_TOKEN or ORG_ID is missing in .env."
    echo "  - Please add CIRCLECI_TOKEN from: https://app.circleci.com/settings/user/tokens"
    echo "  - For ORG_ID, GOTO: https://support.circleci.com/hc/en-us/articles/25148876955163-How-to-Find-Your-Organization-ID"
    exit 1
fi

read -e -p "Enter the path to the circleci .env file and context name (space separated): " env_file_path context_name

if [ ! -f "$env_file_path" ]; then
    echo "ERROR: The specified circle-ci .env file does not exist at the given path: $env_file_path"
    exit 1
fi

# Validate input
if [[ -z "$env_file_path" || -z "$context_name" ]]; then
    echo "ERROR: Both the circleci.env file path and context name are required."
    exit 1
fi

CONTEXT_NAME=$context_name  



# Function to install CircleCI CLI based on the OS
install_circleci_cli() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Installing CircleCI CLI on Linux..."
        curl -fLSs https://circle.ci/cli | sudo bash
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "INFO: Installing CircleCI CLI on macOS..."
        brew install circleci
    else
        echo "ERROR: Unsupported OS for automatic CircleCI CLI installation."
        exit 1
    fi
}

# Check if CircleCI CLI is installed
if ! command -v circleci &> /dev/null; then
    echo "INFO: CircleCI CLI is not installed."
    read -p "Would you like to install CircleCI CLI? (y/n): " install_circleci_choice
    if [[ "$install_circleci_choice" == "y" || "$install_circleci_choice" == "Y" ]]; then
        install_circleci_cli
    else
        echo "ERROR: CircleCI CLI is required to run this script. Please install it manually."
        exit 1
    fi
fi

# Get the context ID using CircleCI CLI
context_output=$(circleci context list --org-id=$ORG_ID --token=$CIRCLECI_TOKEN)

# Check if CONTEXT_NAME exists
if echo "$context_output" | grep -q "$CONTEXT_NAME"; then
    echo "INFO: FOUND THE CONTEXT: $CONTEXT_NAME"
else
    echo "WARNING: Context '$CONTEXT_NAME' not found."
    read -p "Would you like to create a new context? (y/n): " create_context

    if [[ "$create_context" =~ ^[Yy]$ ]]; then
        # Add command to create the context here
        echo "INFO: Creating a new context..."
        echo $(circleci context create $CONTEXT_NAME --org-id=$ORG_ID --token=$CIRCLECI_TOKEN)
    else
        echo "INFO: Exiting without creating a new context."
        exit 1
    fi
fi

echo "Uploading environment variables...."

# Read the circleci.env file and upload each key-value pair as a secret
while IFS= read -r line || [[ -n "$line" ]]; do
    # Skip comments and empty lines
    if [[ $line == \#* ]] || [[ -z $line ]]; then
        continue
    fi

    # Split the line into key and value
    IFS='=' read -r key value <<< "$line"

    if [ -z "$key" ] || [ -z "$value" ]; then
        echo "Invalid line in $env_file_path file: $line"
        continue
    fi

    # Trim whitespace
    key=$(echo "$key" | xargs)
    value=$(echo "$value" | xargs)

    # Store the secret in CircleCI context
    echo "$value" | circleci context store-secret --org-id="$ORG_ID" --token="$CIRCLECI_TOKEN" "$CONTEXT_NAME" "$key"

    if [ $? -eq 0 ]; then
        echo "INFO: Uploaded secret: $key"
    else
        echo "ERROR: Failed to upload secret: $key"
    fi
done < "$env_file_path"

echo "Done uploading environment variables."

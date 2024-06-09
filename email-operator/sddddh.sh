#!/bin/bash

# Define the old and new strings
old_string="email-operator/api"
new_string="email-operator/api"

# Find all files in the current directory and its subdirectories
find . -type f -exec sed -i "s|$old_string|$new_string|g" {} +

echo "Replacement complete."

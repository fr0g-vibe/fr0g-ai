#!/bin/bash

# Unicode Purge Script for fr0g.ai Project
# Removes unicode symbols and replaces them with plain text alternatives

set -e

echo "Starting unicode purge across fr0g.ai project..."

# Define unicode symbol mappings
declare -A UNICODE_MAP=(
    ["✅"]="COMPLETED"
    ["❌"]="FAILED"
    ["🔥"]="PRIORITY"
    ["🚀"]="STARTING"
    ["⚡"]="PERFORMANCE"
    ["🎯"]="TARGET"
    ["🧪"]="TESTING"
    ["🔍"]="CHECKING"
    ["🔨"]="BUILDING"
    ["📦"]="INSTALLING"
    ["🐸"]="fr0g.ai"
    ["🏥"]="HEALTH"
    ["🛠"]="SETUP"
    ["🔒"]="SECURITY"
    ["💡"]="TIP"
    ["⏳"]="WAITING"
    ["🐳"]="DOCKER"
    ["🧹"]="CLEANING"
)

# Function to purge unicode from a file
purge_file() {
    local file="$1"
    echo "Purging unicode from: $file"
    
    # Skip binary files
    if file "$file" | grep -q "ELF\|executable"; then
        echo "  Skipping binary file: $file"
        return
    fi
    
    # Skip JSON data files (they might contain legitimate unicode)
    if [[ "$file" == *.json ]] && [[ "$file" == *"/data/"* ]]; then
        echo "  Skipping data file: $file"
        return
    fi
    
    # Create backup
    cp "$file" "$file.backup"
    
    # Apply replacements
    local temp_file=$(mktemp)
    cp "$file" "$temp_file"
    
    for unicode in "${!UNICODE_MAP[@]}"; do
        replacement="${UNICODE_MAP[$unicode]}"
        sed -i "s/$unicode/$replacement/g" "$temp_file"
    done
    
    # Check if file was modified
    if ! cmp -s "$file" "$temp_file"; then
        mv "$temp_file" "$file"
        echo "  Modified: $file"
    else
        rm "$temp_file"
        rm "$file.backup"
        echo "  No changes: $file"
    fi
}

# First, clean up any existing backup files
echo "Cleaning up existing backup files..."
find . -name '*.backup' -delete 2>/dev/null || true

# Get list of files with unicode symbols (excluding binary data and dot files)
files_with_unicode=$(find . -type f -not -path '*/.*' -not -path '*/data/*' -not -path '*/bin/*' -exec grep -l '✅\|❌\|🔥\|🚀\|⚡\|🎯\|🧪\|🔍\|🔨\|📦\|🐸\|🏥\|🛠\|🔒\|💡\|⏳\|🐳\|🧹' {} \; 2>/dev/null || true)

# Filter out binary files
filtered_files=""
while IFS= read -r file; do
    if [ -f "$file" ] && ! file "$file" | grep -q "ELF\|executable\|binary"; then
        filtered_files="$filtered_files$file"$'\n'
    fi
done <<< "$files_with_unicode"

files_with_unicode="$filtered_files"

if [ -z "$files_with_unicode" ]; then
    echo "No files with unicode symbols found."
    exit 0
fi

echo "Found files with unicode symbols:"
echo "$files_with_unicode"
echo ""

# Process each file
while IFS= read -r file; do
    if [ -f "$file" ]; then
        purge_file "$file"
    fi
done <<< "$files_with_unicode"

echo ""
echo "Unicode purge completed!"
echo "All unicode symbols have been replaced with plain text alternatives."
echo "Files have been processed in-place without backup files."

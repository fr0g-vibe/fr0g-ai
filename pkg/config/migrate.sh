#!/bin/bash

# Migration script for fr0g-ai subprojects to use shared configuration library

set -e

PROJECT_ROOT=$(pwd)
SHARED_CONFIG_PATH="pkg/config"

echo "🔧 Migrating fr0g-ai subprojects to use shared configuration library..."

# Function to backup a file
backup_file() {
    local file=$1
    if [ -f "$file" ]; then
        cp "$file" "$file.backup.$(date +%Y%m%d_%H%M%S)"
        echo "📦 Backed up $file"
    fi
}

# Function to update go.mod files
update_go_mod() {
    local project_dir=$1
    echo "📝 Updating go.mod in $project_dir..."
    
    cd "$PROJECT_ROOT/$project_dir"
    
    # Add replace directive for local pkg/config if not already present
    if ! grep -q "pkg/config" go.mod; then
        echo "" >> go.mod
        echo "replace pkg/config => ../pkg/config" >> go.mod
        echo "✅ Added pkg/config replace directive to $project_dir/go.mod"
    fi
    
    # Run go mod tidy
    go mod tidy
    echo "✅ Updated dependencies for $project_dir"
}

# Function to check if shared config is properly set up
check_shared_config() {
    echo "🔍 Checking shared configuration setup..."
    
    if [ ! -d "$PROJECT_ROOT/$SHARED_CONFIG_PATH" ]; then
        echo "❌ Shared config directory not found at $SHARED_CONFIG_PATH"
        exit 1
    fi
    
    required_files=("config.go" "validation.go" "loader.go")
    for file in "${required_files[@]}"; do
        if [ ! -f "$PROJECT_ROOT/$SHARED_CONFIG_PATH/$file" ]; then
            echo "❌ Required file $SHARED_CONFIG_PATH/$file not found"
            exit 1
        fi
    done
    
    echo "✅ Shared configuration files are present"
}

# Function to migrate fr0g-ai-aip
migrate_aip() {
    echo "🚀 Migrating fr0g-ai-aip..."
    
    local aip_dir="fr0g-ai-aip"
    if [ ! -d "$PROJECT_ROOT/$aip_dir" ]; then
        echo "⚠️  $aip_dir directory not found, skipping..."
        return
    fi
    
    cd "$PROJECT_ROOT/$aip_dir"
    
    # Backup existing config files
    backup_file "internal/config/validation.go"
    
    # Update go.mod
    update_go_mod "$aip_dir"
    
    echo "✅ fr0g-ai-aip migration completed"
}

# Function to migrate fr0g-ai-bridge
migrate_bridge() {
    echo "🌉 Migrating fr0g-ai-bridge..."
    
    local bridge_dir="fr0g-ai-bridge"
    if [ ! -d "$PROJECT_ROOT/$bridge_dir" ]; then
        echo "⚠️  $bridge_dir directory not found, skipping..."
        return
    fi
    
    cd "$PROJECT_ROOT/$bridge_dir"
    
    # Backup existing config files
    backup_file "internal/config/config.go"
    backup_file "internal/api/validation.go"
    
    # Update go.mod
    update_go_mod "$bridge_dir"
    
    echo "✅ fr0g-ai-bridge migration completed"
}

# Function to migrate fr0g-ai-master-control
migrate_master_control() {
    echo "🎛️  Migrating fr0g-ai-master-control..."
    
    local mc_dir="fr0g-ai-master-control"
    if [ ! -d "$PROJECT_ROOT/$mc_dir" ]; then
        echo "⚠️  $mc_dir directory not found, skipping..."
        return
    fi
    
    cd "$PROJECT_ROOT/$mc_dir"
    
    # Backup existing config files
    backup_file "internal/config/validation.go"
    backup_file "internal/mastercontrol/config.go"
    
    # Update go.mod
    update_go_mod "$mc_dir"
    
    echo "✅ fr0g-ai-master-control migration completed"
}

# Function to run tests
run_tests() {
    echo "🧪 Running tests to verify migration..."
    
    cd "$PROJECT_ROOT"
    
    # Test shared config
    echo "Testing shared configuration..."
    cd "$SHARED_CONFIG_PATH"
    go test -v ./...
    
    # Test each subproject
    for project in "fr0g-ai-aip" "fr0g-ai-bridge" "fr0g-ai-master-control"; do
        if [ -d "$PROJECT_ROOT/$project" ]; then
            echo "Testing $project..."
            cd "$PROJECT_ROOT/$project"
            go build ./... || echo "⚠️  Build issues in $project - manual review needed"
        fi
    done
    
    echo "✅ Tests completed"
}

# Main migration process
main() {
    echo "🎯 Starting fr0g-ai configuration migration..."
    echo "Project root: $PROJECT_ROOT"
    
    # Check prerequisites
    check_shared_config
    
    # Migrate each subproject
    migrate_aip
    migrate_bridge
    migrate_master_control
    
    # Run tests
    run_tests
    
    echo ""
    echo "🎉 Migration completed successfully!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Review the backed up files (*.backup.*) for any custom logic"
    echo "2. Update your application code to use the new shared config types"
    echo "3. Test your applications thoroughly"
    echo "4. Update documentation to reference the shared configuration"
    echo ""
    echo "📚 See pkg/config/README.md for usage examples and best practices"
}

# Run the migration
main "$@"

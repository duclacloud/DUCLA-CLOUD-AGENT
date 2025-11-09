# 🖥️ Ducla Cloud Agent CLI Features

## 📋 Overview

Ducla Cloud Agent v1.0.0 bây giờ đã có đầy đủ CLI commands và man page, làm cho việc sử dụng và quản lý agent trở nên dễ dàng hơn nhiều.

## ✨ New Features Added

### 1. 📚 Man Page
- **Complete manual**: `man ducla-agent`
- **Professional documentation** với tất cả options, commands, examples
- **Installed system-wide** trong `/usr/share/man/man1/`
- **Searchable** với man database

### 2. 🖥️ CLI Commands
Comprehensive command-line interface với các nhóm lệnh:

#### **Show Commands**
```bash
ducla-agent show status          # Agent runtime status
ducla-agent show health          # System health information  
ducla-agent show metrics         # Current metrics
ducla-agent show config          # Current configuration
ducla-agent show tasks           # List all tasks
ducla-agent show tasks running   # List only running tasks
ducla-agent show version         # Version information
```

#### **Task Management**
```bash
ducla-agent task create "command"  # Create and execute task
ducla-agent task cancel TASK_ID    # Cancel running task
ducla-agent task logs TASK_ID      # Show task logs
```

#### **File Operations**
```bash
ducla-agent file list PATH        # List files in directory
ducla-agent file copy SRC DEST    # Copy file
ducla-agent file move SRC DEST    # Move file  
ducla-agent file delete PATH      # Delete file/directory
```

#### **Configuration**
```bash
ducla-agent config validate       # Validate config file
ducla-agent config test          # Test config and connectivity
```

### 3. 🔧 Enhanced Options
```bash
-c, --config FILE    # Configuration file path
-d, --debug          # Enable debug logging
-v, --version        # Show version information
-h, --help           # Show help message
```

## 🎯 Usage Examples

### Basic Usage
```bash
# Start agent daemon
ducla-agent

# Show help
ducla-agent --help

# Show version
ducla-agent show version
```

### Monitoring & Status
```bash
# Check agent status
ducla-agent show status

# Check system health
ducla-agent show health

# View current configuration
ducla-agent show config
```

### Task Management
```bash
# Create a simple task
ducla-agent task create "ls -la /tmp"

# List all tasks
ducla-agent show tasks

# List only running tasks
ducla-agent show tasks running
```

### File Operations
```bash
# List files in directory
ducla-agent file list /var/log

# Copy a file
ducla-agent file copy /tmp/source.txt /tmp/backup.txt

# Move a file
ducla-agent file move /tmp/old.txt /tmp/new.txt

# Delete a file
ducla-agent file delete /tmp/unwanted.txt
```

### Configuration Management
```bash
# Validate configuration
ducla-agent config validate

# Test configuration and connectivity
ducla-agent config test

# Use custom config file
ducla-agent -c /path/to/custom.yaml show status
```

## 📖 Documentation

### Man Page Sections
The man page includes comprehensive documentation:

- **NAME & SYNOPSIS**: Command overview
- **DESCRIPTION**: Detailed explanation
- **OPTIONS**: All command-line options
- **COMMANDS**: Complete command reference
- **CONFIGURATION**: Config file documentation
- **API ENDPOINTS**: REST API reference
- **EXAMPLES**: Practical usage examples
- **FILES**: Important file locations
- **ENVIRONMENT**: Environment variables
- **SYSTEMD SERVICE**: Service management
- **MONITORING**: Health checks and metrics
- **SECURITY**: Security features
- **TROUBLESHOOTING**: Common issues and solutions

### Access Documentation
```bash
# View complete manual
man ducla-agent

# Search in manual
man ducla-agent | grep -i "task"

# View specific sections
man ducla-agent | less +/EXAMPLES
```

## 🚀 Benefits

### For Users
- **Easy to use**: Intuitive command structure
- **Self-documenting**: Built-in help and man page
- **Flexible**: Works with or without config files
- **Comprehensive**: All features accessible via CLI

### For Administrators
- **Scriptable**: Perfect for automation scripts
- **Monitorable**: Easy status and health checking
- **Debuggable**: Debug mode and detailed logging
- **Manageable**: Config validation and testing

### For Developers
- **API-compatible**: CLI mirrors REST API functionality
- **Extensible**: Easy to add new commands
- **Consistent**: Follows Unix CLI conventions
- **Professional**: Production-ready interface

## 🔍 Technical Implementation

### Architecture
- **Modular design**: Separate CLI handler (`cli.go`)
- **Config-aware**: Automatically loads configuration
- **Error handling**: Proper error messages and exit codes
- **HTTP client**: Communicates with running agent via REST API

### Error Handling
- **Graceful degradation**: Works with default config if file missing
- **Clear messages**: Descriptive error messages
- **Proper exit codes**: Standard Unix exit codes
- **Debug support**: Verbose logging with `-d` flag

### Performance
- **Fast startup**: Quick command execution
- **Low overhead**: Minimal resource usage for CLI commands
- **Efficient**: Direct HTTP API communication
- **Responsive**: Quick response times

## 📦 Package Integration

### DEB/RPM Packages
- **Man page included**: Automatically installed with package
- **System integration**: Proper file locations
- **Service management**: Works with systemd
- **Clean uninstall**: Proper cleanup on removal

### Installation
```bash
# Install DEB package
sudo dpkg -i ducla-agent_1.0.0_amd64.deb

# Man page automatically available
man ducla-agent

# CLI commands ready to use
ducla-agent --help
```

## 🎉 Demo & Testing

### CLI Demo Script
```bash
# Run CLI demo
./demo-cli.sh
```

### Manual Testing
```bash
# Test all CLI features
./dist/ducla-agent --help
./dist/ducla-agent show version
./dist/ducla-agent show config
./dist/ducla-agent config validate
man ducla-agent
```

## 🚀 Future Enhancements

### Planned Features
- **Interactive mode**: `ducla-agent interactive`
- **Batch operations**: Multiple file operations
- **Advanced filtering**: Complex task queries
- **Export/Import**: Configuration management
- **Plugins**: Extensible command system

### Community Contributions
- **Custom commands**: Plugin system for custom commands
- **Shell completion**: Bash/Zsh completion scripts
- **GUI wrapper**: Optional graphical interface
- **Integration**: IDE and editor plugins

## 📊 Summary

### ✅ Completed Features
- [x] Complete man page documentation
- [x] Comprehensive CLI command set
- [x] Show commands (status, health, metrics, config, tasks, version)
- [x] Task management (create, cancel, logs)
- [x] File operations (list, copy, move, delete)
- [x] Configuration management (validate, test)
- [x] Help system and error handling
- [x] Package integration (DEB/RPM)
- [x] Demo scripts and testing

### 🎯 Impact
- **User Experience**: Dramatically improved usability
- **Administration**: Easier system management
- **Automation**: Better scripting capabilities
- **Documentation**: Professional-grade documentation
- **Adoption**: Lower barrier to entry

**Ducla Cloud Agent v1.0.0** is now a complete, professional-grade cloud agent with full CLI support and comprehensive documentation! 🎊
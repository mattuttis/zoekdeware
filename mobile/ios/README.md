# ZoekDeware iOS App

SwiftUI-based dating app for iPhone 14 and newer (iOS 16+).

## Requirements

- Xcode 15+
- iOS 16.0+ deployment target
- [XcodeGen](https://github.com/yonaskolb/XcodeGen) for project generation

## Setup

```bash
# Install XcodeGen and generate project
make setup

# Open in Xcode
make open
```

## Project Structure

```
ZoekDeware/
├── App/                    # App entry point
│   ├── ZoekDewareApp.swift
│   ├── ContentView.swift
│   └── Config.swift
├── Core/                   # Core infrastructure
│   ├── Network/
│   │   ├── APIClient.swift
│   │   └── Endpoints.swift
│   ├── Storage/
│   │   └── TokenStorage.swift
│   └── Theme/
├── Features/               # Feature modules
│   ├── Auth/
│   │   ├── Models/
│   │   ├── ViewModels/
│   │   └── Views/
│   ├── Profile/
│   ├── Discover/
│   └── Chat/
└── Resources/
    ├── Info.plist
    └── Assets.xcassets
```

## Architecture

- **SwiftUI** for UI
- **MVVM** pattern
- **Swift Concurrency** (async/await)
- **Keychain** for secure token storage

## Building

```bash
# Build for simulator
make build

# Run tests
make test
```

## Configuration

Edit `App/Config.swift` to configure API endpoints:

- Development: `http://localhost:8000/api/v1`
- Staging: `https://staging-api.zoekdeware.app/api/v1`
- Production: `https://api.zoekdeware.app/api/v1`

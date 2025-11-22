import Foundation

enum Config {
    enum Environment {
        case development
        case staging
        case production
    }

    static let current: Environment = {
        #if DEBUG
        return .development
        #else
        return .production
        #endif
    }()

    static var apiBaseURL: URL {
        switch current {
        case .development:
            return URL(string: "http://localhost/api/v1/")!
        case .staging:
            return URL(string: "https://staging-api.zoekdeware.app/api/v1/")!
        case .production:
            return URL(string: "https://api.zoekdeware.app/api/v1/")!
        }
    }

    static var webSocketURL: URL {
        switch current {
        case .development:
            return URL(string: "ws://localhost/api/v1/ws")!
        case .staging:
            return URL(string: "wss://staging-api.zoekdeware.app/api/v1/ws")!
        case .production:
            return URL(string: "wss://api.zoekdeware.app/api/v1/ws")!
        }
    }
}
